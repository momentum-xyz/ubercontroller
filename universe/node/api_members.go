package node

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/dto"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

// @Summary Get members belonging to an object
// @Schemes
// @Description Returns members belonging to the object
// @Tags members
// @Accept json
// @Produce json
// @Param object_id path string true "Object UMID"
// @Success 200 {object} dto.Member
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/objects/{object_id}/members [get]
func (n *Node) apiMembersGetForObject(c *gin.Context) {
	objectID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiMembersGetForObject: failed to parse object umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	allUserObjects := n.GetUserObjects()
	filteredUserObjects, err := allUserObjects.GetUserObjectsByObjectID(objectID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiMembersGetForObject: failed to get user objects")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_user_objects", err, n.log)
		return
	}

	members := make([]dto.Member, 0, len(filteredUserObjects))
	for _, filteredUserObject := range filteredUserObjects {
		user, err := n.LoadUser(filteredUserObject.UserID)
		if err != nil {
			err := errors.WithMessage(err, "Node: apiMembersGetForObject: failed to load user")
			api.AbortRequest(c, http.StatusInternalServerError, "failed_to_load_user", err, n.log)
			return
		}

		var avatarHash, userName, userRole string
		profile := user.GetProfile()
		if profile != nil {
			avatarHash = ""
			if profile.AvatarHash != nil {
				avatarHash = *profile.AvatarHash
			}

			userName = ""
			if profile.Name != nil {
				userName = *profile.Name
			}
		}

		userObjectValue := filteredUserObject.Value
		if userObjectValue != nil {
			userRole = utils.GetFromAnyMap(*userObjectValue, universe.ReservedAttributes.User.Role.Key, "")
		}

		member := dto.Member{
			ObjectID:   &filteredUserObject.ObjectID,
			UserID:     &filteredUserObject.UserID,
			Name:       &userName,
			AvatarHash: &avatarHash,
			Role:       &userRole,
		}

		members = append(members, member)
	}

	c.JSON(http.StatusOK, members)
}
