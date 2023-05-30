package node

import (
	"fmt"
	"net/http"

	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils/umid"

	"github.com/momentum-xyz/ubercontroller/utils/modify"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe/auth"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/dto"
)

type queryPluginAttribute struct {
	PluginID      string `form:"plugin_id" json:"plugin_id" binding:"required"`
	AttributeName string `form:"attribute_name" json:"attribute_name" binding:"required"`
}

// @Summary Get object attribute
// @Schemes
// @Description Returns object attribute
// @Tags objects
// @Produce json
// @Param object_id path string true "Object UMID"
// @Param attribute_id query node.queryPluginAttribute true "query params"
// @Success 200 {object} entry.AttributeValue
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/objects/{object_id}/attributes [get]
func (n *Node) apiGetObjectAttributesValue(c *gin.Context) {
	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetObjectAttributesValue: user from context")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user", err, n.log)
		return
	}

	objectID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetObjectAttributesValue: failed to parse object umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}
	object, ok := n.GetObjectFromAllObjects(objectID)
	if !ok {
		err := errors.Errorf("Node: apiGetObjectAttributesValue: object not found: %s", objectID)
		api.AbortRequest(c, http.StatusNotFound, "object_not_found", err, n.log)
		return
	}

	attrType, attributeID, err := n.apiPluginAttributeFromQuery(c)
	if err != nil {
		err := fmt.Errorf("node: apiGetObjectAttributesValue: failed to get plugin attribute: %w", err)
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_attribute", err, n.log)
		return
	}

	allowed, err := auth.CheckAttributePermissions(
		c, *attrType.GetEntry(), object.GetObjectAttributes(), attributeID, userID,
		auth.ReadOperation)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetObjectAttributesValue: permissions check")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_permissions_check", err, n.log)
		return
	} else if !allowed {
		err := fmt.Errorf("operation not permitted")
		api.AbortRequest(c, http.StatusForbidden, "operation_not_permitted", err, n.log)
		return
	}

	out, ok := object.GetObjectAttributes().GetValue(attributeID)
	if !ok {
		err := errors.Errorf("Node: apiGetObjectAttributesValue: object attribute value not found: %s", attributeID)
		api.AbortRequest(c, http.StatusNotFound, "attribute_not_found", err, n.log)
		return
	}

	c.JSON(http.StatusOK, out)
}

// @Summary Get object and all subobject attributes
// @Schemes
// @Description Returns object and all subobject attributes
// @Tags objects
// @Accept json
// @Produce json
// @Param object_id path string true "Object UMID"
// @Param query query node.queryPluginAttribute true "query params"
// @Success 200 {object} dto.ObjectAttributeValues
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/objects/{object_id}/attributes-with-children [get]
func (n *Node) apiGetObjectWithChildrenAttributeValues(c *gin.Context) {
	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "user from context")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user", err, n.log)
		return
	}

	objectID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetObjectWithChildrenAttributeValues: failed to parse object umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	attrType, attributeID, err := n.apiPluginAttributeFromQuery(c)
	if err != nil {
		err := fmt.Errorf("node: apiGetObjectWithChildrenAttributeValues: failed to get plugin attribute: %w", err)
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_attribute", err, n.log)
		return
	}

	rootObject, ok := n.GetObjectFromAllObjects(objectID)
	if !ok {
		err := errors.Errorf("Node: apiGetObjectWithChildrenAttributeValues: object not found: %s", objectID)
		api.AbortRequest(c, http.StatusNotFound, "object_not_found", err, n.log)
		return
	}

	// TODO: either remove this API method, or implement recursive permission checks...
	allowed, err := auth.CheckReadAllPermissions[entry.AttributeID](
		c, *attrType.GetEntry(), rootObject.GetObjectAttributes(), userID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetObjectAttributesValue: permissions check")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_permissions_check", err, n.log)
		return
	} else if !allowed {
		err := fmt.Errorf("operation not permitted")
		api.AbortRequest(c, http.StatusForbidden, "operation_not_permitted", err, n.log)
		return
	}

	objects := rootObject.GetObjects(true)
	objects[rootObject.GetID()] = rootObject

	objectAttributes := make(dto.ObjectAttributeValues, len(objects))
	for _, object := range objects {
		attributeValue, ok := object.GetObjectAttributes().GetValue(attributeID)
		if !ok || attributeValue == nil {
			continue
		}

		objectAttributes[object.GetID()] = attributeValue
	}

	c.JSON(http.StatusOK, objectAttributes)
}

