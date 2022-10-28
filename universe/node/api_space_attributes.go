package node

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe/api"
)

// @Summary Returns space attributes
// @Schemes
// @Description Returns space attributes
// @Tags spaces
// @Accept json
// @Produce json
// @Param space_id path string true "Space ID"
// @Success 200 {object} entry.AttributeValue
// @Success 500 {object} api.HTTPError
// @Success 400 {object} api.HTTPError
// @Success 404 {object} api.HTTPError
// @Router /api/v4/spaces/{space_id}/attributes [get]
func (n *Node) apiGetSpaceAttributes(c *gin.Context) {
	inQuery := struct {
		PluginID string `form:"plugin_id" binding:"required"`
		Name     string `form:"attribute_name" binding:"required"`
	}{}

	spaceID, err := uuid.Parse(c.Param("spaceID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceAttributes: failed to parse space id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_space_id", err, n.log)
		return
	}

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceAttributes: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, n.log)
		return
	}

	pluginID, err := uuid.Parse(inQuery.PluginID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceAttributes: failed to parse plugin id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}

	space, ok := n.GetSpaceFromAllSpaces(spaceID)
	if !ok {
		err := errors.Errorf("Node: apiGetSpaceAttributes: space not found: %s", spaceID)
		api.AbortRequest(c, http.StatusNotFound, "space_not_found", err, n.log)
		return
	}

	attributeID := entry.NewAttributeID(pluginID, inQuery.Name)
	out, ok := space.GetSpaceAttributeValue(attributeID)
	if !ok {
		err := errors.Errorf("Node: apiGetSpaceSubAttribute: space attribute value not found: %s", attributeID)
		api.AbortRequest(c, http.StatusNotFound, "attribute_not_found", err, n.log)
		return
	}

	c.JSON(http.StatusOK, out)
}

func (n *Node) apiGetSpaceSubAttribute(c *gin.Context) {
	inQuery := struct {
		PluginID        string `form:"plugin_id" binding:"required"`
		Name            string `form:"attribute_name" binding:"required"`
		SubAttributeKey string `form:"sub_attribute_key" binding:"required"`
	}{}

	spaceID, err := uuid.Parse(c.Param("spaceID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceSubAttributes: failed to parse space id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_space_id", err, n.log)
		return
	}

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceSubAttribute: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, n.log)
		return
	}

	pluginID, err := uuid.Parse(inQuery.PluginID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceSubAttributes: failed to parse plugin id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}

	space, ok := n.GetSpaceFromAllSpaces(spaceID)
	if !ok {
		err := errors.Errorf("Node: apiGetSpaceSubAttribute: space not found: %s", spaceID)
		api.AbortRequest(c, http.StatusNotFound, "space_not_found", err, n.log)
		return
	}

	attributeID := entry.NewAttributeID(pluginID, inQuery.Name)
	attributeValue, ok := space.GetSpaceAttributeValue(attributeID)
	if !ok {
		err := errors.Errorf("Node: apiGetSpaceSubAttribute: attribute value not found: %s", attributeID)
		api.AbortRequest(c, http.StatusNotFound, "attribute_value_not_found", err, n.log)
		return
	}

	if attributeValue == nil {
		err := errors.Errorf("Node: apiGetSpaceSubAttribute: attribute value is nil")
		api.AbortRequest(c, http.StatusNotFound, "attribute_value_nil", err, n.log)
		return
	}

	value, ok := (*attributeValue)[inQuery.SubAttributeKey]
	if !ok {
		err := errors.Errorf("Node: apiGetSpaceSubAttribute: attribute key not found: %s", inQuery.SubAttributeKey)
		api.AbortRequest(c, http.StatusNotFound, "attribute_key_not_found", err, n.log)
		return
	}

	c.JSON(http.StatusOK, value)
}
