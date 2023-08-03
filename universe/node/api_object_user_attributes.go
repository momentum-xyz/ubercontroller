package node

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe/attributes"
	"github.com/momentum-xyz/ubercontroller/universe/auth"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/dto"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

// @Summary Get object user attribute
// @Schemes
// @Description Returns object user attribute
// @Tags objects
// @Accept json
// @Produce json
// @Param object_id path string true "Object UMID"
// @Param user_id path string true "User UMID"
// @Param query query attributes.QueryPluginAttribute true "query params"
// @Success 200 {object} entry.AttributeValue
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/objects/{object_id}/{user_id}/attributes [get]
func (n *Node) apiGetObjectUserAttributesValue(c *gin.Context) {
	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "user from context")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user", err, n.log)
		return
	}

	objectID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetObjectUserAttributesValue: failed to parse object umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	targetUserID, err := umid.Parse(c.Param("userID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetObjectUserAttributesValue: failed to parse user umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	attrType, attributeID, err := attributes.PluginAttributeFromQuery(c, n)
	if err != nil {
		err := fmt.Errorf("failed to get plugin attribute from query: %w", err)
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_attribute", err, n.log)
		return
	}

	objectUserAttributeID := entry.NewObjectUserAttributeID(attributeID, objectID, targetUserID)

	allowed, err := auth.CheckAttributePermissions(
		c, *attrType.GetEntry(), n.GetObjectUserAttributes(), objectUserAttributeID, userID,
		auth.ReadOperation)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetObjectUserAttributesValue: permissions check")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_permissions_check", err, n.log)
		return
	} else if !allowed {
		err := fmt.Errorf("operation not permitted")
		api.AbortRequest(c, http.StatusForbidden, "operation_not_permitted", err, n.log)
		return
	}

	out, ok := n.GetObjectUserAttributes().GetValue(objectUserAttributeID)
	if !ok {
		err := errors.Errorf("Node: apiGetObjectUserAttributesValue: object attribute value not found: %s", attributeID)
		api.AbortRequest(c, http.StatusNotFound, "attribute_not_found", err, n.log)
		return
	}

	c.JSON(http.StatusOK, out)
}

// @Summary Set object user attribute
// @Schemes
// @Description Sets entire object user attribute
// @Tags objects
// @Accept json
// @Produce json
// @Param object_id path string true "Object UMID"
// @Param user_id path string true "User UMID"
// @Param body body node.apiSetObjectUserAttributesValue.InBody true "body params"
// @Success 202 {object} entry.AttributeValue
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/objects/{object_id}/{user_id}/attributes [post]
func (n *Node) apiSetObjectUserAttributesValue(c *gin.Context) {
	type InBody struct {
		attributes.QueryPluginAttribute
		AttributeValue map[string]any `json:"attribute_value" binding:"required"`
	}

	inBody := InBody{}

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiSetObjectUserAttributesValue: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	objectID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetObjectUserAttributesValue: failed to parse object umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	targetUserID, err := umid.Parse(c.Param("userID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetObjectUserAttributesValue: failed to parse user umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	pluginID, err := umid.Parse(inBody.PluginID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetObjectUserAttributesValue: failed to parse plugin umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}
	attrType, ok := n.GetAttributeTypes().GetAttributeType(entry.AttributeTypeID{pluginID, inBody.AttributeName})
	if !ok {
		err := fmt.Errorf("attribute type not found")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_attribute", err, n.log)
		return
	}
	attributeID := entry.NewAttributeID(pluginID, inBody.AttributeName)
	objectUserAttributeID := entry.NewObjectUserAttributeID(attributeID, objectID, targetUserID)
	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "user from context")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user", err, n.log)
		return
	}

	allowed, err := auth.CheckAttributePermissions(
		c, *attrType.GetEntry(), n.GetObjectUserAttributes(), objectUserAttributeID, userID,
		auth.WriteOperation)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetObjectUserAttributesValue: permissions check")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_permissions_check", err, n.log)
		return
	} else if !allowed {
		err := fmt.Errorf("operation not permitted")
		api.AbortRequest(c, http.StatusForbidden, "operation_not_permitted", err, n.log)
		return
	}

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

	objectUserAttribute, err := n.GetObjectUserAttributes().Upsert(objectUserAttributeID, modifyFn, true)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiSetObjectUserAttributesValue: failed to upsert object user attribute")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_upsert", err, n.log)
		return
	}

	c.JSON(http.StatusAccepted, objectUserAttribute.Value)
}

