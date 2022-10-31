package node

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe/api"
	"github.com/momentum-xyz/ubercontroller/universe/api/dto"
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
func (n *Node) apiGetSpaceAttributesValue(c *gin.Context) {
	inQuery := struct {
		PluginID      string `form:"plugin_id" binding:"required"`
		AttributeName string `form:"attribute_name" binding:"required"`
	}{}

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceAttributesValue: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, n.log)
		return
	}

	spaceID, err := uuid.Parse(c.Param("spaceID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceAttributesValue: failed to parse space id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_space_id", err, n.log)
		return
	}

	pluginID, err := uuid.Parse(inQuery.PluginID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceAttributesValue: failed to parse plugin id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}

	space, ok := n.GetSpaceFromAllSpaces(spaceID)
	if !ok {
		err := errors.Errorf("Node: apiGetSpaceAttributesValue: space not found: %s", spaceID)
		api.AbortRequest(c, http.StatusNotFound, "space_not_found", err, n.log)
		return
	}

	attributeID := entry.NewAttributeID(pluginID, inQuery.AttributeName)
	out, ok := space.GetSpaceAttributeValue(attributeID)
	if !ok {
		err := errors.Errorf("Node: apiGetSpaceAttributeSubValue: space attribute value not found: %s", attributeID)
		api.AbortRequest(c, http.StatusNotFound, "attribute_not_found", err, n.log)
		return
	}

	c.JSON(http.StatusOK, out)
}

func (n *Node) apiGetSpaceAttributeSubValue(c *gin.Context) {
	inQuery := struct {
		PluginID        string `form:"plugin_id" binding:"required"`
		AttributeName   string `form:"attribute_name" binding:"required"`
		SubAttributeKey string `form:"sub_attribute_key" binding:"required"`
	}{}

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceAttributeSubValue: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, n.log)
		return
	}

	spaceID, err := uuid.Parse(c.Param("spaceID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceSubAttributes: failed to parse space id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_space_id", err, n.log)
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
		err := errors.Errorf("Node: apiGetSpaceAttributeSubValue: space not found: %s", spaceID)
		api.AbortRequest(c, http.StatusNotFound, "space_not_found", err, n.log)
		return
	}

	attributeID := entry.NewAttributeID(pluginID, inQuery.AttributeName)
	attributeValue, ok := space.GetSpaceAttributeValue(attributeID)
	if !ok {
		err := errors.Errorf("Node: apiGetSpaceAttributeSubValue: attribute value not found: %s", attributeID)
		api.AbortRequest(c, http.StatusNotFound, "attribute_value_not_found", err, n.log)
		return
	}

	if attributeValue == nil {
		err := errors.Errorf("Node: apiGetSpaceAttributeSubValue: attribute value is nil")
		api.AbortRequest(c, http.StatusNotFound, "attribute_value_nil", err, n.log)
		return
	}

	out := dto.SpaceSubAttributes{
		inQuery.SubAttributeKey: (*attributeValue)[inQuery.SubAttributeKey],
	}

	c.JSON(http.StatusOK, out)
}

func (n *Node) apiSetSpaceAttributeSubValue(c *gin.Context) {
	inBody := struct {
		PluginID          string `json:"plugin_id" binding:"required"`
		AttributeName     string `json:"attribute_name" binding:"required"`
		SubAttributeKey   string `json:"sub_attribute_key" binding:"required"`
		SubAttributeValue string `json:"sub_attribute_value" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiSetSpaceAttributeSubValue: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	spaceID, err := uuid.Parse(c.Param("spaceID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetSpaceAttributeSubValue: failed to parse space id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_space_id", err, n.log)
		return
	}

	pluginID, err := uuid.Parse(inBody.PluginID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetSpaceAttributeSubValue: failed to parse plugin id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}

	space, ok := n.GetSpaceFromAllSpaces(spaceID)
	if !ok {
		err := errors.Errorf("Node: apiSetSpaceAttributeSubValue: space not found: %s", spaceID)
		api.AbortRequest(c, http.StatusNotFound, "space_not_found", err, n.log)
		return
	}

	attributeID := entry.NewAttributeID(pluginID, inBody.AttributeName)

	modifyFn := func(current *entry.AttributePayload) (*entry.AttributePayload, error) {
		newValue := func() *entry.AttributeValue {
			value := entry.NewAttributeValue()
			(*value)[inBody.SubAttributeKey] = inBody.SubAttributeValue
			return value
		}

		if current == nil {
			return entry.NewAttributePayload(newValue(), nil), nil
		}

		if current.Value == nil {
			current.Value = newValue()
			return current, nil
		}

		(*current.Value)[inBody.SubAttributeKey] = inBody.SubAttributeValue

		return current, nil
	}

	spaceAttribute, err := space.UpsertSpaceAttribute(attributeID, modifyFn, true)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiSetSpaceAttributeSubValue: failed to upsert space attribute")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_upsert", err, n.log)
		return
	}

	out := dto.SpaceSubAttributes{
		inBody.SubAttributeKey: (*spaceAttribute.Value)[inBody.SubAttributeKey],
	}

	c.JSON(http.StatusCreated, out)
}

func (n *Node) apiRemoveSpaceAttributeSubValue(c *gin.Context) {
	inQuery := struct {
		PluginID        string `form:"plugin_id" binding:"required"`
		AttributeName   string `form:"attribute_name" binding:"required"`
		SubAttributeKey string `form:"sub_attribute_key" binding:"required"`
	}{}

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveSpaceAttributeSubValue: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, n.log)
		return
	}

	spaceID, err := uuid.Parse(c.Param("spaceID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveSpaceAttributeSubValue: failed to parse space id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_space_id", err, n.log)
		return
	}

	pluginID, err := uuid.Parse(inQuery.PluginID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveSpaceAttributeSubValue: failed to parse plugin id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}

	space, ok := n.GetSpaceFromAllSpaces(spaceID)
	if !ok {
		err := errors.Errorf("Node: apiRemoveSpaceAttributeSubValue: space not found: %s", spaceID)
		api.AbortRequest(c, http.StatusNotFound, "space_not_found", err, n.log)
		return
	}

	attributeID := entry.NewAttributeID(pluginID, inQuery.AttributeName)

	modifyFn := func(current *entry.AttributeValue) (*entry.AttributeValue, error) {
		if current == nil {
			return current, nil
		}

		delete(*current, inQuery.SubAttributeKey)

		return current, nil
	}

	if err := space.UpdateSpaceAttributeValue(attributeID, modifyFn, true); err != nil {
		err = errors.WithMessage(err, "Node: apiRemoveSpaceAttributeSubValue: failed to remove space attribute")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_upsert", err, n.log)
		return
	}

	c.JSON(http.StatusOK, nil)
}
