package node

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/dto"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

// @Summary Get object user attribute
// @Schemes
// @Description Returns object user attribute
// @Tags spaces
// @Accept json
// @Produce json
// @Param object_id path string true "Object ID"
// @Param query query node.apiGetSpaceAttributesValue.InQuery true "query params"
// @Success 200 {object} entry.AttributeValue
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces/{object_id}/{user_id}/attributes [get]
func (n *Node) apiGetSpaceUserAttributesValue(c *gin.Context) {
	type InQuery struct {
		PluginID      string `form:"plugin_id" binding:"required"`
		AttributeName string `form:"attribute_name" binding:"required"`
	}

	inQuery := InQuery{}

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceUserAttributesValue: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, n.log)
		return
	}

	objectID, err := uuid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceUserAttributesValue: failed to parse object id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	userID, err := uuid.Parse(c.Param("userID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetSpaceUserAttributesValue: failed to parse user id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	pluginID, err := uuid.Parse(inQuery.PluginID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceUserAttributesValue: failed to parse plugin id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}

	attributeID := entry.NewAttributeID(pluginID, inQuery.AttributeName)
	spaceUserAttributeID := entry.NewObjectUserAttributeID(attributeID, objectID, userID)
	out, ok := n.GetObjectUserAttributes().GetValue(spaceUserAttributeID)
	if !ok {
		err := errors.Errorf("Node: apiGetSpaceUserAttributesValue: object attribute value not found: %s", attributeID)
		api.AbortRequest(c, http.StatusNotFound, "attribute_not_found", err, n.log)
		return
	}

	c.JSON(http.StatusOK, out)
}

// @Summary Set object user attribute
// @Schemes
// @Description Sets entire object user attribute
// @Tags spaces
// @Accept json
// @Produce json
// @Param object_id path string true "Object ID"
// @Param body body node.apiSetSpaceAttributesValue.InBody true "body params"
// @Success 202 {object} entry.AttributeValue
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces/{object_id}/{user_id}/attributes [post]
func (n *Node) apiSetSpaceUserAttributesValue(c *gin.Context) {
	type InBody struct {
		PluginID       string         `json:"plugin_id" binding:"required"`
		AttributeName  string         `json:"attribute_name" binding:"required"`
		AttributeValue map[string]any `json:"attribute_value" binding:"required"`
	}

	inBody := InBody{}

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiSetSpaceUserAttributesValue: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	objectID, err := uuid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetSpaceUserAttributesValue: failed to parse object id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	userID, err := uuid.Parse(c.Param("userID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetSpaceUserAttributesValue: failed to parse user id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	pluginID, err := uuid.Parse(inBody.PluginID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetSpaceUserAttributesValue: failed to parse plugin id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}

	attributeID := entry.NewAttributeID(pluginID, inBody.AttributeName)
	spaceUserAttributeID := entry.NewObjectUserAttributeID(attributeID, objectID, userID)

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

	spaceUserAttribute, err := n.GetObjectUserAttributes().Upsert(spaceUserAttributeID, modifyFn, true)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiSetSpaceUserAttributesValue: failed to upsert object user attribute")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_upsert", err, n.log)
		return
	}

	c.JSON(http.StatusAccepted, spaceUserAttribute.Value)
}

// @Summary Get object sub attributes
// @Schemes
// @Description Returns object sub attributes
// @Tags spaces
// @Accept json
// @Produce json
// @Param object_id path string true "Object ID"
// @Param query query node.apiGetSpaceAttributeSubValue.InQuery true "query params"
// @Success 200 {object} dto.SpaceSubAttributes
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces/{object_id}/{user_id}/attributes/sub [get]
func (n *Node) apiGetSpaceUserAttributeSubValue(c *gin.Context) {
	type InQuery struct {
		PluginID        string `form:"plugin_id" binding:"required"`
		AttributeName   string `form:"attribute_name" binding:"required"`
		SubAttributeKey string `form:"sub_attribute_key" binding:"required"`
	}

	inQuery := InQuery{}

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceUserAttributeSubValue: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, n.log)
		return
	}

	objectID, err := uuid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceUserAttributeSubValue: failed to parse object id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	userID, err := uuid.Parse(c.Param("userID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceUserAttributeSubValue: failed to parse user id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	pluginID, err := uuid.Parse(inQuery.PluginID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceUserAttributeSubValue: failed to parse plugin id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}

	attributeID := entry.NewAttributeID(pluginID, inQuery.AttributeName)
	spaceUserAttributeID := entry.NewObjectUserAttributeID(attributeID, objectID, userID)
	spaceUserAttributeValue, ok := n.GetObjectUserAttributes().GetValue(spaceUserAttributeID)
	if !ok {
		err := errors.Errorf("Node: apiGetSpaceUserAttributeSubValue: object user attribute value not found: %s", attributeID)
		api.AbortRequest(c, http.StatusNotFound, "attribute_value_not_found", err, n.log)
		return
	}

	if spaceUserAttributeValue == nil {
		err := errors.Errorf("Node: apiGetSpaceUserAttributeSubValue: attribute value is nil")
		api.AbortRequest(c, http.StatusNotFound, "attribute_value_nil", err, n.log)
		return
	}

	out := dto.SpaceSubAttributes{
		inQuery.SubAttributeKey: (*spaceUserAttributeValue)[inQuery.SubAttributeKey],
	}

	c.JSON(http.StatusOK, out)
}

// @Summary Set object sub attribute
// @Schemes
// @Description Sets a object sub attribute
// @Tags spaces
// @Accept json
// @Produce json
// @Param object_id path string true "Object ID"
// @Param body body node.apiSetSpaceAttributeSubValue.Body true "body params"
// @Success 202 {object} dto.SpaceSubAttributes
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces/{object_id}/{user_id}/attributes/sub [post]
func (n *Node) apiSetSpaceUserAttributeSubValue(c *gin.Context) {
	type Body struct {
		PluginID          string `json:"plugin_id" binding:"required"`
		AttributeName     string `json:"attribute_name" binding:"required"`
		SubAttributeKey   string `json:"sub_attribute_key" binding:"required"`
		SubAttributeValue any    `json:"sub_attribute_value" binding:"required"`
	}

	inBody := Body{}

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiSetSpaceUserAttributeSubValue: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	objectID, err := uuid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetSpaceUserAttributeSubValue: failed to parse object id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	userID, err := uuid.Parse(c.Param("userID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetSpaceUserAttributeSubValue: failed to parse user id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	pluginID, err := uuid.Parse(inBody.PluginID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetSpaceUserAttributeSubValue: failed to parse plugin id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}

	attributeID := entry.NewAttributeID(pluginID, inBody.AttributeName)
	spaceUserAttributeID := entry.NewObjectUserAttributeID(attributeID, objectID, userID)

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

	spaceUserAttribute, err := n.GetObjectUserAttributes().Upsert(spaceUserAttributeID, modifyFn, true)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiSetSpaceUserAttributeSubValue: failed to upsert object user attribute")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_upsert", err, n.log)
		return
	}

	out := dto.SpaceSubAttributes{
		inBody.SubAttributeKey: (*spaceUserAttribute.Value)[inBody.SubAttributeKey],
	}

	c.JSON(http.StatusAccepted, out)
}

// @Summary Delete object user sub attribute
// @Schemes
// @Description Deletes a object user sub attribute
// @Tags spaces
// @Accept json
// @Produce json
// @Param object_id path string true "Object ID"
// @Param body body node.apiRemoveSpaceAttributeSubValue.Body true "body params"
// @Success 200 {object} nil
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces/{object_id}/{user_id}/attributes/sub [delete]
func (n *Node) apiRemoveSpaceUserAttributeSubValue(c *gin.Context) {
	type Body struct {
		PluginID        string `json:"plugin_id" binding:"required"`
		AttributeName   string `json:"attribute_name" binding:"required"`
		SubAttributeKey string `json:"sub_attribute_key" binding:"required"`
	}

	inBody := Body{}

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiRemoveSpaceUserAttributeSubValue: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	objectID, err := uuid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveSpaceUserAttributeSubValue: failed to parse object id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	userID, err := uuid.Parse(c.Param("userID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetSpaceUserAttributeSubValue: failed to parse user id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	pluginID, err := uuid.Parse(inBody.PluginID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveSpaceUserAttributeSubValue: failed to parse plugin id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}

	attributeID := entry.NewAttributeID(pluginID, inBody.AttributeName)
	spaceUserAttributeID := entry.NewObjectUserAttributeID(attributeID, objectID, userID)

	modifyFn := func(current *entry.AttributeValue) (*entry.AttributeValue, error) {
		if current == nil {
			return current, nil
		}

		delete(*current, inBody.SubAttributeKey)

		return current, nil
	}

	if _, err := n.GetObjectUserAttributes().UpdateValue(spaceUserAttributeID, modifyFn, true); err != nil {
		err = errors.WithMessage(err, "Node: apiRemoveSpaceUserAttributeSubValue: failed to update object user attribute")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_update", err, n.log)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary Delete object user attribute
// @Schemes
// @Description Deletes a object attribute
// @Tags spaces
// @Accept json
// @Produce json
// @Param object_id path string true "Object ID"
// @Param body body node.apiRemoveSpaceAttributeValue.Body true "body params"
// @Success 200 {object} nil
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces/{object_id}/{user_id}/attributes [delete]
func (n *Node) apiRemoveSpaceUserAttributeValue(c *gin.Context) {
	type Body struct {
		PluginID      string `json:"plugin_id" binding:"required"`
		AttributeName string `json:"attribute_name" binding:"required"`
	}

	var inBody Body
	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiRemoveSpaceUserAttributeValue: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	objectID, err := uuid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveSpaceUserAttributeValue: failed to parse object id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	userID, err := uuid.Parse(c.Param("userID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetSpaceUserAttributeSubValue: failed to parse user id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	pluginID, err := uuid.Parse(inBody.PluginID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveSpaceUserAttributeValue: failed to parse plugin id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}

	attributeID := entry.NewAttributeID(pluginID, inBody.AttributeName)
	spaceUserAttributeID := entry.NewObjectUserAttributeID(attributeID, objectID, userID)

	if _, err := n.GetObjectUserAttributes().UpdateValue(
		spaceUserAttributeID, modify.ReplaceWith[entry.AttributeValue](nil), true,
	); err != nil {
		err = errors.WithMessage(err, "Node: apiRemoveSpaceUserAttributeValue: failed to update object user attribute")
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
// @Param object_id path string true "Object ID"
// @Param query query node.apiGetSpaceAllUsersAttributeValuesList.InQuery true "query params"
// @Success 200 {object} map[uuid.UUID]entry.AttributeValue
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces/{object_id}/all-users/attributes [get]
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

	objectID, err := uuid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceAllUsersAttributeValuesList: failed to parse object id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	pluginID, err := uuid.Parse(inQuery.PluginID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceAllUsersAttributeValuesList: failed to parse plugin id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}

	sua, err := n.db.GetObjectUserAttributesDB().GetObjectUserAttributesByObjectAttributeID(
		n.ctx, entry.NewObjectAttributeID(entry.NewAttributeID(pluginID, inQuery.AttributeName), objectID),
	)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpaceAllUsersAttributeValuesList: failed to get object user attributes")
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
