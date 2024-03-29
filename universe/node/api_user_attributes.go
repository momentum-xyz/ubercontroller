package node

import (
	"fmt"
	"net/http"

	"github.com/momentum-xyz/ubercontroller/universe/attributes"
	"github.com/momentum-xyz/ubercontroller/universe/auth"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/dto"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
	"github.com/momentum-xyz/ubercontroller/utils/umid"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
)

// @Summary Get user attribute for own user based on token
// @Description Returns user attribute
// @Tags attributes,users
// @Security Bearer
// @Param query query attributes.QueryPluginAttribute true "query params"
// @Success 200 {object} entry.AttributeValue
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/users/me/attributes [get]
func (n *Node) apiGetMeUserAttributeValue(c *gin.Context) {
	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiGetMeUserAttributeValue: failed to get user umid")
		api.AbortRequest(c, http.StatusInternalServerError, "get_user_id_failed", err, n.log)
		return
	}

	_, attributeID, err := attributes.PluginAttributeFromQuery(c, n)
	if err != nil {
		err := fmt.Errorf("failed to get plugin attribute: %w", err)
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_attribute", err, n.log)
		return
	}

	userAttributeID := entry.NewUserAttributeID(attributeID, userID)

	// TODO: permission check? In theory we could have 'hidden' user attrs used by admins or plugins.

	out, ok := n.GetUserAttributes().GetValue(userAttributeID)
	if !ok {
		err := errors.Errorf("Node: apiGetMeUserAttributeValue: user attribute value not found: %s", attributeID)
		api.AbortRequest(c, http.StatusNotFound, "attribute_not_found", err, n.log)
		return
	}

	c.JSON(http.StatusOK, out)
}

// @Summary Get user attribute
// @Description Returns user attribute
// @Tags attributes,users
// @Security Bearer
// @Param user_id path string true "User UMID"
// @Param query query attributes.QueryPluginAttribute true "query params"
// @Success 200 {object} entry.AttributeValue
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/users/{user_id}/attributes [get]
func (n *Node) apiGetUserAttributeValue(c *gin.Context) {
	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiGetUserAttributeValue: failed to get user umid")
		api.AbortRequest(c, http.StatusInternalServerError, "get_user_id_failed", err, n.log)
		return
	}

	targetUserID, err := umid.Parse(c.Param("userID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetUserAttributeValue: failed to parse user umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	attrType, attributeID, err := attributes.PluginAttributeFromQuery(c, n)
	if err != nil {
		err := fmt.Errorf("failed to get plugin attribute: %w", err)
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_attribute", err, n.log)
		return
	}

	userAttributeID := entry.NewUserAttributeID(attributeID, targetUserID)

	allowed, err := auth.CheckAttributePermissions(
		c, *attrType.GetEntry(), n.GetUserAttributes(), userAttributeID, userID,
		auth.ReadOperation,
	)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetUserAttributeValue: permissions check")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_permissions_check", err, n.log)
		return
	} else if !allowed {
		err := fmt.Errorf("operation not permitted")
		api.AbortRequest(c, http.StatusForbidden, "operation_not_permitted", err, n.log)
		return
	}
	out, ok := n.GetUserAttributes().GetValue(userAttributeID)
	if !ok {
		err := errors.Errorf("Node: apiGetUserAttributeValue: user attribute value not found: %s", attributeID)
		api.AbortRequest(c, http.StatusNotFound, "attribute_not_found", err, n.log)
		return
	}

	c.JSON(http.StatusOK, out)
}

// @Summary Set user attribute
// @Description Sets entire user attribute
// @Tags attributes,users
// @Security Bearer
// @Param user_id path string true "User UMID"
// @Param body body node.apiSetUserAttributeValue.InBody true "body params"
// @Success 202 {object} entry.AttributeValue
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/users/{user_id}/attributes [post]
func (n *Node) apiSetUserAttributeValue(c *gin.Context) {
	type InBody struct {
		attributes.QueryPluginAttribute
		AttributeValue map[string]any `json:"attribute_value" binding:"required"`
	}

	inBody := InBody{}

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiSetUserAttributeValue: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	targetUserID, err := umid.Parse(c.Param("userID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetUserAttributeValue: failed to parse user umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	pluginID, err := umid.Parse(inBody.PluginID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetUserAttributeValue: failed to parse plugin umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}
	attrType, ok := n.GetAttributeTypes().GetAttributeType(entry.AttributeTypeID{pluginID, inBody.AttributeName})
	if !ok {
		err := fmt.Errorf("attribute type not found")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_attribute", err, n.log)
		return
	}
	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiGetUserAttributeValue: failed to get user umid")
		api.AbortRequest(c, http.StatusInternalServerError, "get_user_id_failed", err, n.log)
		return
	}
	attributeID := entry.NewAttributeID(pluginID, inBody.AttributeName)
	userAttributeID := entry.NewUserAttributeID(attributeID, targetUserID)

	allowed, err := auth.CheckAttributePermissions(
		c, *attrType.GetEntry(), n.GetUserAttributes(), userAttributeID, userID,
		auth.WriteOperation,
	)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetUserAttributeValue: permissions check")
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

	payload, err := n.GetUserAttributes().Upsert(userAttributeID, modifyFn, true)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiSetUserAttributeValue: failed to upsert user attribute")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_upsert", err, n.log)
		return
	}

	c.JSON(http.StatusAccepted, payload.Value)
}

