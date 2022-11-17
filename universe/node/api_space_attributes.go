package node

import (
	"github.com/momentum-xyz/ubercontroller/utils/modify"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe/common/api"
	"github.com/momentum-xyz/ubercontroller/universe/common/api/dto"
)

// @Summary Returns space attributes
// @Schemes
// @Description Returns space attributes
// @Tags spaces
// @Accept json
// @Produce json
// @Param space_id path string true "Space ID"
// @Param plugin_id query string true "Plugin ID"
// @Param attribute_name query string true "Attribute Name"
// @Success 200 {object} entry.AttributeValue
// @Success 500 {object} api.HTTPError
// @Success 400 {object} api.HTTPError
// @Success 404 {object} api.HTTPError
// @Router /api/v4/spaces/{space_id}/attributes [get]
func (n *Node) apiGetSpaceAttributesValue(c *gin.Context) {
	type Query struct {
		PluginID      string `form:"plugin_id" binding:"required"`
		AttributeName string `form:"attribute_name" binding:"required"`
	}

	inQuery := Query{}

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

// @Summary Returns space and all subspace attributes
// @Schemes
// @Description Returns space and all subspace attributes
// @Tags spaces
// @Accept json
// @Produce json
// @Param space_id path string true "Space ID"
// @Param plugin_id query string true "Plugin ID"
// @Param attribute_name query string true "Attribute Name"
// @Success 200 {object} dto.SpaceAttributeValues
// @Success 500 {object} api.HTTPError
// @Success 400 {object} api.HTTPError
// @Success 404 {object} api.HTTPError
// @Router /api/v4/spaces/{space_id}/attributes-with-children [get]
func (n *Node) apiGetSpaceWithChildrenAttributeValues(c *gin.Context) {
	type Query struct {
		PluginID      string `form:"plugin_id" binding:"required"`
		AttributeName string `form:"attribute_name" binding:"required"`
	}

	inQuery := Query{}

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceWithChildrenAttributeValues: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, n.log)
		return
	}

	spaceID, err := uuid.Parse(c.Param("spaceID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceWithChildrenAttributeValues: failed to parse space id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_space_id", err, n.log)
		return
	}

	pluginID, err := uuid.Parse(inQuery.PluginID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceWithChildrenAttributeValues: failed to parse plugin id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}

	rootSpace, ok := n.GetSpaceFromAllSpaces(spaceID)
	if !ok {
		err := errors.Errorf("Node: apiGetSpaceWithChildrenAttributeValues: space not found: %s", spaceID)
		api.AbortRequest(c, http.StatusNotFound, "space_not_found", err, n.log)
		return
	}

	spaces := rootSpace.GetSpaces(true)
	spaceAttributes := make(dto.SpaceAttributeValues, len(spaces))

	spaces[rootSpace.GetID()] = rootSpace

	for _, space := range spaces {
		attributeID := entry.NewAttributeID(pluginID, inQuery.AttributeName)
		attributeValue, ok := space.GetSpaceAttributeValue(attributeID)
		if !ok {
			continue
		}

		if attributeValue != nil {
			spaceAttributes[space.GetID()] = attributeValue
		}
	}

	c.JSON(http.StatusOK, spaceAttributes)
}

// @Summary Updates entire space attribute value
// @Schemes
// @Description Updates entire space attribute value
// @Tags spaces
// @Accept json
// @Produce json
// @Param space_id path string true "Space ID"
// @Param plugin_id body string true "Plugin ID"
// @Param attribute_name body string true "Attribute Name"
// @Param attribute_value body []string true "[Attribute Value]"
// @Success 202 {object} entry.AttributeValue
// @Success 500 {object} api.HTTPError
// @Success 400 {object} api.HTTPError
// @Success 404 {object} api.HTTPError
// @Router /api/v4/spaces/{space_id}/attributes [post]
func (n *Node) apiSetSpaceAttributesValue(c *gin.Context) {
	type Body struct {
		PluginID       string         `json:"plugin_id" binding:"required"`
		AttributeName  string         `json:"attribute_name" binding:"required"`
		AttributeValue map[string]any `json:"attribute_value" binding:"required"`
	}

	inBody := Body{}

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiSetSpaceAttributesValue: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	spaceID, err := uuid.Parse(c.Param("spaceID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetSpaceAttributesValue: failed to parse space id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_space_id", err, n.log)
		return
	}

	pluginID, err := uuid.Parse(inBody.PluginID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetSpaceAttributesValue: failed to parse plugin id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}

	space, ok := n.GetSpaceFromAllSpaces(spaceID)
	if !ok {
		err := errors.Errorf("Node: apiSetSpaceAttributesValue: space not found: %s", spaceID)
		api.AbortRequest(c, http.StatusNotFound, "space_not_found", err, n.log)
		return
	}

	attributeID := entry.NewAttributeID(pluginID, inBody.AttributeName)

	modifyFn := func(current *entry.AttributePayload) (*entry.AttributePayload, error) {
		newValue := func() *entry.AttributeValue {
			value := entry.NewAttributeValue()
			*value = inBody.AttributeValue
			return value
		}

		if current == nil {
			return entry.NewAttributePayload(newValue(), nil), nil
		}

		if current.Value == nil {
			current.Value = newValue()
			return current, nil
		}

		*current.Value = inBody.AttributeValue

		return current, nil
	}

	spaceAttribute, err := space.UpsertSpaceAttribute(attributeID, modifyFn, true)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiSetSpaceAttributesValue: failed to upsert space attribute")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_upsert", err, n.log)
		return
	}

	c.JSON(http.StatusAccepted, spaceAttribute.Value)
}