// @Summary Get object user sub attribute
// @Schemes
// @Description Returns object user sub attributes
// @Tags objects
// @Accept json
// @Produce json
// @Param object_id path string true "Object UMID"
// @Param user_id path string true "User UMID"
// @Param query query node.apiGetObjectUserAttributeSubValue.InQuery true "query params"
// @Success 200 {object} dto.ObjectSubAttributes
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/objects/{object_id}/{user_id}/attributes/sub [get]
func (n *Node) apiGetObjectUserAttributeSubValue(c *gin.Context) {
	type InQuery struct {
		attributes.QueryPluginAttribute
		SubAttributeKey string `form:"sub_attribute_key" binding:"required"`
	}

	inQuery := InQuery{}

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Node: apiGetObjectUserAttributeSubValue: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, n.log)
		return
	}

	objectID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetObjectUserAttributeSubValue: failed to parse object umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	targetUserID, err := umid.Parse(c.Param("userID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetObjectUserAttributeSubValue: failed to parse user umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	pluginID, err := umid.Parse(inQuery.PluginID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetObjectUserAttributeSubValue: failed to parse plugin umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}
	attrType, ok := n.GetAttributeTypes().GetAttributeType(entry.AttributeTypeID{pluginID, inQuery.AttributeName})
	if !ok {
		err := fmt.Errorf("attribute type not found")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_attribute", err, n.log)
		return
	}
	attributeID := entry.NewAttributeID(pluginID, inQuery.AttributeName)
	objectUserAttributeID := entry.NewObjectUserAttributeID(attributeID, objectID, targetUserID)
	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "user from context")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user", err, n.log)
		return
	}

	allowed, err := auth.CheckAttributePermissions(
		c, *attrType.GetEntry(), n.GetObjectUserAttributes(), objectUserAttributeID, userID,
		auth.ReadOperation)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetObjectUserAttributesValue: permissions check")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_permissions_check", err, n.log)
		return
	} else if !allowed {
		err := fmt.Errorf("operation not permitted")
		api.AbortRequest(c, http.StatusForbidden, "operation_not_permitted", err, n.log)
		return
	}

	objectUserAttributeValue, ok := n.GetObjectUserAttributes().GetValue(objectUserAttributeID)
	if !ok {
		err := errors.Errorf(
			"Node: apiGetObjectUserAttributeSubValue: object user attribute value not found: %s", attributeID,
		)
		api.AbortRequest(c, http.StatusNotFound, "attribute_value_not_found", err, n.log)
		return
	}

	if objectUserAttributeValue == nil {
		err := errors.Errorf("Node: apiGetObjectUserAttributeSubValue: attribute value is nil")
		api.AbortRequest(c, http.StatusNotFound, "attribute_value_nil", err, n.log)
		return
	}

	out := dto.ObjectSubAttributes{
		inQuery.SubAttributeKey: (*objectUserAttributeValue)[inQuery.SubAttributeKey],
	}

	c.JSON(http.StatusOK, out)
}

