package node

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"

	"github.com/momentum-xyz/ubercontroller/universe/api"
	"github.com/momentum-xyz/ubercontroller/utils"
)

// @BasePath /api/v4

// @Summary Edits a profile
// @Schemes
// @Description Edits a profile
// @Tags profile
// @Accept json
// @Produce json
// @Param request body node.apiProfileUpdate.Body true "body params"
// @Success 200 {object} dto.User
// @Success 500 {object} api.HTTPError
// @Success 400 {object} api.HTTPError
// @Router /api/v4/profile [patch]
func (n *Node) apiProfileUpdate(c *gin.Context) {
	type Body struct {
		Name    *string `json:"name"`
		Profile *struct {
			Name        *string `json:"name"`
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
		err = errors.WithMessage(err, "Node: apiProfileUpdate: failed to get user id from context")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	userProfile, err := n.db.UsersGetUserProfileByUserID(c, userID)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiProfileUpdate: failed to get user profile by user id")
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

	if err := n.db.UsersUpdateUserProfile(c, userID, userProfile); err != nil {
		err = errors.WithMessage(err, "failed to update user profile")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_update_user_profile", err, n.log)
		return
	}

	n.apiUsersGetMe(c)
}

// @BasePath /api/v4

// @Summary Uploads a user avatar to media-manager
// @Schemes
// @Description Sends an image file to the media manager and returns a hash
// @Tags profile
// @Accept json
// @Produce string
// @Param request body node.apiProfileUpdate.Body true "body params"
// @Success 200 {object} dto.User
// @Success 500 {object} api.HTTPError
// @Success 400 {object} api.HTTPError
// @Router /api/v4/profile [patch]
func (n *Node) apiProfileUploadAvatar(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		err := errors.WithMessage(err, "Node: apiProfileUploadAvatar: failed to read file")
		api.AbortRequest(c, http.StatusBadRequest, "failed_to_read", err, n.log)
		return
	}

	openedFile, err := file.Open()
	defer openedFile.Close()
	if err != nil {
		err := errors.WithMessage(err, "Node: apiProfileUploadAvatar: failed to open file")
		api.AbortRequest(c, http.StatusBadRequest, "failed_to_open", err, n.log)
		return
	}

	req, err := http.NewRequest("POST", n.cfg.Common.RenderDefaultUrl, openedFile)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiProfileUploadAvatar: failed to post data to media-manager")
		api.AbortRequest(c, http.StatusBadRequest, "failed_to_post", err, n.log)
		return
	}

	req.Header.Set("Content-Type", "image/png")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiProfileUploadAvatar: failed to post data to media-manager")
		api.AbortRequest(c, http.StatusBadRequest, "failed_to_post", err, n.log)
		return
	}

	fmt.Sprintln(resp)
}
