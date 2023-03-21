package node

import (
	"net/http"

	"github.com/momentum-xyz/ubercontroller/utils/umid"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
)

// @Summary Get user attribute
// @Schemes
// @Description Returns user attribute
// @Tags users
// @Accept json
// @Produce json
// @Param user_id path string true "User UMID"
// @Param query query node.apiGetUserAttributeValue.InQuery true "query params"
// @Success 200 {user} entry.AttributeValue
// @Failure 500 {user} api.HTTPError
// @Failure 400 {user} api.HTTPError
// @Failure 404 {user} api.HTTPError
// @Router /api/v4/users/{user_id}/attributes [get]
func (n *Node) apiGetUserAttributeValue(c *gin.Context) {
	type InQuery struct {
		PluginID      string `form:"plugin_id" binding:"required"`
		AttributeName string `form:"attribute_name" binding:"required"`
	}

	inQuery := InQuery{}

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Node: apiGetUserAttributeValue: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, n.log)
		return
	}

	userID, err := umid.Parse(c.Param("userID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetUserAttributeValue: failed to parse user umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	pluginID, err := umid.Parse(inQuery.PluginID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetUserAttributeValue: failed to parse plugin umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}

	attributeID := entry.NewAttributeID(pluginID, inQuery.AttributeName)
	userAttributeID := entry.NewUserAttributeID(attributeID, userID)
	
	out, ok := n.GetUserAttributes().GetValue(userAttributeID)
	if !ok {
		err := errors.Errorf("Node: apiGetUserAttributeValue: user attribute value not found: %s", attributeID)
		api.AbortRequest(c, http.StatusNotFound, "attribute_not_found", err, n.log)
		return
	}

	c.JSON(http.StatusOK, out)
}

// @Summary Set user attribute
// @Schemes
// @Description Sets entire user attribute
// @Tags users
// @Accept json
// @Produce json
// @Param user_id path string true "User UMID"
// @Param body body node.apiSetUserAttributeValue.InBody true "body params"
// @Success 202 {user} entry.AttributeValue
// @Failure 500 {user} api.HTTPError
// @Failure 400 {user} api.HTTPError
// @Failure 404 {user} api.HTTPError
// @Router /api/v4/users/{user_id}/attributes [post]
//func (n *Node) apiSetUserAttributeValue(c *gin.Context) {
//	type InBody struct {
//		PluginID       string         `json:"plugin_id" binding:"required"`
//		AttributeName  string         `json:"attribute_name" binding:"required"`
//		AttributeValue map[string]any `json:"attribute_value" binding:"required"`
//	}
//
//	inBody := InBody{}
//
//	if err := c.ShouldBindJSON(&inBody); err != nil {
//		err = errors.WithMessage(err, "Node: apiSetUserAttributeValue: failed to bind json")
//		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
//		return
//	}
//
//	userID, err := umid.Parse(c.Param("userID"))
//	if err != nil {
//		err := errors.WithMessage(err, "Node: apiSetUserAttributeValue: failed to parse user umid")
//		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
//		return
//	}
//
//	pluginID, err := umid.Parse(inBody.PluginID)
//	if err != nil {
//		err := errors.WithMessage(err, "Node: apiSetUserAttributeValue: failed to parse plugin umid")
//		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
//		return
//	}
//
//	user, ok := n.GetUserFromAllUsers(userID)
//	if !ok {
//		err := errors.Errorf("Node: apiSetUserAttributeValue: user not found: %s", userID)
//		api.AbortRequest(c, http.StatusNotFound, "user_not_found", err, n.log)
//		return
//	}
//
//	attributeID := entry.NewAttributeID(pluginID, inBody.AttributeName)
//
//	modifyFn := func(current *entry.AttributePayload) (*entry.AttributePayload, error) {
//		newValue := func() *entry.AttributeValue {
//			value := entry.NewAttributeValue()
//			*value = inBody.AttributeValue
//			return value
//		}
//
//		if current == nil {
//			return entry.NewAttributePayload(newValue(), nil), nil
//		}
//
//		if current.Value == nil {
//			current.Value = newValue()
//			return current, nil
//		}
//
//		*current.Value = inBody.AttributeValue
//
//		return current, nil
//	}
//
//	payload, err := user.GetUserAttributes().Upsert(attributeID, modifyFn, true)
//	if err != nil {
//		err = errors.WithMessage(err, "Node: apiSetUserAttributeValue: failed to upsert user attribute")
//		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_upsert", err, n.log)
//		return
//	}
//
//	c.JSON(http.StatusAccepted, payload.Value)
//}
//
//// @Summary Get user sub attributes
//// @Schemes
//// @Description Returns user sub attributes
//// @Tags users
//// @Accept json
//// @Produce json
//// @Param user_id path string true "User UMID"
//// @Param query query node.apiGetUserAttributeSubValue.InQuery true "query params"
//// @Success 200 {user} dto.UserSubAttributes
//// @Failure 500 {user} api.HTTPError
//// @Failure 400 {user} api.HTTPError
//// @Failure 404 {user} api.HTTPError
//// @Router /api/v4/users/{user_id}/attributes/sub [get]
//func (n *Node) apiGetUserAttributeSubValue(c *gin.Context) {
//	type InQuery struct {
//		PluginID        string `form:"plugin_id" binding:"required"`
//		AttributeName   string `form:"attribute_name" binding:"required"`
//		SubAttributeKey string `form:"sub_attribute_key" binding:"required"`
//	}
//
//	inQuery := InQuery{}
//
//	if err := c.ShouldBindQuery(&inQuery); err != nil {
//		err := errors.WithMessage(err, "Node: apiGetUserAttributeSubValue: failed to bind query")
//		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, n.log)
//		return
//	}
//
//	userID, err := umid.Parse(c.Param("userID"))
//	if err != nil {
//		err := errors.WithMessage(err, "Node: apiGetUserSubAttributes: failed to parse user umid")
//		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
//		return
//	}
//
//	pluginID, err := umid.Parse(inQuery.PluginID)
//	if err != nil {
//		err := errors.WithMessage(err, "Node: apiGetUserSubAttributes: failed to parse plugin umid")
//		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
//		return
//	}
//
//	user, ok := n.GetUserFromAllUsers(userID)
//	if !ok {
//		err := errors.Errorf("Node: apiGetUserAttributeSubValue: user not found: %s", userID)
//		api.AbortRequest(c, http.StatusNotFound, "user_not_found", err, n.log)
//		return
//	}
//
//	attributeID := entry.NewAttributeID(pluginID, inQuery.AttributeName)
//	attributeValue, ok := user.GetUserAttributes().GetValue(attributeID)
//	if !ok {
//		err := errors.Errorf("Node: apiGetUserAttributeSubValue: attribute value not found: %s", attributeID)
//		api.AbortRequest(c, http.StatusNotFound, "attribute_value_not_found", err, n.log)
//		return
//	}
//
//	if attributeValue == nil {
//		err := errors.Errorf("Node: apiGetUserAttributeSubValue: attribute value is nil")
//		api.AbortRequest(c, http.StatusNotFound, "attribute_value_nil", err, n.log)
//		return
//	}
//
//	out := dto.UserSubAttributes{
//		inQuery.SubAttributeKey: (*attributeValue)[inQuery.SubAttributeKey],
//	}
//
//	c.JSON(http.StatusOK, out)
//}
//
//// @Summary Set user sub attribute
//// @Schemes
//// @Description Sets a user sub attribute
//// @Tags users
//// @Accept json
//// @Produce json
//// @Param user_id path string true "User UMID"
//// @Param body body node.apiSetUserAttributeSubValue.Body true "body params"
//// @Success 202 {user} dto.UserSubAttributes
//// @Failure 500 {user} api.HTTPError
//// @Failure 400 {user} api.HTTPError
//// @Failure 404 {user} api.HTTPError
//// @Router /api/v4/users/{user_id}/attributes/sub [post]
//func (n *Node) apiSetUserAttributeSubValue(c *gin.Context) {
//	type Body struct {
//		PluginID          string `json:"plugin_id" binding:"required"`
//		AttributeName     string `json:"attribute_name" binding:"required"`
//		SubAttributeKey   string `json:"sub_attribute_key" binding:"required"`
//		SubAttributeValue any    `json:"sub_attribute_value" binding:"required"`
//	}
//
//	inBody := Body{}
//
//	if err := c.ShouldBindJSON(&inBody); err != nil {
//		err = errors.WithMessage(err, "Node: apiSetUserAttributeSubValue: failed to bind json")
//		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
//		return
//	}
//
//	userID, err := umid.Parse(c.Param("userID"))
//	if err != nil {
//		err := errors.WithMessage(err, "Node: apiSetUserAttributeSubValue: failed to parse user umid")
//		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
//		return
//	}
//
//	pluginID, err := umid.Parse(inBody.PluginID)
//	if err != nil {
//		err := errors.WithMessage(err, "Node: apiSetUserAttributeSubValue: failed to parse plugin umid")
//		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
//		return
//	}
//
//	user, ok := n.GetUserFromAllUsers(userID)
//	if !ok {
//		err := errors.Errorf("Node: apiSetUserAttributeSubValue: user not found: %s", userID)
//		api.AbortRequest(c, http.StatusNotFound, "user_not_found", err, n.log)
//		return
//	}
//
//	attributeID := entry.NewAttributeID(pluginID, inBody.AttributeName)
//
//	modifyFn := func(current *entry.AttributePayload) (*entry.AttributePayload, error) {
//		newValue := func() *entry.AttributeValue {
//			value := entry.NewAttributeValue()
//			(*value)[inBody.SubAttributeKey] = inBody.SubAttributeValue
//			return value
//		}
//
//		if current == nil {
//			return entry.NewAttributePayload(newValue(), nil), nil
//		}
//
//		if current.Value == nil {
//			current.Value = newValue()
//			return current, nil
//		}
//
//		(*current.Value)[inBody.SubAttributeKey] = inBody.SubAttributeValue
//
//		return current, nil
//	}
//
//	payload, err := user.GetUserAttributes().Upsert(attributeID, modifyFn, true)
//	if err != nil {
//		err = errors.WithMessage(err, "Node: apiSetUserAttributeSubValue: failed to upsert user attribute")
//		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_upsert", err, n.log)
//		return
//	}
//
//	out := dto.UserSubAttributes{
//		inBody.SubAttributeKey: (*payload.Value)[inBody.SubAttributeKey],
//	}
//
//	c.JSON(http.StatusAccepted, out)
//}
//
//// @Summary Delete user sub attribute
//// @Schemes
//// @Description Deletes a user sub attribute
//// @Tags users
//// @Accept json
//// @Produce json
//// @Param user_id path string true "User UMID"
//// @Param body body node.apiRemoveUserAttributeSubValue.Body true "body params"
//// @Success 200 {user} nil
//// @Failure 500 {user} api.HTTPError
//// @Failure 400 {user} api.HTTPError
//// @Failure 404 {user} api.HTTPError
//// @Router /api/v4/users/{user_id}/attributes/sub [delete]
//func (n *Node) apiRemoveUserAttributeSubValue(c *gin.Context) {
//	type Body struct {
//		PluginID        string `json:"plugin_id" binding:"required"`
//		AttributeName   string `json:"attribute_name" binding:"required"`
//		SubAttributeKey string `json:"sub_attribute_key" binding:"required"`
//	}
//
//	inBody := Body{}
//
//	if err := c.ShouldBindJSON(&inBody); err != nil {
//		err = errors.WithMessage(err, "Node: apiRemoveUserAttributeSubValue: failed to bind json")
//		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
//		return
//	}
//
//	userID, err := umid.Parse(c.Param("userID"))
//	if err != nil {
//		err := errors.WithMessage(err, "Node: apiRemoveUserAttributeSubValue: failed to parse user umid")
//		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
//		return
//	}
//
//	pluginID, err := umid.Parse(inBody.PluginID)
//	if err != nil {
//		err := errors.WithMessage(err, "Node: apiRemoveUserAttributeSubValue: failed to parse plugin umid")
//		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
//		return
//	}
//
//	user, ok := n.GetUserFromAllUsers(userID)
//	if !ok {
//		err := errors.Errorf("Node: apiRemoveUserAttributeSubValue: user not found: %s", userID)
//		api.AbortRequest(c, http.StatusNotFound, "user_not_found", err, n.log)
//		return
//	}
//
//	attributeID := entry.NewAttributeID(pluginID, inBody.AttributeName)
//
//	modifyFn := func(current *entry.AttributeValue) (*entry.AttributeValue, error) {
//		if current == nil {
//			return current, nil
//		}
//
//		delete(*current, inBody.SubAttributeKey)
//
//		return current, nil
//	}
//
//	if _, err := user.GetUserAttributes().UpdateValue(attributeID, modifyFn, true); err != nil {
//		err = errors.WithMessage(err, "Node: apiRemoveUserAttributeSubValue: failed to update user attribute")
//		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_update", err, n.log)
//		return
//	}
//
//	c.JSON(http.StatusOK, nil)
//}
//
//// @Summary Delete user attribute
//// @Schemes
//// @Description Deletes a user attribute
//// @Tags users
//// @Accept json
//// @Produce json
//// @Param user_id path string true "User UMID"
//// @Param body body node.apiRemoveUserAttributeValue.Body true "body params"
//// @Success 200 {user} nil
//// @Failure 500 {user} api.HTTPError
//// @Failure 400 {user} api.HTTPError
//// @Failure 404 {user} api.HTTPError
//// @Router /api/v4/users/{user_id}/attributes [delete]
//func (n *Node) apiRemoveUserAttributeValue(c *gin.Context) {
//	type Body struct {
//		PluginID      string `json:"plugin_id" binding:"required"`
//		AttributeName string `json:"attribute_name" binding:"required"`
//	}
//
//	var inBody Body
//	if err := c.ShouldBindJSON(&inBody); err != nil {
//		err = errors.WithMessage(err, "Node: apiRemoveUserAttributeValue: failed to bind json")
//		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
//		return
//	}
//
//	userID, err := umid.Parse(c.Param("userID"))
//	if err != nil {
//		err := errors.WithMessage(err, "Node: apiRemoveUserAttributeValue: failed to parse user umid")
//		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
//		return
//	}
//
//	pluginID, err := umid.Parse(inBody.PluginID)
//	if err != nil {
//		err := errors.WithMessage(err, "Node: apiRemoveUserAttributeValue: failed to parse plugin umid")
//		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
//		return
//	}
//
//	user, ok := n.GetUserFromAllUsers(userID)
//	if !ok {
//		err := errors.Errorf("Node: apiRemoveUserAttributeValue: user not found: %s", userID)
//		api.AbortRequest(c, http.StatusNotFound, "user_not_found", err, n.log)
//		return
//	}
//
//	attributeID := entry.NewAttributeID(pluginID, inBody.AttributeName)
//	if _, err := user.GetUserAttributes().UpdateValue(
//		attributeID, modify.ReplaceWith[entry.AttributeValue](nil), true,
//	); err != nil {
//		err = errors.WithMessage(err, "Node: apiRemoveUserAttributeValue: failed to update user attribute")
//		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_update", err, n.log)
//		return
//	}
//
//	c.JSON(http.StatusOK, nil)
//}
