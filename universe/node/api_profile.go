package node

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe/api"
	"github.com/momentum-xyz/ubercontroller/utils"
)

func (n *Node) apiProfileEdit(c *gin.Context) {
	inBody := struct {
		Name    *string `json:"name"`
		Profile *struct {
			Bio         *string `json:"bio"`
			Location    *string `json:"location"`
			AvatarHash  *string `json:"avatarHash"`
			ProfileLink *string `json:"profileLink"`
			OnBoarded   *bool   `json:"onBoarded"`
		} `json:"profile"`
	}{}

	if err := c.ShouldBindJSON(&inBody); err != nil {
		n.log.Warn(errors.WithMessage(err, "Node: apiProfileEdit: failed to bind json"))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
		})
		return
	}

	userID, err := api.GetUserIDFromRequest(c)
	if err != nil {
		n.log.Warn(errors.WithMessage(err, "Node: apiProfileEdit: failed to get user id from request"))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "invalid access token",
		})
		return
	}

	userProfile, err := n.db.UsersGetUserProfileByUserID(c, userID)
	if err != nil {
		n.log.Warn(errors.WithMessage(err, "Node: apiProfileEdit: failed to get user profile by user id"))
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "not found",
		})
		return
	}

	userProfile.OnBoarded = utils.GetPTR(true)

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

	if err := n.db.UsersUpdateUserProfile(c, userID, userProfile); err != nil {
		n.log.Warn(errors.WithMessage(err, "failed to update user profile"))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "failed to update user profile",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"userOnboarded": userProfile.OnBoarded,
	})
}
