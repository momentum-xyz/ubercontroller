package node

import (
	"context"
	"encoding/json"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/types/generic"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/momentum-xyz/ubercontroller/utils"
)

var updateProfileStore = generic.NewSyncMap[umid.UMID, UpdateProfileStoreItem](0)

type UpdateProfileStoreItem struct {
	Status      string
	NodeJSOut   *NodeJSOut
	UserID      umid.UMID
	UserProfile *entry.UserProfile
	Error       error
}

// @Summary Check update user profile job by Job UMID
// @Schemes
// @Description Returns Update Profile Job UMID status
// @Tags profile
// @Accept json
// @Produce json
// @Param job_id path string true "Job UMID"
// @Success 200 {object} node.apiProfileUpdateCheckJob.Out
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/profile/check-job/{job_id} [get]
func (n *Node) apiProfileUpdateCheckJob(c *gin.Context) {
	type Out struct {
		NodeJSOut *NodeJSOut         `json:"nodeJSOut"`
		Status    string             `json:"status"`
		JobID     umid.UMID          `json:"job_id"`
		UserID    umid.UMID          `json:"user_id"`
		Error     *string            `json:"error"`
		Profile   *entry.UserProfile `json:"profile"`
	}

	jobID, err := umid.Parse(c.Param("jobID"))
	if err != nil {
		err = errors.WithMessage(err, "Node: apiProfileUpdateCheckJob: failed to parse uuid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_param", err, n.log)
		return
	}

	item, ok := updateProfileStore.Load(jobID)
	if !ok {
		item.Status = "job not found"
	}

	var message *string
	if item.Error != nil {
		e := item.Error.Error()
		message = &e
	}

	out := Out{
		JobID:     jobID,
		Status:    item.Status,
		NodeJSOut: item.NodeJSOut,
		Profile:   item.UserProfile,
		UserID:    item.UserID,
		Error:     message,
	}

	c.JSON(http.StatusOK, out)
}

// @Summary Edit user profile
// @Schemes
// @Description Edits a user profile
// @Tags profile
// @Accept json
// @Produce json
// @Param body body node.apiProfileUpdate.Body true "body params"
// @Success 200 {object} node.apiProfileUpdate.Out
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Router /api/v4/profile [patch]
func (n *Node) apiProfileUpdate(c *gin.Context) {
	type Body struct {
		Name    *string `json:"name"`
		Profile *struct {
			Bio         *string `json:"bio"`
			ProfileLink *string `json:"profileLink"`
			Location    *string `json:"location"`
			AvatarHash  *string `json:"avatarHash"`
		} `json:"profile"`
	}

	inBody := Body{}

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiProfileUpdate: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiProfileUpdate: failed to get user umid from context")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	userProfile, err := n.db.GetUsersDB().GetUserProfileByUserID(c, userID)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiProfileUpdate: failed to get user profile by user umid")
		api.AbortRequest(c, http.StatusNotFound, "not_found", err, n.log)
		return
	}

	nameChanged := inBody.Name != nil && *userProfile.Name != *inBody.Name
	avatarChanged := inBody.Profile != nil && inBody.Profile.AvatarHash != nil && userProfile.AvatarHash != nil && *userProfile.AvatarHash != *inBody.Profile.AvatarHash
	shouldUpdateNFT := nameChanged || avatarChanged

	if inBody.Name != nil {
		// TODO: check name unique
		userProfile.Name = inBody.Name
	}

	inProfile := inBody.Profile
	if inProfile != nil {
		if inProfile.Bio != nil {
			userProfile.Bio = inProfile.Bio
		}
		if inProfile.Location != nil {
			userProfile.Location = inProfile.Location
		}
		if inProfile.AvatarHash != nil {
			userProfile.AvatarHash = inProfile.AvatarHash
		}
		if inProfile.ProfileLink != nil {
			userProfile.ProfileLink = inProfile.ProfileLink
		}
	}

	userProfile.OnBoarded = utils.GetPTR(true)

	type Out struct {
		JobID  *umid.UMID `json:"job_id"`
		UserID umid.UMID  `json:"user_id"`
	}

	if !shouldUpdateNFT {
		// If no need to update NFT meta, execute update synchronously and return JobID:null
		if err := n.db.GetUsersDB().UpdateUserProfile(c, userID, userProfile); err != nil {
			err = errors.WithMessage(err, "Node: apiProfileUpdate: failed to update user profile")
			api.AbortRequest(c, http.StatusNotFound, "not_found", err, n.log)
			return
		}

		out := Out{
			JobID:  nil,
			UserID: userID,
		}

		c.JSON(http.StatusOK, out)
		return
	}

	jobID := umid.New()
	updateProfileStore.Store(
		jobID, UpdateProfileStoreItem{
			Status:    StatusInProgress,
			NodeJSOut: nil,
		},
	)

	// Can not use gin context, because worker go-routine should continue after response
	ctx := context.Background()
	go n.updateUserProfileWorker(ctx, jobID, userID, userProfile)

	out := Out{
		JobID:  &jobID,
		UserID: userID,
	}

	c.JSON(http.StatusOK, out)
}

func (n *Node) updateUserProfileWorker(
	ctx context.Context, jobID umid.UMID, userID umid.UMID, userProfile *entry.UserProfile,
) {
	item := UpdateProfileStoreItem{
		Status:    "",
		NodeJSOut: nil,
		Error:     nil,
	}

	wallet, err := n.db.GetUsersDB().GetUserWalletByUserID(ctx, userID)
	if err != nil {
		err = errors.WithMessage(err, "failed to get user wallet by userID")
		{
			item.Status = StatusFailed
			item.Error = err
			updateProfileStore.Store(jobID, item)
		}
		log.Error(err)
		return
	}

	meta := NFTMeta{
		Name:  "",
		Image: "",
	}

	if userProfile.Name != nil {
		meta.Name = *userProfile.Name
	}

	if userProfile.AvatarHash != nil {
		meta.Image = *userProfile.AvatarHash
	}

	b, err := json.Marshal(meta)
	if err != nil {
		err = errors.WithMessage(err, "failed to json.Marshal meta to nodejs input")
		{
			item.Status = StatusFailed
			item.Error = err
			updateProfileStore.Store(jobID, item)
		}
		log.Error(err)
		return
	}

	output, err := exec.Command(
		"node", "./nodejs/check-nft/update-nft.js", *wallet, n.cfg.Common.MnemonicPhrase, string(b), userID.String(),
	).Output()
	if err != nil {
		err = errors.WithMessage(err, "failed to execute node script update-nft.js")
		{
			item.Status = StatusFailed
			item.Error = err
			updateProfileStore.Store(jobID, item)
		}
		log.Error(err)
		return
	}

	var nodeJSOut NodeJSOut
	if err := json.Unmarshal(output, &nodeJSOut); err != nil {
		err = errors.WithMessage(err, "failed to unmarshal output update-nft.js")
		{
			item.Status = StatusFailed
			item.NodeJSOut = &nodeJSOut
			item.Error = err
			updateProfileStore.Store(jobID, item)
		}
		log.Error(err)
		return
	}

	// Update DB
	if err := n.db.GetUsersDB().UpdateUserProfile(ctx, userID, userProfile); err != nil {
		err = errors.WithMessage(err, "failed to update user profile")
		{
			item.Status = StatusFailed
			item.Error = err
			updateProfileStore.Store(jobID, item)
		}
		log.Error(err)
		return
	}

	item.Status = StatusDone
	item.Error = nil
	item.UserProfile = userProfile
	item.UserID = userID
	updateProfileStore.Store(jobID, item)
}