// @Summary Set object sub attribute
// @Schemes
// @Description Sets a object sub attribute
// @Tags objects
// @Accept json
// @Produce json
// @Param object_id path string true "Object UMID"
// @Param user_id path string true "User UMID"
// @Param body body node.apiSetObjectUserAttributeSubValue.Body true "body params"
// @Success 202 {object} dto.ObjectSubAttributes
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/objects/{object_id}/{user_id}/attributes/sub [post]
func (n *Node) apiSetObjectUserAttributeSubValue(c *gin.Context) {
	type Body struct {
		attributes.QueryPluginAttribute
		SubAttributeKey   string `json:"sub_attribute_key" binding:"required"`
		SubAttributeValue any    `json:"sub_attribute_value" binding:"required"`
	}

	inBody := Body{}

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiSetObjectUserAttributeSubValue: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	objectID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetObjectUserAttributeSubValue: failed to parse object umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	targetUserID, err := umid.Parse(c.Param("userID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetObjectUserAttributeSubValue: failed to parse user umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	pluginID, err := umid.Parse(inBody.PluginID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetObjectUserAttributeSubValue: failed to parse plugin umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}
	attrType, ok := n.GetAttributeTypes().GetAttributeType(entry.AttributeTypeID{pluginID, inBody.AttributeName})
	if !ok {
		err := fmt.Errorf("attribute type not found")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_attribute", err, n.log)
		return
	}
	attributeID := entry.NewAttributeID(pluginID, inBody.AttributeName)
	objectUserAttributeID := entry.NewObjectUserAttributeID(attributeID, objectID, targetUserID)
	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "user from context")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user", err, n.log)
		return
	}

	allowed, err := auth.CheckAttributePermissions(
		c, *attrType.GetEntry(), n.GetObjectUserAttributes(), objectUserAttributeID, userID,
		auth.WriteOperation)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetObjectUserAttributesValue: permissions check")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_permissions_check", err, n.log)
		return
	} else if !allowed {
		err := fmt.Errorf("operation not permitted")
		api.AbortRequest(c, http.StatusForbidden, "operation_not_permitted", err, n.log)
		return
	}

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

	objectUserAttribute, err := n.GetObjectUserAttributes().Upsert(objectUserAttributeID, modifyFn, true)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiSetObjectUserAttributeSubValue: failed to upsert object user attribute")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_upsert", err, n.log)
		return
	}

	out := dto.ObjectSubAttributes{
		inBody.SubAttributeKey: (*objectUserAttribute.Value)[inBody.SubAttributeKey],
	}

	c.JSON(http.StatusAccepted, out)
}

