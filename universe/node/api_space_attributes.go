package node

import (
	"net/http"

	"github.com/momentum-xyz/ubercontroller/utils/modify"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe/common/api"
	"github.com/momentum-xyz/ubercontroller/universe/common/api/dto"
)

// @Summary Get count of subkeys in space attribute
// @Schemes
// @Description Returns total number ob subkeys in given space and attribute_name
// @Tags spaces
// @Accept json
// @Produce json
// @Param space_id path string true "Space ID"
// @Param query query node.apiGetSpaceAttributesSubsCount.InQuery true "query params"
// @Success 200 {object} node.apiGetSpaceAttributesSubsCount.Out
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces/{space_id}/attributes/subs-count [get]
func (n *Node) apiGetSpaceAttributesSubsCount(c *gin.Context) {
	type InQuery struct {
		PluginID      string `form:"plugin_id" binding:"required"`
		AttributeName string `form:"attribute_name" binding:"required"`
	}

	inQuery := InQuery{}

	type Out struct {
		Count int `json:"count"`
	}

	out := Out{
		Count: 0,
	}

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
	attrValue, ok := space.GetSpaceAttributeValue(attributeID)
	if !ok {
		err := errors.Errorf("Node: apiGetSpaceAttributesValue: space attribute value not found: %s", attributeID)
		api.AbortRequest(c, http.StatusNotFound, "attribute_not_found", err, n.log)
		return
	}

	if attrValue != nil {
		o := *attrValue
		for _ = range o {
			out.Count++
		}
	}

	c.JSON(http.StatusOK, out)
}

// @Summary Get space attribute
// @Schemes
// @Description Returns space attribute
// @Tags spaces
// @Accept json
// @Produce json
// @Param space_id path string true "Space ID"
// @Param query query node.apiGetSpaceAttributesValue.InQuery true "query params"
// @Success 200 {object} entry.AttributeValue
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces/{space_id}/attributes [get]
func (n *Node) apiGetSpaceAttributesValue(c *gin.Context) {
	type InQuery struct {
		PluginID      string `form:"plugin_id" binding:"required"`
		AttributeName string `form:"attribute_name" binding:"required"`
	}

	inQuery := InQuery{}

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
		err := errors.Errorf("Node: apiGetSpaceAttributesValue: space attribute value not found: %s", attributeID)
		api.AbortRequest(c, http.StatusNotFound, "attribute_not_found", err, n.log)
		return
	}

	c.JSON(http.StatusOK, out)
}

// @Summary Get space and all subspace attributes
// @Schemes
// @Description Returns space and all subspace attributes
// @Tags spaces
// @Accept json
// @Produce json
// @Param space_id path string true "Space ID"
// @Param query query node.apiGetSpaceWithChildrenAttributeValues.InQuery true "query params"
// @Success 200 {object} dto.SpaceAttributeValues
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces/{space_id}/attributes-with-children [get]
func (n *Node) apiGetSpaceWithChildrenAttributeValues(c *gin.Context) {
	type InQuery struct {
		PluginID      string `form:"plugin_id" binding:"required"`
		AttributeName string `form:"attribute_name" binding:"required"`
	}

	inQuery := InQuery{}

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

// @Summary Set space attribute
// @Schemes
// @Description Sets entire space attribute
// @Tags spaces
// @Accept json
// @Produce json
// @Param space_id path string true "Space ID"
// @Param body body node.apiSetSpaceAttributesValue.InBody true "body params"
// @Success 202 {object} entry.AttributeValue
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces/{space_id}/attributes [post]
func (n *Node) apiSetSpaceAttributesValue(c *gin.Context) {
	type InBody struct {
		PluginID       string         `json:"plugin_id" binding:"required"`
		AttributeName  string         `json:"attribute_name" binding:"required"`
		AttributeValue map[string]any `json:"attribute_value" binding:"required"`
	}

	inBody := InBody{}

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

// @Summary Get space sub attributes
// @Schemes
// @Description Returns space sub attributes
// @Tags spaces
// @Accept json
// @Produce json
// @Param space_id path string true "Space ID"
// @Param query query node.apiGetSpaceAttributeSubValue.InQuery true "query params"
// @Success 200 {object} dto.SpaceSubAttributes
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces/{space_id}/attributes/sub [get]
func (n *Node) apiGetSpaceAttributeSubValue(c *gin.Context) {
	type InQuery struct {
		PluginID        string `form:"plugin_id" binding:"required"`
		AttributeName   string `form:"attribute_name" binding:"required"`
		SubAttributeKey string `form:"sub_attribute_key" binding:"required"`
	}

	inQuery := InQuery{}

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

// @Summary Set space sub attribute
// @Schemes
// @Description Sets a space sub attribute
// @Tags spaces
// @Accept json
// @Produce json
// @Param space_id path string true "Space ID"
// @Param body body node.apiSetSpaceAttributeSubValue.Body true "body params"
// @Success 202 {object} dto.SpaceSubAttributes
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
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

// @Summary Delete space sub attribute
// @Schemes
// @Description Deletes a space sub attribute
// @Tags spaces
// @Accept json
// @Produce json
// @Param space_id path string true "Space ID"
// @Param body body node.apiRemoveSpaceAttributeSubValue.Body true "body params"
// @Success 200 {object} nil
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
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

// @Summary Delete space attribute
// @Schemes
// @Description Deletes a space attribute
// @Tags spaces
// @Accept json
// @Produce json
// @Param space_id path string true "Space ID"
// @Param body body node.apiRemoveSpaceAttributeValue.Body true "body params"
// @Success 200 {object} nil
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces/{space_id}/attributes [delete]
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
