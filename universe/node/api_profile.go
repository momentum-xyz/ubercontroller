package node

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe/api"
	"github.com/momentum-xyz/ubercontroller/utils"
)

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
		n.log.Debug(err)
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err)
		return
	}

	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiProfileUpdate: failed to get user id from context")
		n.log.Debug(err)
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err)
		return
	}

	userProfile, err := n.db.UsersGetUserProfileByUserID(c, userID)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiProfileUpdate: failed to get user profile by user id")
		n.log.Debug(err)
		api.AbortRequest(c, http.StatusNotFound, "not_found", err)
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
		n.log.Error(err)
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_update_user_profile", err)
		return
	}

	n.apiUsersGetMe(c)
}
