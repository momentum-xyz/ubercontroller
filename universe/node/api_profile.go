package node

import (
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe/api"
	"github.com/momentum-xyz/ubercontroller/utils"
)

type Form struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

func (n *Node) apiProfileUpdate(c *gin.Context) {
	inBody := struct {
		Name    *string `json:"name"`
		Profile *struct {
			Name        *string `json:"name"`
			Bio         *string `json:"bio"`
			ProfileLink *string `json:"profileLink"`
			Location    *string `json:"location"`
			AvatarHash  *string `json:"avatarHash"`
		} `json:"profile"`
	}{}

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

func (n *Node) apiProfileUploadAvatar(c *gin.Context) {
	var form Form
	_ := c.ShouldBind(&form)
	openedFile, _ := form.File.Open()
	file, _ := io.ReadAll(openedFile)

	resp, err := http.PostForm("https://httpbin.org/post", file)
	if err != nil {
		err = errors.WithMessage(err, "failed to update user avatar")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_update_user_avatar", err, n.log)
		return
	}

	var res map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&res)
	// c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}
