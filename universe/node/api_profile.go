package node

import (
	"net/http"

	"github.com/momentum-xyz/ubercontroller/utils/umid"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/momentum-xyz/ubercontroller/utils"
)

// @Summary Edit user profile
// @Description Edits a user profile
// @Tags users,profile
// @Security Bearer
// @Param body body node.apiProfileUpdate.Body true "body params"
// @Success 200 {object} node.apiProfileUpdate.Out
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
