package node

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/dto"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
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

// @Summary Add member to object
// @Schemes
// @Description Add member to object
// @Tags members
// @Accept json
// @Produce json
// @Param body body node.apiPostMemberForObject.Body true "body params"
// @Success 200 {object} nil
// @Failure 400 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// /api/v4/objects/{object_id}/members [post]
func (n *Node) apiPostMemberForObject(c *gin.Context) {

	objectID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiPostMemberForObject: failed to parse object umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	type Body struct {
		UserID *umid.UMID `json:"user_id"`
		Wallet *string    `json:"wallet"`
		Role   string     `json:"role" binding:"required"`
	}

	var inBody Body
	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiPostMemberForObject: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	if inBody.Role != "admin" {
		err := errors.New("Node: apiPostMemberForObject: role not allowed")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	if inBody.Wallet == nil && inBody.UserID == nil {
		err := errors.New("Node: apiPostMemberForObject: user_id or wallet must be provided")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	if inBody.Wallet != nil && inBody.UserID != nil {
		err := errors.New("Node: apiPostMemberForObject: only one parameter should be provided: user_id or wallet")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	var userID umid.UMID
	if inBody.Wallet != nil {
		user, err := n.db.GetUsersDB().GetUserByWallet(c, *inBody.Wallet)
		if err != nil {
			err = errors.WithMessage(err, "Node: apiPostMemberForObject: failed to GetUserByWallet")
			api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
			return
		}
		userID = user.UserID
	} else {
		userID = *inBody.UserID
	}

	id := entry.NewUserObjectID(userID, objectID)

	var modifyFunc modify.Fn[entry.UserObjectValue]
	modifyFunc = func(v *entry.UserObjectValue) (*entry.UserObjectValue, error) {
		if v == nil {
			v = &entry.UserObjectValue{}
		}
		(*v)["role"] = inBody.Role

		return v, nil
	}

	_, err = n.userObjects.Upsert(id, modifyFunc, true)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiPostMemberForObject: failed to upsert user_object")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary Delete member from object
// @Schemes
// @Description Delete member from object
// @Tags members
// @Accept json
// @Produce json
// @Param objectID path string true "ObjectID UMID"
// @Param userID path string true "UserID UMID"
// @Success 200 {object} int
// @Failure 400 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /api/v4/objects/{objectID}/members/{:userID} [delete]
func (n *Node) apiDeleteMemberFromObject(c *gin.Context) {

	objectID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiDeleteMemberFromObject: failed to parse object umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	userID, err := umid.Parse(c.Param("userID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiDeleteMemberFromObject: failed to parse object umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	id := entry.NewUserObjectID(userID, objectID)
	_, err = n.userObjects.Remove(id, true)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiDeleteMemberFromObject: failed to remove user_object")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}

	c.JSON(http.StatusOK, nil)
}
