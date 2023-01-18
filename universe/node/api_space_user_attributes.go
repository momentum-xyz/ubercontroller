package node

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe/common/api"
	"github.com/momentum-xyz/ubercontroller/universe/common/api/dto"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

// @Summary Get space user attribute
// @Schemes
// @Description Returns space user attribute
// @Tags spaces
// @Accept json
// @Produce json
// @Param space_id path string true "Space ID"
// @Param user_id path string true "User ID"
// @Param plugin_id path string true "Plugin ID"
// @Param attribute_name path string true "Attribute Name"
// @Success 200 {object} entry.AttributeValue
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces/{space_id}/{user_id}/attributes/{plugin_id}/{attribute_name} [get]
func (n *Node) apiGetSpaceUserAttributesValue(c *gin.Context) {
	pluginID, err := uuid.Parse(c.Param("pluginID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceUserAttributesValue: failed to parse plugin id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}

	attributeName := c.Param("attributeName")
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceUserAttributesValue: failed to get attribute name from path parameters")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_attribute_name", err, n.log)
		return
	}

	spaceID, err := uuid.Parse(c.Param("spaceID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceUserAttributesValue: failed to parse space id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_space_id", err, n.log)
		return
	}

	userID, err := uuid.Parse(c.Param("userID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetSpaceUserAttributesValue: failed to parse user id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	attributeID := entry.NewAttributeID(pluginID, attributeName)
	spaceUserAttributeID := entry.NewSpaceUserAttributeID(attributeID, spaceID, userID)
	out, ok := n.GetSpaceUserAttributes().GetValue(spaceUserAttributeID)
	if !ok {
		err := errors.Errorf("Node: apiGetSpaceUserAttributesValue: space attribute value not found: %s", attributeID)
		api.AbortRequest(c, http.StatusNotFound, "attribute_not_found", err, n.log)
		return
	}

	c.JSON(http.StatusOK, out)
}

// @Summary Set space user attribute
// @Schemes
// @Description Sets entire space user attribute
// @Tags spaces
// @Accept json
// @Produce json
// @Param space_id path string true "Space ID"
// @Param user_id path string true "User ID"
// @Param plugin_id path string true "Plugin ID"
// @Param attribute_name path string true "Attribute Name"
// @Param body body node.apiSetSpaceAttributesValue.InBody true "body params"
// @Success 202 {object} entry.AttributeValue
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces/{space_id}/{user_id}/attributes/{plugin_id}/{attribute_name} [post]
func (n *Node) apiSetSpaceUserAttributesValue(c *gin.Context) {
	type InBody struct {
		AttributeValue map[string]any `json:"attribute_value" binding:"required"`
	}

	inBody := InBody{}

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiSetSpaceUserAttributesValue: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	pluginID, err := uuid.Parse(c.Param("pluginID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetSpaceUserAttributesValue: failed to parse plugin id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}

	attributeName := c.Param("attributeName")
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetSpaceUserAttributesValue: failed to get attribute name from path parameters")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_attribute_name", err, n.log)
		return
	}

	spaceID, err := uuid.Parse(c.Param("spaceID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetSpaceUserAttributesValue: failed to parse space id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_space_id", err, n.log)
		return
	}

	userID, err := uuid.Parse(c.Param("userID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetSpaceUserAttributesValue: failed to parse user id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	attributeID := entry.NewAttributeID(pluginID, attributeName)
	spaceUserAttributeID := entry.NewSpaceUserAttributeID(attributeID, spaceID, userID)

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

	spaceUserAttribute, err := n.GetSpaceUserAttributes().Upsert(spaceUserAttributeID, modifyFn, true)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiSetSpaceUserAttributesValue: failed to upsert space user attribute")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_upsert", err, n.log)
		return
	}

	c.JSON(http.StatusAccepted, spaceUserAttribute.Value)
}

// @Summary Get space sub attributes
// @Schemes
// @Description Returns space sub attributes
// @Tags spaces
// @Accept json
// @Produce json
// @Param space_id path string true "Space ID"
// @Param user_id path string true "User ID"
// @Param plugin_id path string true "Plugin ID"
// @Param attribute_name path string true "Attribute Name"
// @Param sub_attribute_key path string true "Sub Attribute Key"
// @Success 200 {object} dto.SpaceSubAttributes
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces/{space_id}/{user_id}/attributes/{plugin_id}/{attribute_name}/sub/{sub_attribute_key} [get]
func (n *Node) apiGetSpaceUserAttributeSubValue(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Param("spaceID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceUserAttributeSubValue: failed to parse space id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_space_id", err, n.log)
		return
	}

	userID, err := uuid.Parse(c.Param("userID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceUserAttributeSubValue: failed to parse user id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	pluginID, err := uuid.Parse(c.Param("pluginID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceUserAttributeSubValue: failed to parse plugin id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}

	attributeName := c.Param("attributeName")
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceUserAttributeSubValue: failed to get attribute name from path parameters")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_attribute_name", err, n.log)
		return
	}

	subAttributeKey := c.Param("subAttributeKey")
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceUserAttributeSubValue: failed to get sub-attribute key from path parameters")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_sub_attribute_key", err, n.log)
		return
	}

	attributeID := entry.NewAttributeID(pluginID, attributeName)
	spaceUserAttributeID := entry.NewSpaceUserAttributeID(attributeID, spaceID, userID)
	spaceUserAttributeValue, ok := n.GetSpaceUserAttributes().GetValue(spaceUserAttributeID)
	if !ok {
		err := errors.Errorf("Node: apiGetSpaceUserAttributeSubValue: space user attribute value not found: %s", attributeID)
		api.AbortRequest(c, http.StatusNotFound, "attribute_value_not_found", err, n.log)
		return
	}

	if spaceUserAttributeValue == nil {
		err := errors.Errorf("Node: apiGetSpaceUserAttributeSubValue: attribute value is nil")
		api.AbortRequest(c, http.StatusNotFound, "attribute_value_nil", err, n.log)
		return
	}

	out := dto.SpaceSubAttributes{
		subAttributeKey: (*spaceUserAttributeValue)[subAttributeKey],
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
// @Param user_id path string true "User ID"
// @Param plugin_id path string true "Plugin ID"
// @Param attribute_name path string true "Attribute Name"
// @Param sub_attribute_key path string true "Sub Attribute Key"
// @Param body body node.apiSetSpaceAttributeSubValue.Body true "body params"
// @Success 202 {object} dto.SpaceSubAttributes
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces/{space_id}/{user_id}/attributes/{plugin_id}/{attribute_name}/sub/{sub_attribute_key} [post]
func (n *Node) apiSetSpaceUserAttributeSubValue(c *gin.Context) {
	type Body struct {
		SubAttributeValue any `json:"sub_attribute_value" binding:"required"`
	}

	inBody := Body{}

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiSetSpaceUserAttributeSubValue: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	spaceID, err := uuid.Parse(c.Param("spaceID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetSpaceUserAttributeSubValue: failed to parse space id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_space_id", err, n.log)
		return
	}

	userID, err := uuid.Parse(c.Param("userID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetSpaceUserAttributeSubValue: failed to parse user id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	pluginID, err := uuid.Parse(c.Param("pluginID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetSpaceUserAttributeSubValue: failed to parse plugin id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}

	attributeName := c.Param("attributeName")
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetSpaceUserAttributeSubValue: failed to get attribute name from path parameters")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_attribute_name", err, n.log)
		return
	}

	subAttributeKey := c.Param("subAttributeKey")
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetSpaceUserAttributeSubValue: failed to get sub-attribute key from path parameters")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_sub_attribute_key", err, n.log)
		return
	}

	attributeID := entry.NewAttributeID(pluginID, attributeName)
	spaceUserAttributeID := entry.NewSpaceUserAttributeID(attributeID, spaceID, userID)

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

	spaceUserAttribute, err := n.GetSpaceUserAttributes().Upsert(spaceUserAttributeID, modifyFn, true)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiSetSpaceUserAttributeSubValue: failed to upsert space user attribute")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_upsert", err, n.log)
		return
	}

	out := dto.SpaceSubAttributes{
		subAttributeKey: (*spaceUserAttribute.Value)[subAttributeKey],
	}

	c.JSON(http.StatusAccepted, out)
}