// @Summary Delete object user sub attribute
// @Schemes
// @Description Deletes a object user sub attribute
// @Tags objects
// @Accept json
// @Produce json
// @Param object_id path string true "Object UMID"
// @Param user_id path string true "User UMID"
// @Param body body node.apiRemoveObjectUserAttributeSubValue.Body true "body params"
// @Success 200 {object} nil
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/objects/{object_id}/{user_id}/attributes/sub [delete]
func (n *Node) apiRemoveObjectUserAttributeSubValue(c *gin.Context) {
	type Body struct {
		attributes.QueryPluginAttribute
		SubAttributeKey string `json:"sub_attribute_key" binding:"required"`
	}

	inBody := Body{}

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiRemoveObjectUserAttributeSubValue: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	objectID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveObjectUserAttributeSubValue: failed to parse object umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	targetUserID, err := umid.Parse(c.Param("userID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetObjectUserAttributeSubValue: failed to parse user umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	pluginID, err := umid.Parse(inBody.PluginID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveObjectUserAttributeSubValue: failed to parse plugin umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}
	attrType, ok := n.GetAttributeTypes().GetAttributeType(entry.AttributeTypeID{pluginID, inBody.AttributeName})
	if !ok {
		err := fmt.Errorf("attribute type not found")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_attribute", err, n.log)
		return
	}
	attributeID := entry.NewAttributeID(pluginID, inBody.AttributeName)
	objectUserAttributeID := entry.NewObjectUserAttributeID(attributeID, objectID, targetUserID)
	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "user from context")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user", err, n.log)
		return
	}

	allowed, err := auth.CheckAttributePermissions(
		c, *attrType.GetEntry(), n.GetObjectUserAttributes(), objectUserAttributeID, userID,
		auth.WriteOperation)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetObjectUserAttributesValue: permissions check")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_permissions_check", err, n.log)
		return
	} else if !allowed {
		err := fmt.Errorf("operation not permitted")
		api.AbortRequest(c, http.StatusForbidden, "operation_not_permitted", err, n.log)
		return
	}

	modifyFn := func(current *entry.AttributeValue) (*entry.AttributeValue, error) {
		if current == nil {
			return current, nil
		}

		delete(*current, inBody.SubAttributeKey)

		return current, nil
	}

	if _, err := n.GetObjectUserAttributes().UpdateValue(objectUserAttributeID, modifyFn, true); err != nil {
		err = errors.WithMessage(
			err, "Node: apiRemoveObjectUserAttributeSubValue: failed to update object user attribute",
		)
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_update", err, n.log)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary Delete object user attribute
// @Schemes
// @Description Deletes a object attribute
// @Tags objects
// @Accept json
// @Produce json
// @Param object_id path string true "Object UMID"
// @Param user_id path string true "User UMID"
// @Param body body node.apiRemoveObjectUserAttributeValue.Body true "body params"
// @Success 200 {object} nil
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/objects/{object_id}/{user_id}/attributes [delete]
func (n *Node) apiRemoveObjectUserAttributeValue(c *gin.Context) {
	type Body struct {
		attributes.QueryPluginAttribute
	}

	var inBody Body
	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiRemoveObjectUserAttributeValue: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	objectID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveObjectUserAttributeValue: failed to parse object umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	targetUserID, err := umid.Parse(c.Param("userID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetObjectUserAttributeSubValue: failed to parse user umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	pluginID, err := umid.Parse(inBody.PluginID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveObjectUserAttributeValue: failed to parse plugin umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}
	attrType, ok := n.GetAttributeTypes().GetAttributeType(entry.AttributeTypeID{pluginID, inBody.AttributeName})
	if !ok {
		err := fmt.Errorf("attribute type not found")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_attribute", err, n.log)
		return
	}
	attributeID := entry.NewAttributeID(pluginID, inBody.AttributeName)
	objectUserAttributeID := entry.NewObjectUserAttributeID(attributeID, objectID, targetUserID)
	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "user from context")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user", err, n.log)
		return
	}

	allowed, err := auth.CheckAttributePermissions(
		c, *attrType.GetEntry(), n.GetObjectUserAttributes(), objectUserAttributeID, userID,
		auth.WriteOperation)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetObjectUserAttributesValue: permissions check")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_permissions_check", err, n.log)
		return
	} else if !allowed {
		err := fmt.Errorf("operation not permitted")
		api.AbortRequest(c, http.StatusForbidden, "operation_not_permitted", err, n.log)
		return
	}

	if _, err := n.GetObjectUserAttributes().UpdateValue(
		objectUserAttributeID, modify.ReplaceWith[entry.AttributeValue](nil), true,
	); err != nil {
		err = errors.WithMessage(err, "Node: apiRemoveObjectUserAttributeValue: failed to update object user attribute")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_update", err, n.log)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary Get list of attributes for all users limited by object, plugin and attribute_name
// @Schemes
// @Description Returns map with key as userID and value as Attribute Value
// @Tags objects
// @Accept json
// @Produce json
// @Param object_id path string true "Object UMID"
// @Param query query attributes.QueryPluginAttribute true "query params"
// @Success 200 {object} map[umid.UMID]entry.AttributeValue
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/objects/{object_id}/all-users/attributes [get]
func (n *Node) apiGetObjectAllUsersAttributeValuesList(c *gin.Context) {
	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "user from context")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user", err, n.log)
		return
	}

	objectID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetObjectAllUsersAttributeValuesList: failed to parse object umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	attrType, attributeID, err := attributes.PluginAttributeFromQuery(c, n)
	if err != nil {
		err := fmt.Errorf("failed to get plugin attribute: %w", err)
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_attribute", err, n.log)
		return
	}

	allowed, err := auth.CheckReadAllPermissions[entry.ObjectUserAttributeID](
		c, *attrType.GetEntry(), n.GetObjectUserAttributes(), userID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetObjectUserAttributesValue: permissions check")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_permissions_check", err, n.log)
		return
	} else if !allowed {
		err := fmt.Errorf("operation not permitted")
		api.AbortRequest(c, http.StatusForbidden, "operation_not_permitted", err, n.log)
		return
	}

	sua, err := n.db.GetObjectUserAttributesDB().GetObjectUserAttributesByObjectAttributeID(
		n.ctx, entry.NewObjectAttributeID(attributeID, objectID))
	if err != nil {
		err := errors.WithMessage(
			err, "Node: apiGetObjectAllUsersAttributeValuesList: failed to get object user attributes",
		)
		api.AbortRequest(c, http.StatusInternalServerError, "get_object_user_attributes_failed", err, n.log)
		return
	}

	out := make(map[umid.UMID]*entry.AttributeValue)
	for i := range sua {
		if sua[i].AttributePayload == nil || sua[i].AttributePayload.Value == nil {
			continue
		}
		out[sua[i].UserID] = sua[i].AttributePayload.Value
	}

	c.JSON(http.StatusOK, out)
}