// @Summary Makes object attribute available to public
// @Schemes
// @Description Changes object permissions to be {"write":"user", "read":"any"}
// @Tags objects
// @Accept json
// @Produce json
// @Param object_id path string true "Object UMID"
// @Param body body node.queryPluginAttribute true "body params"
// @Success 202 {object} dto.ObjectSubOptions
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/objects/{object_id}/attributes/publicize [post]
func (n *Node) apiSetObjectAttributesPublic(c *gin.Context) {
	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetObjectAttributesPublic: user from context")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user", err, n.log)
		return
	}

	objectID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetObjectAttributesPublic: failed to parse object umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	object, ok := n.GetObjectFromAllObjects(objectID)
	if !ok {
		err := errors.Errorf("Node: apiSetObjectAttributesPublic: object not found: %s", objectID)
		api.AbortRequest(c, http.StatusNotFound, "object_not_found", err, n.log)
		return
	}

	attrType, attributeID, err := n.apiPluginAttributeFromQuery(c)
	if err != nil {
		err := fmt.Errorf("node: apiGetObjectWithChildrenAttributeValues: failed to get plugin attribute: %w", err)
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_attribute", err, n.log)
		return
	}

	allowed, err := auth.CheckAttributePermissions(
		c, *attrType.GetEntry(), object.GetObjectAttributes(), attributeID, userID,
		auth.WriteOperation)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetObjectAttributesValue: permissions check")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_permissions_check", err, n.log)
		return
	} else if !allowed {
		err := fmt.Errorf("operation not permitted")
		api.AbortRequest(c, http.StatusForbidden, "operation_not_permitted", err, n.log)
		return
	}

	modifyFn := func(current *entry.AttributePayload) (*entry.AttributePayload, error) {
		newOptions := func() *entry.AttributeOptions {
			options := entry.NewAttributeOptions()
			*options = entry.AttributeOptions{
				"permissions": map[string]any{
					"read":  "any",
					"write": "user",
				},
			}
			return options
		}

		if current == nil {
			return entry.NewAttributePayload(current.Value, newOptions()), nil
		}

		if current.Options == nil {
			current.Options = newOptions()
			return current, nil
		}

		*current.Options = entry.AttributeOptions{
			"permissions": map[string]any{
				"read":  "any",
				"write": "user",
			},
		}

		return current, nil
	}

	if _, err := object.GetObjectAttributes().Upsert(attributeID, modifyFn, true); err != nil {
		err := errors.WithMessage(err, "Node: apiSetObjectAttributesPublic: failed to set options")
		api.AbortRequest(c, http.StatusInternalServerError, "set_options_failed", err, n.log)
		return
	}

	c.JSON(http.StatusAccepted, nil)
}

// @Summary Set object attribute
// @Schemes
// @Description Sets entire object attribute
// @Tags objects
// @Accept json
// @Produce json
// @Param object_id path string true "Object UMID"
// @Param body body node.apiSetObjectAttributesValue.InBody true "body params"
// @Success 202 {object} entry.AttributeValue
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/objects/{object_id}/attributes [post]
func (n *Node) apiSetObjectAttributesValue(c *gin.Context) {
	type InBody struct {
		queryPluginAttribute
		AttributeValue map[string]any `json:"attribute_value" binding:"required"`
	}

	inBody := InBody{}

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiSetObjectAttributesValue: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	objectID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetObjectAttributesValue: failed to parse object umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	pluginID, err := umid.Parse(inBody.PluginID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetObjectAttributesValue: failed to parse plugin umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}
	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "user from context")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user", err, n.log)
		return
	}
	attrType, ok := n.GetAttributeTypes().GetAttributeType(entry.AttributeTypeID{pluginID, inBody.AttributeName})
	if !ok {
		err := fmt.Errorf("attribute type not found")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_attribute", err, n.log)
		return
	}
	attributeID := entry.NewAttributeID(pluginID, inBody.AttributeName)
	object, ok := n.GetObjectFromAllObjects(objectID)
	if !ok {
		err := errors.Errorf("Node: apiSetObjectAttributesValue: object not found: %s", objectID)
		api.AbortRequest(c, http.StatusNotFound, "object_not_found", err, n.log)
		return
	}

	allowed, err := auth.CheckAttributePermissions(
		c, *attrType.GetEntry(), object.GetObjectAttributes(), attributeID, userID,
		auth.WriteOperation)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetObjectAttributesValue: permissions check")
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

	payload, err := object.GetObjectAttributes().Upsert(attributeID, modifyFn, true)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiSetObjectAttributesValue: failed to upsert object attribute")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_upsert", err, n.log)
		return
	}

	c.JSON(http.StatusAccepted, payload.Value)
}