// @Summary Returns space attributes sub value
// @Schemes
// @Description Returns space attributes sub value
// @Tags spaces
// @Accept json
// @Produce json
// @Param space_id path string true "Space ID"
// @Param plugin_id query string true "Plugin ID"
// @Param attribute_name query string true "Name"
// @Param sub_attribute_key query string true "Sub Attribute Key"
// @Success 200 {object} dto.SpaceSubAttributes
// @Success 500 {object} api.HTTPError
// @Success 400 {object} api.HTTPError
// @Success 404 {object} api.HTTPError
// @Router /api/v4/spaces/{space_id}/attributes/sub [get]
func (n *Node) apiGetSpaceAttributeSubValue(c *gin.Context) {
	type Query struct {
		PluginID        string `form:"plugin_id" binding:"required"`
		AttributeName   string `form:"attribute_name" binding:"required"`
		SubAttributeKey string `form:"sub_attribute_key" binding:"required"`
	}

	inQuery := Query{}

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

// @Summary Adds a space attribute sub value
// @Schemes
// @Description Adds a space attribute sub value
// @Tags spaces
// @Accept json
// @Produce json
// @Param space_id path string true "Space ID"
// @Param plugin_id body string true "Plugin ID"
// @Param attribute_name body string true "Name"
// @Param sub_attribute_key body string true "Sub Attribute Key"
// @Param sub_attribute_value body string true "Sub Attribute Value"
// @Success 202 {object} dto.SpaceSubAttributes
// @Success 500 {object} api.HTTPError
// @Success 400 {object} api.HTTPError
// @Success 404 {object} api.HTTPError
// @Router /api/v4/spaces/{space_id}/attributes/sub [post]
func (n *Node) apiSetSpaceAttributeSubValue(c *gin.Context) {
	type Body struct {
		PluginID          string `json:"plugin_id" binding:"required"`
		AttributeName     string `json:"attribute_name" binding:"required"`
		SubAttributeKey   string `json:"sub_attribute_key" binding:"required"`
		SubAttributeValue any    `json:"sub_attribute_value" binding:"required"`
	}

	inBody := Body{}

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

	c.JSON(http.StatusAccepted, out)
}

// @Summary Deletes a space attribute sub value
// @Schemes
// @Description Deletes a space attribute sub value
// @Tags spaces
// @Accept json
// @Produce json
// @Param space_id path string true "Space ID"
// @Param plugin_id body string true "Plugin ID"
// @Param attribute_name body string true "Name"
// @Param sub_attribute_key body string true "Sub Attribute Key"
// @Success 200 {object} nil
// @Success 500 {object} api.HTTPError
// @Success 400 {object} api.HTTPError
// @Success 404 {object} api.HTTPError
// @Router /api/v4/spaces/{space_id}/attributes/sub [delete]
func (n *Node) apiRemoveSpaceAttributeSubValue(c *gin.Context) {
	type Body struct {
		PluginID        string `json:"plugin_id" binding:"required"`
		AttributeName   string `json:"attribute_name" binding:"required"`
		SubAttributeKey string `json:"sub_attribute_key" binding:"required"`
	}

	inBody := Body{}

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiRemoveSpaceAttributeSubValue: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	spaceID, err := uuid.Parse(c.Param("spaceID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveSpaceAttributeSubValue: failed to parse space id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_space_id", err, n.log)
		return
	}

	pluginID, err := uuid.Parse(inBody.PluginID)
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

	attributeID := entry.NewAttributeID(pluginID, inBody.AttributeName)

	modifyFn := func(current *entry.AttributeValue) (*entry.AttributeValue, error) {
		if current == nil {
			return current, nil
		}

		delete(*current, inBody.SubAttributeKey)

		return current, nil
	}

	if _, err := space.UpdateSpaceAttributeValue(attributeID, modifyFn, true); err != nil {
		err = errors.WithMessage(err, "Node: apiRemoveSpaceAttributeSubValue: failed to update space attribute")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_update", err, n.log)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (n *Node) apiRemoveSpaceAttributeValue(c *gin.Context) {
	type Body struct {
		PluginID      string `json:"plugin_id" binding:"required"`
		AttributeName string `json:"attribute_name" binding:"required"`
	}

	var inBody Body
	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiRemoveSpaceAttributeValue: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	spaceID, err := uuid.Parse(c.Param("spaceID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveSpaceAttributeValue: failed to parse space id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_space_id", err, n.log)
		return
	}

	pluginID, err := uuid.Parse(inBody.PluginID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveSpaceAttributeValue: failed to parse plugin id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}

	space, ok := n.GetSpaceFromAllSpaces(spaceID)
	if !ok {
		err := errors.Errorf("Node: apiRemoveSpaceAttributeValue: space not found: %s", spaceID)
		api.AbortRequest(c, http.StatusNotFound, "space_not_found", err, n.log)
		return
	}

	attributeID := entry.NewAttributeID(pluginID, inBody.AttributeName)
	if _, err := space.UpdateSpaceAttributeValue(
		attributeID, modify.ReplaceWith[entry.AttributeValue](nil), true,
	); err != nil {
		err = errors.WithMessage(err, "Node: apiRemoveSpaceAttributeValue: failed to update space attribute")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_update", err, n.log)
		return
	}
	
	c.JSON(http.StatusOK, nil)
}