// @Summary Delete space user sub attribute
// @Schemes
// @Description Deletes a space user sub attribute
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
// @Router /api/v4/spaces/{space_id}/{user_id}/attributes/{plugin_id}/{attribute_name}/sub/{sub_attribute_key} [delete]
func (n *Node) apiRemoveSpaceUserAttributeSubValue(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Param("spaceID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveSpaceUserAttributeSubValue: failed to parse space id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_space_id", err, n.log)
		return
	}

	userID, err := uuid.Parse(c.Param("userID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveSpaceUserAttributeSubValue: failed to parse user id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	pluginID, err := uuid.Parse(c.Param("pluginID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveSpaceUserAttributeSubValue: failed to parse plugin id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}

	attributeName := c.Param("attributeName")
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveSpaceUserAttributeSubValue: failed to get attribute name from path parameters")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_attribute_name", err, n.log)
		return
	}

	subAttributeKey := c.Param("subAttributeKey")
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveSpaceUserAttributeSubValue: failed to get sub-attribute key from path parameters")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_sub_attribute_key", err, n.log)
		return
	}

	attributeID := entry.NewAttributeID(pluginID, attributeName)
	spaceUserAttributeID := entry.NewSpaceUserAttributeID(attributeID, spaceID, userID)

	modifyFn := func(current *entry.AttributeValue) (*entry.AttributeValue, error) {
		if current == nil {
			return current, nil
		}

		delete(*current, subAttributeKey)

		return current, nil
	}

	if _, err := n.GetSpaceUserAttributes().UpdateValue(spaceUserAttributeID, modifyFn, true); err != nil {
		err = errors.WithMessage(err, "Node: apiRemoveSpaceUserAttributeSubValue: failed to update space user attribute")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_update", err, n.log)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary Delete space user attribute
// @Schemes
// @Description Deletes a space attribute
// @Tags spaces
// @Accept json
// @Produce json
// @Param space_id path string true "Space ID"
// @Param user_id path string true "User ID"
// @Param plugin_id path string true "Plugin ID"
// @Param attribute_name path string true "Attribute Name"
// @Param body body node.apiRemoveSpaceAttributeValue.Body true "body params"
// @Success 200 {object} nil
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces/{space_id}/{user_id}/attributes/{plugin_id}/{attribute_name} [delete]
func (n *Node) apiRemoveSpaceUserAttributeValue(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Param("spaceID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveSpaceUserAttributeValue: failed to parse space id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_space_id", err, n.log)
		return
	}

	userID, err := uuid.Parse(c.Param("userID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveSpaceUserAttributeValue: failed to parse user id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	pluginID, err := uuid.Parse(c.Param("pluginID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveSpaceUserAttributeValue: failed to parse plugin id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}

	attributeName := c.Param("attributeName")
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveSpaceUserAttributeValue: failed to get attribute name from path parameters")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_attribute_name", err, n.log)
		return
	}

	attributeID := entry.NewAttributeID(pluginID, attributeName)
	spaceUserAttributeID := entry.NewSpaceUserAttributeID(attributeID, spaceID, userID)

	if _, err := n.GetSpaceUserAttributes().UpdateValue(
		spaceUserAttributeID, modify.ReplaceWith[entry.AttributeValue](nil), true,
	); err != nil {
		err = errors.WithMessage(err, "Node: apiRemoveSpaceUserAttributeValue: failed to update space user attribute")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_update", err, n.log)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary Get list of attributes for all users limited by space, plugin and attribute_name
// @Schemes
// @Description Returns map with key as userID and value as Attribute Value
// @Tags spaces
// @Accept json
// @Produce json
// @Param space_id path string true "Space ID"
// @Param plugin_id path string true "Plugin ID"
// @Param attribute_name path string true "Attribute Name"
// @Success 200 {object} map[uuid.UUID]entry.AttributeValue
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces/{space_id}/all-users/attributes/{plugin_id}/{attribute_name} [get]
func (n *Node) apiGetSpaceAllUsersAttributeValuesList(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Param("spaceID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceAllUsersAttributeValuesList: failed to parse space id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_space_id", err, n.log)
		return
	}

	pluginID, err := uuid.Parse(c.Param("pluginID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceAllUsersAttributeValuesList: failed to parse plugin id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}

	attributeName := c.Param("attributeName")
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceAllUsersAttributeValuesList: failed to get attribute name from path parameters")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_attribute_name", err, n.log)
		return
	}

	sua, err := n.db.SpaceUserAttributesGetSpaceUserAttributesByPluginIDAndNameAndSpaceID(
		n.ctx, pluginID, attributeName, spaceID,
	)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceAllUsersAttributeValuesList: failed to get space user attributes")
		api.AbortRequest(c, http.StatusInternalServerError, "get_space_user_attributes_failed", err, n.log)
		return
	}

	out := make(map[uuid.UUID]*entry.AttributeValue)
	for i := range sua {
		if sua[i].AttributePayload == nil || sua[i].AttributePayload.Value == nil {
			continue
		}
		out[sua[i].UserID] = sua[i].AttributePayload.Value
	}

	c.JSON(http.StatusOK, out)
}