// @Summary Get object sub attributes
// @Schemes
// @Description Returns object sub attributes
// @Tags objects
// @Accept json
// @Produce json
// @Param object_id path string true "Object UMID"
// @Param query query node.apiGetObjectAttributeSubValue.InQuery true "query params"
// @Success 200 {object} dto.ObjectSubAttributes
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/objects/{object_id}/attributes/sub [get]
func (n *Node) apiGetObjectAttributeSubValue(c *gin.Context) {
	type InQuery struct {
		queryPluginAttribute
		SubAttributeKey string `form:"sub_attribute_key" binding:"required"`
	}

	inQuery := InQuery{}

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Node: apiGetObjectAttributeSubValue: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, n.log)
		return
	}

	objectID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetObjectSubAttributes: failed to parse object umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	pluginID, err := umid.Parse(inQuery.PluginID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetObjectSubAttributes: failed to parse plugin umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}

	object, ok := n.GetObjectFromAllObjects(objectID)
	if !ok {
		err := errors.Errorf("Node: apiGetObjectAttributeSubValue: object not found: %s", objectID)
		api.AbortRequest(c, http.StatusNotFound, "object_not_found", err, n.log)
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
		err := errors.WithMessage(err, "user from context")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user", err, n.log)
		return
	}
	attributeID := entry.NewAttributeID(pluginID, inQuery.AttributeName)

	allowed, err := auth.CheckAttributePermissions(
		c, *attrType.GetEntry(), object.GetObjectAttributes(), attributeID, userID,
		auth.ReadOperation)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetObjectAttributesValue: permissions check")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_permissions_check", err, n.log)
		return
	} else if !allowed {
		err := fmt.Errorf("operation not permitted")
		api.AbortRequest(c, http.StatusForbidden, "operation_not_permitted", err, n.log)
		return
	}

	attributeValue, ok := object.GetObjectAttributes().GetValue(attributeID)
	if !ok {
		err := errors.Errorf("Node: apiGetObjectAttributeSubValue: attribute value not found: %s", attributeID)
		api.AbortRequest(c, http.StatusNotFound, "attribute_value_not_found", err, n.log)
		return
	}

	if attributeValue == nil {
		err := errors.Errorf("Node: apiGetObjectAttributeSubValue: attribute value is nil")
		api.AbortRequest(c, http.StatusNotFound, "attribute_value_nil", err, n.log)
		return
	}

	out := dto.ObjectSubAttributes{
		inQuery.SubAttributeKey: (*attributeValue)[inQuery.SubAttributeKey],
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
// @Param body body node.apiSetObjectAttributeSubValue.Body true "body params"
// @Success 202 {object} dto.ObjectSubAttributes
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/objects/{object_id}/attributes/sub [post]
func (n *Node) apiSetObjectAttributeSubValue(c *gin.Context) {
	type Body struct {
		queryPluginAttribute
		SubAttributeKey   string `json:"sub_attribute_key" binding:"required"`
		SubAttributeValue any    `json:"sub_attribute_value" binding:"required"`
	}

	inBody := Body{}

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiSetObjectAttributeSubValue: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	objectID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetObjectAttributeSubValue: failed to parse object umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	pluginID, err := umid.Parse(inBody.PluginID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetObjectAttributeSubValue: failed to parse plugin umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}
	object, ok := n.GetObjectFromAllObjects(objectID)
	if !ok {
		err := errors.Errorf("Node: apiSetObjectAttributeSubValue: object not found: %s", objectID)
		api.AbortRequest(c, http.StatusNotFound, "object_not_found", err, n.log)
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
		err := errors.WithMessage(err, "user from context")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user", err, n.log)
		return
	}
	attributeID := entry.NewAttributeID(pluginID, inBody.AttributeName)

	allowed, err := auth.CheckAttributePermissions(
		c, *attrType.GetEntry(), object.GetObjectAttributes(), attributeID, userID,
		auth.WriteOperation)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetObjectAttributesValue: permissions check")
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

	payload, err := object.GetObjectAttributes().Upsert(attributeID, modifyFn, true)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiSetObjectAttributeSubValue: failed to upsert object attribute")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_upsert", err, n.log)
		return
	}

	out := dto.ObjectSubAttributes{
		inBody.SubAttributeKey: (*payload.Value)[inBody.SubAttributeKey],
	}

	c.JSON(http.StatusAccepted, out)
}

