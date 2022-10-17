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
		n.log.Debug(errors.WithMessage(err, "Node: apiProfileUpdate: failed to bind json"))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
		})
		return
	}

	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		n.log.Debug(errors.WithMessage(err, "Node: apiProfileUpdate: failed to get user id from context"))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "invalid user id",
		})
		return
	}

	userProfile, err := n.db.UsersGetUserProfileByUserID(c, userID)
	if err != nil {
		n.log.Debug(errors.WithMessage(err, "Node: apiProfileUpdate: failed to get user profile by user id"))
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "not found",
		})
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
		n.log.Error(errors.WithMessage(err, "failed to update user profile"))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "failed to update user profile",
		})
		return
	}

	c.JSON(http.StatusOK, nil)
}
