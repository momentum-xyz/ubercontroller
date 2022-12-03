package node

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe/common/api"
)

// @Summary Get list of attributes for all users limited by space, plugin and attribute_name
// @Schemes
// @Description Returns map with key as userID and value as Attribute Value
// @Tags spaces
// @Accept json
// @Produce json
// @Param space_id path string true "Space ID"
// @Param query query node.apiGetSpaceAllUsersAttributeValuesList.InQuery true "query params"
// @Success 200 {object} map[uuid.UUID]entry.AttributeValue
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces/{space_id}/all-users/attributes [get]
func (n *Node) apiGetSpaceAllUsersAttributeValuesList(c *gin.Context) {
	type InQuery struct {
		PluginID      string `form:"plugin_id" binding:"required"`
		AttributeName string `form:"attribute_name" binding:"required"`
	}

	inQuery := InQuery{}

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceAllUsersAttributeValuesList: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, n.log)
		return
	}

	spaceID, err := uuid.Parse(c.Param("spaceID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceAllUsersAttributeValuesList: failed to parse space id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_space_id", err, n.log)
		return
	}

	pluginID, err := uuid.Parse(inQuery.PluginID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceUserAttributesValue: failed to parse plugin id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}

	sua, err := n.db.SpaceUserAttributesGetSpaceUserAttributesBySpaceID(context.Background(), spaceID)
	if err != nil {
		panic(err)
	}

	out := make(map[uuid.UUID]entry.AttributeValue)
	for i := range sua {
		if sua[i] == nil {
			break
		}
		item := *sua[i]

		if item.PluginID != pluginID {
			break
		}

		if item.AttributeID.Name != inQuery.AttributeName {
			break
		}

		value := item.AttributePayload.Value
		if value != nil {
			out[item.UserID] = *value
		}
	}

	c.JSON(http.StatusOK, out)
}