// @Summary Delete object sub attribute
// @Schemes
// @Description Deletes a object sub attribute
// @Tags objects
// @Accept json
// @Produce json
// @Param object_id path string true "Object UMID"
// @Param body body node.apiRemoveObjectAttributeSubValue.Body true "body params"
// @Success 200 {object} nil
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/objects/{object_id}/attributes/sub [delete]
func (n *Node) apiRemoveObjectAttributeSubValue(c *gin.Context) {
	type Body struct {
		queryPluginAttribute
		SubAttributeKey string `json:"sub_attribute_key" binding:"required"`
	}

	inBody := Body{}

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiRemoveObjectAttributeSubValue: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	objectID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveObjectAttributeSubValue: failed to parse object umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	pluginID, err := umid.Parse(inBody.PluginID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveObjectAttributeSubValue: failed to parse plugin umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}
	object, ok := n.GetObjectFromAllObjects(objectID)
	if !ok {
		err := errors.Errorf("Node: apiRemoveObjectAttributeSubValue: object not found: %s", objectID)
		api.AbortRequest(c, http.StatusNotFound, "object_not_found", err, n.log)
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
		err := errors.WithMessage(err, "user from context")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user", err, n.log)
		return
	}
	attributeID := entry.NewAttributeID(pluginID, inBody.AttributeName)

	allowed, err := auth.CheckAttributePermissions(
		c, *attrType.GetEntry(), object.GetObjectAttributes(), attributeID, userID,
		auth.WriteOperation)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetObjectAttributesValue: permissions check")
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

	if _, err := object.GetObjectAttributes().UpdateValue(attributeID, modifyFn, true); err != nil {
		err = errors.WithMessage(err, "Node: apiRemoveObjectAttributeSubValue: failed to update object attribute")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_update", err, n.log)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary Delete object attribute
// @Schemes
// @Description Deletes a object attribute
// @Tags objects
// @Accept json
// @Produce json
// @Param object_id path string true "Object UMID"
// @Param body body node.queryPluginAttribute true "body params"
// @Success 200 {object} nil
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/objects/{object_id}/attributes [delete]
func (n *Node) apiRemoveObjectAttributeValue(c *gin.Context) {
	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveObjectAttributeValue: user from context")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user", err, n.log)
		return
	}

	objectID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveObjectAttributeValue: failed to parse object umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}
	object, ok := n.GetObjectFromAllObjects(objectID)
	if !ok {
		err := errors.Errorf("Node: apiRemoveObjectAttributeValue: object not found: %s", objectID)
		api.AbortRequest(c, http.StatusNotFound, "object_not_found", err, n.log)
		return
	}

	attrType, attributeID, err := n.apiPluginAttributeFromQuery(c)
	if err != nil {
		err := fmt.Errorf("node: apiGetObjectAttributesValue: failed to get plugin attribute: %w", err)
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_attribute", err, n.log)
		return
	}

	allowed, err := auth.CheckAttributePermissions(
		c, *attrType.GetEntry(), object.GetObjectAttributes(), attributeID, userID,
		auth.WriteOperation)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetObjectAttributesValue: permissions check")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_permissions_check", err, n.log)
		return
	} else if !allowed {
		err := fmt.Errorf("operation not permitted")
		api.AbortRequest(c, http.StatusForbidden, "operation_not_permitted", err, n.log)
		return
	}

	if _, err := object.GetObjectAttributes().UpdateValue(
		attributeID, modify.ReplaceWith[entry.AttributeValue](nil), true,
	); err != nil {
		err = errors.WithMessage(err, "Node: apiRemoveObjectAttributeValue: failed to update object attribute")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_update", err, n.log)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// TODO: refactor this whole file to an object_attribute package,
// so we don't need these kinds of silly names to avoid collisions.
func (n *Node) apiPluginAttributeFromQuery(c *gin.Context) (universe.AttributeType, entry.AttributeID, error) {
	var attrID entry.AttributeID
	inQuery := queryPluginAttribute{}
	if err := c.ShouldBindQuery(&inQuery); err != nil {
		return nil, attrID, fmt.Errorf("failed to bind query: %w", err)
	}
	pluginID, err := umid.Parse(inQuery.PluginID)
	if err != nil {
		return nil, attrID, fmt.Errorf("failed to parse plugin ID: %w", err)
	}
	attrType, ok := n.GetAttributeTypes().GetAttributeType(
		entry.AttributeTypeID{PluginID: pluginID, Name: inQuery.AttributeName})
	if !ok {
		return nil, attrID, fmt.Errorf("attribute type for %+v not found", inQuery)
	}
	attrID = entry.NewAttributeID(pluginID, inQuery.AttributeName)
	return attrType, attrID, nil
}
