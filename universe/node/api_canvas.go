package node

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/dto"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

// @Summary Gets user contributions by object id
// @Description Returns an object with a nested items array of contributions
// @Tags canvas
// @Security Bearer
// @Param object_id path string true "ObjectID string"
// @Success 200 {object} node.apiCanvasGetUserContributions.Out
// @Failure 400 {object} api.HTTPError
// @Router /api/v4/canvas/{object_id}/user-contributions [get]
func (n *Node) apiCanvasGetUserContributions(c *gin.Context) {
	type InQuery struct {
		OrderBy string `form:"order" json:"order"`
		Limit   uint   `form:"limit,default=10" json:"limit"`
		Offset  uint   `form:"offset" json:"offset"`
		Search  string `form:"q" json:"q"`
	}
	var inQuery InQuery
	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Node: apiCanvasGetUserContributions: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_query", err, n.log)
	}
	var limit uint
	if inQuery.Limit > 100 { // TODO: go 1.21 max function
		limit = 100
	} else {
		limit = inQuery.Limit
	}

	parentID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiCanvasGetUserContributions: failed to parse object umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	parent, ok := n.GetObjectFromAllObjects(parentID)
	if !ok {
		err := errors.Errorf("Node: apiCanvasGetUserContributions: parent object not found: %s", parentID)
		api.AbortRequest(c, http.StatusNotFound, "object_not_found", err, n.log)
		return
	}

	childrenIDs := parent.GetChildIDs()

	attrNames := []string{universe.ReservedAttributes.Object.CanvasContribution.Name}
	ouaDB := n.db.GetObjectUserAttributesDB()
	canvasContributionObjectUserAttributes, err := ouaDB.GetObjectUserAttributesByObjectIDsAttributeIDs(c, attrNames, childrenIDs, limit)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiCanvasGetUserContributions: failed to get canvasContributionObjectUserAttributes")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_attributes", err, n.log)
		return
	}

	userCanvasContributions := make([]dto.UserCanvasContributionItem, 0, len(canvasContributionObjectUserAttributes))

	for _, oua := range canvasContributionObjectUserAttributes {
		user, err := n.LoadUser(oua.UserID)
		if err != nil {
			err := errors.WithMessage(err, "Node: apiCanvasGetUserContributions: failed to get load user by id")
			api.AbortRequest(c, http.StatusInternalServerError, "failed_to_load_user", err, n.log)
			return
		}

		profile := user.GetProfile()

		var name string
		if profile != nil && profile.Name != nil {
			name = *profile.Name
		}

		voteObjectUserAttributesCount, err := ouaDB.GetObjectUserAttributesCountByObjectIDNullable(c, oua.ObjectID, universe.ReservedAttributes.Object.Vote.Name)
		if err != nil {
			err := errors.WithMessage(err, "Node: apiCanvasGetUserContributions: failed to get vote count")
			api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_count", err, n.log)
			return
		}

		commentObjectUserAttributesCount, err := ouaDB.GetObjectUserAttributesCountByObjectIDNullable(c, oua.ObjectID, universe.ReservedAttributes.Object.Comment.Name)
		if err != nil {
			err := errors.WithMessage(err, "Node: apiCanvasGetUserContributions: failed to get comment count")
			api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_count", err, n.log)
			return
		}

		userCanvasContribution := dto.UserCanvasContributionItem{
			ObjectID: oua.ObjectID,
			User: dto.User{
				ID:   umid.UMID{},
				Name: name,
				Profile: dto.Profile{
					Bio:         profile.Bio,
					AvatarHash:  profile.AvatarHash,
					ProfileLink: profile.ProfileLink,
				},
			},
			Type:      oua.AttributeID,
			Value:     oua.Value,
			Votes:     voteObjectUserAttributesCount,
			Comments:  commentObjectUserAttributesCount,
			CreatedAt: oua.CreatedAt,
			UpdatedAt: oua.UpdatedAt,
		}

		userCanvasContributions = append(userCanvasContributions, userCanvasContribution)
	}

	out := dto.UserCanvasContributions{
		Count:  uint(len(userCanvasContributions)),
		Limit:  limit,
		Offset: 0,
		Items:  userCanvasContributions,
	}

	c.JSON(http.StatusOK, out)
}