// @Summary Get user sub attributes
// @Description Returns user sub attributes
// @Tags attributes,users
// @Security Bearer
// @Param user_id path string true "User UMID"
// @Param query query node.apiGetUserAttributeSubValue.InQuery true "query params"
// @Success 200 {object} dto.UserSubAttributes
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/users/{user_id}/attributes/sub [get]
func (n *Node) apiGetUserAttributeSubValue(c *gin.Context) {
	type InQuery struct {
		attributes.QueryPluginAttribute
		SubAttributeKey string `form:"sub_attribute_key" json:"sub_attribute_key" binding:"required"`
	}

	inQuery := InQuery{}

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Node: apiGetUserAttributeSubValue: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, n.log)
		return
	}

	targetUserID, err := umid.Parse(c.Param("userID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetUserSubAttributes: failed to parse user umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	pluginID, err := umid.Parse(inQuery.PluginID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetUserSubAttributes: failed to parse plugin umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}
	attrType, ok := n.GetAttributeTypes().GetAttributeType(entry.AttributeTypeID{pluginID, inQuery.AttributeName})
	if !ok {
		err := fmt.Errorf("attribute type not found")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_attribute", err, n.log)
		return
	}
	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiGetUserAttributeValue: failed to get user umid")
		api.AbortRequest(c, http.StatusInternalServerError, "get_user_id_failed", err, n.log)
		return
	}
	attributeID := entry.NewAttributeID(pluginID, inQuery.AttributeName)
	userAttributeID := entry.NewUserAttributeID(attributeID, targetUserID)

	allowed, err := auth.CheckAttributePermissions(
		c, *attrType.GetEntry(), n.GetUserAttributes(), userAttributeID, userID,
		auth.ReadOperation,
	)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetUserAttributeValue: permissions check")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_permissions_check", err, n.log)
		return
	} else if !allowed {
		err := fmt.Errorf("operation not permitted")
		api.AbortRequest(c, http.StatusForbidden, "operation_not_permitted", err, n.log)
		return
	}

	attributeValue, ok := n.GetUserAttributes().GetValue(userAttributeID)
	if !ok {
		err := errors.Errorf("Node: apiGetUserAttributeSubValue: attribute value not found: %s", attributeID)
		api.AbortRequest(c, http.StatusNotFound, "attribute_value_not_found", err, n.log)
		return
	}

	if attributeValue == nil {
		err := errors.Errorf("Node: apiGetUserAttributeSubValue: attribute value is nil")
		api.AbortRequest(c, http.StatusNotFound, "attribute_value_nil", err, n.log)
		return
	}

	out := dto.UserSubAttributes{
		inQuery.SubAttributeKey: (*attributeValue)[inQuery.SubAttributeKey],
	}

	c.JSON(http.StatusOK, out)
}

// @Summary Set user sub attribute
// @Description Sets a user sub attribute
// @Tags attributes,users
// @Security Bearer
// @Param user_id path string true "User UMID"
// @Param body body node.apiSetUserAttributeSubValue.Body true "body params"
// @Success 202 {object} dto.UserSubAttributes
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/users/{user_id}/attributes/sub [post]
func (n *Node) apiSetUserAttributeSubValue(c *gin.Context) {
	type Body struct {
		attributes.QueryPluginAttribute
		SubAttributeKey   string `json:"sub_attribute_key" binding:"required"`
		SubAttributeValue any    `json:"sub_attribute_value" binding:"required"`
	}

	inBody := Body{}

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiSetUserAttributeSubValue: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	targetUserID, err := umid.Parse(c.Param("userID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetUserAttributeSubValue: failed to parse user umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	pluginID, err := umid.Parse(inBody.PluginID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetUserAttributeSubValue: failed to parse plugin umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}
	attrType, ok := n.GetAttributeTypes().GetAttributeType(entry.AttributeTypeID{pluginID, inBody.AttributeName})
	if !ok {
		err := fmt.Errorf("attribute type not found")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_attribute", err, n.log)
		return
	}
	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiGetUserAttributeValue: failed to get user umid")
		api.AbortRequest(c, http.StatusInternalServerError, "get_user_id_failed", err, n.log)
		return
	}
	attributeID := entry.NewAttributeID(pluginID, inBody.AttributeName)
	userAttributeID := entry.NewUserAttributeID(attributeID, targetUserID)

	allowed, err := auth.CheckAttributePermissions(
		c, *attrType.GetEntry(), n.GetUserAttributes(), userAttributeID, userID,
		auth.WriteOperation,
	)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetUserAttributeValue: permissions check")
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

	payload, err := n.GetUserAttributes().Upsert(userAttributeID, modifyFn, true)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiSetUserAttributeSubValue: failed to upsert user attribute")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_upsert", err, n.log)
		return
	}

	out := dto.UserSubAttributes{
		inBody.SubAttributeKey: (*payload.Value)[inBody.SubAttributeKey],
	}

	c.JSON(http.StatusAccepted, out)
}