// @Summary Get object user attribute count
// @Schemes
// @Description Returns the number of (non null) object user attributes.
// @Tags objects
// @Accept json
// @Produce json
// @Param object_id path string true "Object UMID"
// @Param query query node.apiGetObjectUserAttributeCount.InQuery true "query params"
// @Success 200 {object} dto.AttributeCount
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/objects/{object_id}/all-users/count [get]
func (n *Node) apiGetObjectUserAttributeCount(c *gin.Context) {
	type InQuery struct {
		attributes.QueryPluginAttribute
		Since string `form:"since"`
	}

	inQuery := InQuery{}

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Node: apiGetObjectUserAttributeCount: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, n.log)
		return
	}

	var sinceTime *time.Time
	if inQuery.Since != "" {
		since, err := time.Parse(time.RFC3339, inQuery.Since)
		if err != nil {
			err := errors.WithMessage(err, "Node: apiGetObjectUserAttributeCount: failed to parse 'since'")
			api.AbortRequest(c, http.StatusBadRequest, "invalid_since_time", err, n.log)
			return
		}

		sinceTime = &since
	}

	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetObjectUserAttributeCount: failed to get user from context")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user", err, n.log)
		return
	}

	objectID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetObjectUserAttributeCount: failed to parse object umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	attrType, attributeID, err := attributes.PluginAttributeFromQuery(c, n)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetObjectUserAttributeCount: failed to get plugin attribute from query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_attribute", err, n.log)
		return
	}

	objectUserAttributeID := entry.NewObjectUserAttributeID(attributeID, objectID, umid.Nil)

	allowed, err := auth.CheckAttributePermissions(
		c, *attrType.GetEntry(), n.GetObjectUserAttributes(), objectUserAttributeID, userID,
		auth.ReadOperation)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetObjectUserAttributeCount: permissions check")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_permissions_check", err, n.log)
		return
	} else if !allowed {
		err := errors.WithMessage(err, "Node: apiGetObjectUserAttributeCount: operation not permitted")
		api.AbortRequest(c, http.StatusForbidden, "operation_not_permitted", err, n.log)
		return
	}

	count, ok := n.GetObjectUserAttributes().GetCountByObjectID(objectID, inQuery.AttributeName, sinceTime)
	if !ok {
		err := errors.Errorf("Node: apiGetObjectUserAttributeCount: object attribute value not found: %s", attributeID)
		api.AbortRequest(c, http.StatusNotFound, "attribute_not_found", err, n.log)
		return
	}

	out := dto.AttributeCount{Count: count}

	c.JSON(http.StatusOK, out)
}

