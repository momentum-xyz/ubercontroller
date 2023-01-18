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

// @Summary Get space attribute
// @Schemes
// @Description Returns space attribute
// @Tags spaces
// @Accept json
// @Produce json
// @Param space_id path string true "Space ID"
// @Param plugin_id path string true "Plugin ID"
// @Param attribute_name path string true "Attribute Name"
// @Success 200 {object} entry.AttributeValue
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces/{space_id}/attributes/{plugin_id}/{attribute_name} [get]
func (n *Node) apiGetSpaceAttributesValue(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Param("spaceID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceAttributesValue: failed to parse space id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_space_id", err, n.log)
		return
	}

	pluginID, err := uuid.Parse(c.Param("pluginID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceAttributesValue: failed to parse plugin id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}

	attributeName := c.Param("attributeName")
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceAttributesValue: failed to get attribute name from path parameters")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_attribute_name", err, n.log)
		return
	}

	space, ok := n.GetSpaceFromAllSpaces(spaceID)
	if !ok {
		err := errors.Errorf("Node: apiGetSpaceAttributesValue: space not found: %s", spaceID)
		api.AbortRequest(c, http.StatusNotFound, "space_not_found", err, n.log)
		return
	}

	attributeID := entry.NewAttributeID(pluginID, attributeName)
	out, ok := space.GetSpaceAttributes().GetValue(attributeID)
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
// @Param plugin_id path string true "Plugin ID"
// @Param attribute_name path string true "Attribute Name"
// @Success 200 {object} dto.SpaceAttributeValues
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces/{space_id}/attributes-with-children/{plugin_id}/{attribute_name} [get]
func (n *Node) apiGetSpaceWithChildrenAttributeValues(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Param("spaceID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceWithChildrenAttributeValues: failed to parse space id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_space_id", err, n.log)
		return
	}

	pluginID, err := uuid.Parse(c.Param("pluginID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceWithChildrenAttributeValues: failed to parse plugin id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}

	attributeName := c.Param("attributeName")
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceWithChildrenAttributeValues: failed to get attribute name from path parameters")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_attribute_name", err, n.log)
		return
	}

	rootSpace, ok := n.GetSpaceFromAllSpaces(spaceID)
	if !ok {
		err := errors.Errorf("Node: apiGetSpaceWithChildrenAttributeValues: space not found: %s", spaceID)
		api.AbortRequest(c, http.StatusNotFound, "space_not_found", err, n.log)
		return
	}

	spaces := rootSpace.GetSpaces(true)
	spaces[rootSpace.GetID()] = rootSpace

	attributeID := entry.NewAttributeID(pluginID, attributeName)
	spaceAttributes := make(dto.SpaceAttributeValues, len(spaces))
	for _, space := range spaces {
		attributeValue, ok := space.GetSpaceAttributes().GetValue(attributeID)
		if !ok || attributeValue == nil {
			continue
		}

		spaceAttributes[space.GetID()] = attributeValue
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
// @Param plugin_id path string true "Plugin ID"
// @Param attribute_name path string true "Attribute Name"
// @Param body body node.apiSetSpaceAttributesValue.InBody true "body params"
// @Success 202 {object} entry.AttributeValue
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces/{space_id}/attributes/{plugin_id}/{attribute_name} [post]
func (n *Node) apiSetSpaceAttributesValue(c *gin.Context) {
	type InBody struct {
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

	pluginID, err := uuid.Parse(c.Param("pluginID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetSpaceAttributesValue: failed to parse plugin id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}

	attributeName := c.Param("attributeName")
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceAttributesValue: failed to get attribute name from path parameters")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_attribute_name", err, n.log)
		return
	}

	space, ok := n.GetSpaceFromAllSpaces(spaceID)
	if !ok {
		err := errors.Errorf("Node: apiSetSpaceAttributesValue: space not found: %s", spaceID)
		api.AbortRequest(c, http.StatusNotFound, "space_not_found", err, n.log)
		return
	}

	attributeID := entry.NewAttributeID(pluginID, attributeName)

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

	payload, err := space.GetSpaceAttributes().Upsert(attributeID, modifyFn, true)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiSetSpaceAttributesValue: failed to upsert space attribute")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_upsert", err, n.log)
		return
	}

	c.JSON(http.StatusAccepted, payload.Value)
}