// @Summary Delete user sub attribute
// @Description Deletes a user sub attribute
// @Tags users
// @Security Bearer
// @Param user_id path string true "User UMID"
// @Param body body node.apiRemoveUserAttributeSubValue.Body true "body params"
// @Success 200 ""
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/users/{user_id}/attributes/sub [delete]
func (n *Node) apiRemoveUserAttributeSubValue(c *gin.Context) {
	type Body struct {
		attributes.QueryPluginAttribute
		SubAttributeKey string `json:"sub_attribute_key" binding:"required"`
	}

	inBody := Body{}

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiRemoveUserAttributeSubValue: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	targetUserID, err := umid.Parse(c.Param("userID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveUserAttributeSubValue: failed to parse user umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	pluginID, err := umid.Parse(inBody.PluginID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveUserAttributeSubValue: failed to parse plugin umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}
	attrType, ok := n.GetAttributeTypes().GetAttributeType(entry.AttributeTypeID{pluginID, inBody.AttributeName})
	if !ok {
		err := fmt.Errorf("attribute type not found")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_attribute", err, n.log)
		return
	}
	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiGetUserAttributeValue: failed to get user umid")
		api.AbortRequest(c, http.StatusInternalServerError, "get_user_id_failed", err, n.log)
		return
	}
	attributeID := entry.NewAttributeID(pluginID, inBody.AttributeName)
	userAttributeID := entry.NewUserAttributeID(attributeID, targetUserID)

	allowed, err := auth.CheckAttributePermissions(
		c, *attrType.GetEntry(), n.GetUserAttributes(), userAttributeID, userID,
		auth.WriteOperation,
	)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetUserAttributeValue: permissions check")
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

	if _, err := n.GetUserAttributes().UpdateValue(userAttributeID, modifyFn, true); err != nil {
		err = errors.WithMessage(err, "Node: apiRemoveUserAttributeSubValue: failed to update user attribute")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_update", err, n.log)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary Delete user attribute
// @Description Deletes a user attribute
// @Tags users
// @Security Bearer
// @Param user_id path string true "User UMID"
// @Param body body node.apiRemoveUserAttributeValue.Body true "body params"
// @Success 200 ""
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/users/{user_id}/attributes [delete]
func (n *Node) apiRemoveUserAttributeValue(c *gin.Context) {
	type Body struct {
		attributes.QueryPluginAttribute
	}

	var inBody Body
	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiRemoveUserAttributeValue: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	targetUserID, err := umid.Parse(c.Param("userID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveUserAttributeValue: failed to parse user umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	pluginID, err := umid.Parse(inBody.PluginID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveUserAttributeValue: failed to parse plugin umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}
	attrType, ok := n.GetAttributeTypes().GetAttributeType(entry.AttributeTypeID{pluginID, inBody.AttributeName})
	if !ok {
		err := fmt.Errorf("attribute type not found")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_attribute", err, n.log)
		return
	}
	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiGetUserAttributeValue: failed to get user umid")
		api.AbortRequest(c, http.StatusInternalServerError, "get_user_id_failed", err, n.log)
		return
	}
	attributeID := entry.NewAttributeID(pluginID, inBody.AttributeName)
	userAttributeID := entry.NewUserAttributeID(attributeID, targetUserID)

	allowed, err := auth.CheckAttributePermissions(
		c, *attrType.GetEntry(), n.GetUserAttributes(), userAttributeID, userID,
		auth.WriteOperation,
	)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetUserAttributeValue: permissions check")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_permissions_check", err, n.log)
		return
	} else if !allowed {
		err := fmt.Errorf("operation not permitted")
		api.AbortRequest(c, http.StatusForbidden, "operation_not_permitted", err, n.log)
		return
	}

	if _, err := n.GetUserAttributes().UpdateValue(
		userAttributeID, modify.ReplaceWith[entry.AttributeValue](nil), true,
	); err != nil {
		err = errors.WithMessage(err, "Node: apiRemoveUserAttributeValue: failed to update user attribute")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_update", err, n.log)
		return
	}

	c.JSON(http.StatusOK, nil)
}