// @Summary Get combined list of json entries from all user's values.
// @Description Allow showing a combined (sorted) list from all users.
// @Description The attribute value is assumed to be a JSON (map-like) object, with some ID as key and the value nested JSON object.
// @Description The fields params allows selecting some fields to directly return in the list.
// @Description The limit and offset params allow pagination.
// @Description Limit defaults to 10, maximun allowed is 100.
// @Tags objects
// @Accept json
// @Produce json
// @Param object_id path string true "Object UMID"
// @Param plugin_id path string true "Plugin MID"
// @Param attribute_name path string true "Name of the plugin attribute"
// @Param query query node.apiObjectUserAttributeValueEntries.InQuery false "query params"
// @Success 202 {object} node.apiObjectUserAttributeValueEntries.JsonResult
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/objects/{object_id}/all-users/attributes/{plugin_id}/{attribute_name}/entries [get]
func (n *Node) apiObjectUserAttributeValueEntries(c *gin.Context) {
	objectID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := fmt.Errorf("invalid object ID: %w", err)
		api.AbortRequest(c, http.StatusBadRequest, "invalid_param", err, n.log)
		return
	}
	attrType, attrID, err := attributes.PluginAttributeFromURL(c, n)
	if err != nil {
		err := fmt.Errorf("plugin attribute: %w", err)
		api.AbortRequest(c, http.StatusNotFound, "invalid_param", err, n.log)
		return
	}
	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := fmt.Errorf("user from context: %w", err)
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user", err, n.log)
		return
	}
	allowed, err := auth.CheckReadAllPermissions[entry.ObjectUserAttributeID](
		c, *attrType.GetEntry(), n.GetObjectUserAttributes(), userID)
	if err != nil {
		err := fmt.Errorf("check read permissions: %w", err)
		api.AbortRequest(c, http.StatusInternalServerError, "failed_permissions_check", err, n.log)
		return
	}
	if !allowed {
		err := fmt.Errorf("operation not permitted")
		api.AbortRequest(c, http.StatusForbidden, "operation_not_permitted", err, n.log)
		return
	}

	type InQuery struct {
		Fields  []string `form:"fields[]"`
		OrderBy string   `form:"order"`
		Limit   uint     `form:"limit,default=10"`
		Offset  uint     `form:"offset"`
	}
	var q InQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		err := fmt.Errorf("bind query: %w", err)
		api.AbortRequest(c, http.StatusBadRequest, "invalid_query", err, n.log)
	}
	var limit uint
	if q.Limit > 100 { // TODO: go 1.21 max function
		limit = 100
	} else {
		limit = q.Limit
	}
	objAttrID := entry.NewObjectAttributeID(attrID, objectID)
	oua := n.db.GetObjectUserAttributesDB()
	count, err := oua.ValueEntriesCount(c, objAttrID)
	if err != nil {
		err := fmt.Errorf("count query: %w", err)
		api.AbortRequest(c, http.StatusInternalServerError, "invalid_params", err, n.log)
		return
	}
	type JsonResult struct {
		Count  uint                     `json:"count"`
		Limit  uint                     `json:"limit"`
		Offset uint                     `json:"offset"`
		Items  []map[string]interface{} `json:"items"`
	}
	if count == 0 {
		c.JSON(http.StatusOK, JsonResult{Limit: limit})
		return
	}
	order := q.OrderBy
	desc := false
	if q.OrderBy != "" {
		if strings.HasPrefix(order, "-") {
			desc = true
			order = strings.TrimPrefix(order, "-")
		}
		if !slices.Contains(q.Fields, order) { // TODO: handle this, instead of error
			err := errors.New("order field should be included in fields")
			api.AbortRequest(c, http.StatusBadRequest, "invalid_params", err, n.log)
			return
		}
	}
	itemList, err := oua.ValueEntries(
		c, objAttrID, q.Fields, order, desc, limit, q.Offset)
	if err != nil {
		err := fmt.Errorf("query: %w", err)
		api.AbortRequest(c, http.StatusInternalServerError, "invalid query", err, n.log)
		return
	}

	result := JsonResult{
		Count:  count,
		Limit:  limit,
		Offset: q.Offset,
		Items:  itemList,
	}
	c.JSON(http.StatusOK, result)
}