// @Summary Get space sub attributes
// @Schemes
// @Description Returns space sub attributes
// @Tags spaces
// @Accept json
// @Produce json
// @Param space_id path string true "Space ID"
// @Param plugin_id path string true "Plugin ID"
// @Param attribute_name path string true "Attribute Name"
// @Param sub_attribute_key path string true "Sub Attribute Key"
// @Success 200 {object} dto.SpaceSubAttributes
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces/{space_id}/attributes/{plugin_id}/{attribute_name}/sub/{sub_attribute_key} [get]
func (n *Node) apiGetSpaceAttributeSubValue(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Param("spaceID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceAttributeSubValue: failed to parse space id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_space_id", err, n.log)
		return
	}

	pluginID, err := uuid.Parse(c.Param("pluginID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceAttributeSubValue: failed to parse plugin id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}

	attributeName := c.Param("attributeName")
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceAttributeSubValue: failed to get attribute name from path parameters")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_attribute_name", err, n.log)
		return
	}

	subAttributeKey := c.Param("subAttributeKey")
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceAttributeSubValue: failed to get sub-attribute key from path parameters")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_sub_attribute_key", err, n.log)
		return
	}

	space, ok := n.GetSpaceFromAllSpaces(spaceID)
	if !ok {
		err := errors.Errorf("Node: apiGetSpaceAttributeSubValue: space not found: %s", spaceID)
		api.AbortRequest(c, http.StatusNotFound, "space_not_found", err, n.log)
		return
	}

	attributeID := entry.NewAttributeID(pluginID, attributeName)
	attributeValue, ok := space.GetSpaceAttributes().GetValue(attributeID)
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
		subAttributeKey: (*attributeValue)[subAttributeKey],
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
// @Param plugin_id path string true "Plugin ID"
// @Param attribute_name path string true "Attribute Name"
// @Param sub_attribute_key path string true "Sub Attribute Key"
// @Param body body node.apiSetSpaceAttributeSubValue.Body true "body params"
// @Success 202 {object} dto.SpaceSubAttributes
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces/{space_id}/attributes/{plugin_id}/{attribute_name}/sub/{sub_attribute_key} [post]
func (n *Node) apiSetSpaceAttributeSubValue(c *gin.Context) {
	type Body struct {
		SubAttributeValue any `json:"sub_attribute_value" binding:"required"`
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

	pluginID, err := uuid.Parse(c.Param("pluginID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetSpaceAttributeSubValue: failed to parse plugin id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}

	attributeName := c.Param("attributeName")
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetSpaceAttributeSubValue: failed to get attribute name from path parameters")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_attribute_name", err, n.log)
		return
	}

	subAttributeKey := c.Param("subAttributeKey")
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetSpaceAttributeSubValue: failed to get sub-attribute key from path parameters")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_sub_attribute_key", err, n.log)
		return
	}

	space, ok := n.GetSpaceFromAllSpaces(spaceID)
	if !ok {
		err := errors.Errorf("Node: apiSetSpaceAttributeSubValue: space not found: %s", spaceID)
		api.AbortRequest(c, http.StatusNotFound, "space_not_found", err, n.log)
		return
	}

	attributeID := entry.NewAttributeID(pluginID, attributeName)

	modifyFn := func(current *entry.AttributePayload) (*entry.AttributePayload, error) {
		newValue := func() *entry.AttributeValue {
			value := entry.NewAttributeValue()
			(*value)[subAttributeKey] = inBody.SubAttributeValue
			return value
		}

		if current == nil {
			return entry.NewAttributePayload(newValue(), nil), nil
		}

		if current.Value == nil {
			current.Value = newValue()
			return current, nil
		}

		(*current.Value)[subAttributeKey] = inBody.SubAttributeValue

		return current, nil
	}

	payload, err := space.GetSpaceAttributes().Upsert(attributeID, modifyFn, true)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiSetSpaceAttributeSubValue: failed to upsert space attribute")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_upsert", err, n.log)
		return
	}

	out := dto.SpaceSubAttributes{
		subAttributeKey: (*payload.Value)[subAttributeKey],
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
// @Param plugin_id path string true "Plugin ID"
// @Param attribute_name path string true "Attribute Name"
// @Param body body node.apiRemoveSpaceAttributeSubValue.Body true "body params"
// @Success 200 {object} nil
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces/{space_id}/attributes/{plugin_id}/{attribute_name}/sub/{sub_attribute_key} [delete]
func (n *Node) apiRemoveSpaceAttributeSubValue(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Param("spaceID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveSpaceAttributeSubValue: failed to parse space id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_space_id", err, n.log)
		return
	}

	pluginID, err := uuid.Parse(c.Param("pluginID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveSpaceAttributeSubValue: failed to parse plugin id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}

	attributeName := c.Param("attributeName")
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveSpaceAttributeSubValue: failed to get attribute name from path parameters")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_attribute_name", err, n.log)
		return
	}

	subAttributeKey := c.Param("subAttributeKey")
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveSpaceAttributeSubValue: failed to get sub-attribute key from path parameters")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_sub_attribute_key", err, n.log)
		return
	}

	space, ok := n.GetSpaceFromAllSpaces(spaceID)
	if !ok {
		err := errors.Errorf("Node: apiRemoveSpaceAttributeSubValue: space not found: %s", spaceID)
		api.AbortRequest(c, http.StatusNotFound, "space_not_found", err, n.log)
		return
	}

	attributeID := entry.NewAttributeID(pluginID, attributeName)

	modifyFn := func(current *entry.AttributeValue) (*entry.AttributeValue, error) {
		if current == nil {
			return current, nil
		}

		delete(*current, subAttributeKey)

		return current, nil
	}

	if _, err := space.GetSpaceAttributes().UpdateValue(attributeID, modifyFn, true); err != nil {
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
// @Param plugin_id path string true "Plugin ID"
// @Param attribute_name path string true "Attribute Namw"
// @Param body body node.apiRemoveSpaceAttributeValue.Body true "body params"
// @Success 200 {object} nil
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces/{space_id}/attributes/{plugin_id}/{attribute_name} [delete]
func (n *Node) apiRemoveSpaceAttributeValue(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Param("spaceID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveSpaceAttributeValue: failed to parse space id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_space_id", err, n.log)
		return
	}

	pluginID, err := uuid.Parse(c.Param("pluginID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveSpaceAttributeValue: failed to parse plugin id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}

	attributeName := c.Param("attributeName")
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveSpaceAttributeValue: failed to get attribute name from path parameters")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_attribute_name", err, n.log)
		return
	}

	space, ok := n.GetSpaceFromAllSpaces(spaceID)
	if !ok {
		err := errors.Errorf("Node: apiRemoveSpaceAttributeValue: space not found: %s", spaceID)
		api.AbortRequest(c, http.StatusNotFound, "space_not_found", err, n.log)
		return
	}

	attributeID := entry.NewAttributeID(pluginID, attributeName)
	if _, err := space.GetSpaceAttributes().UpdateValue(
		attributeID, modify.ReplaceWith[entry.AttributeValue](nil), true,
	); err != nil {
		err = errors.WithMessage(err, "Node: apiRemoveSpaceAttributeValue: failed to update space attribute")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_update", err, n.log)
		return
	}

	c.JSON(http.StatusOK, nil)
}
