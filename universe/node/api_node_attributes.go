package node

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe/attributes"
	"github.com/momentum-xyz/ubercontroller/universe/auth"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

// @Summary Get node attribute
// @Description Returns node attribute
// @Tags attributes,node
// @Security Bearer
// @Param attribute_id query attributes.QueryPluginAttribute true "query params"
// @Success 200 {object} entry.AttributeValue
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/node/attributes [get]
func (n *Node) apiNodeGetAttributesValue(c *gin.Context) {
	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiNodeGetAttributesValue: user from context")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user", err, n.log)
		return
	}

	attrType, attributeID, err := attributes.PluginAttributeFromQuery(c, n)
	if err != nil {
		err := fmt.Errorf("node: apiNodeGetAttributesValue: failed to get plugin attribute: %w", err)
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_attribute", err, n.log)
		return
	}

	allowed, err := auth.CheckAttributePermissions(
		c, *attrType.GetEntry(), n.GetNodeAttributes(), attributeID, userID,
		auth.ReadOperation)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiNodeGetAttributesValue: permissions check")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_permissions_check", err, n.log)
		return
	} else if !allowed {
		err := fmt.Errorf("operation not permitted")
		api.AbortRequest(c, http.StatusForbidden, "operation_not_permitted", err, n.log)
		return
	}

	out, ok := n.GetNodeAttributes().GetValue(attributeID)
	if !ok {
		err := errors.Errorf("Node: apiNodeGetAttributesValue: object attribute value not found: %s", attributeID)
		api.AbortRequest(c, http.StatusNotFound, "attribute_not_found", err, n.log)
		return
	}

	c.JSON(http.StatusOK, out)
}

// @Summary Set node attribute
// @Description Sets entire node attribute
// @Tags attributes,node
// @Security Bearer
// @Param body body node.apiNodeSetAttributesValue.InBody true "body params"
// @Success 202 {object} entry.AttributeValue
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/node/attributes [post]
func (n *Node) apiNodeSetAttributesValue(c *gin.Context) {
	type InBody struct {
		attributes.QueryPluginAttribute
		AttributeValue map[string]any `json:"attribute_value" binding:"required"`
	}

	inBody := InBody{}

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiNodeSetAttributesValue: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	pluginID, err := umid.Parse(inBody.PluginID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiNodeSetAttributesValue: failed to parse plugin umid")
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

	allowed, err := auth.CheckAttributePermissions(
		c, *attrType.GetEntry(), n.GetNodeAttributes(), attributeID, userID,
		auth.WriteOperation)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiNodeSetAttributesValue: permissions check")
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

	payload, err := n.GetNodeAttributes().Upsert(attributeID, modifyFn, true)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiNodeSetAttributesValue: failed to upsert object attribute")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_upsert", err, n.log)
		return
	}

	c.JSON(http.StatusAccepted, payload.Value)
}

// @Summary Delete node attribute
// @Description Deletes a node attribute
// @Tags attributes,node
// @Security Bearer
// @Param body body attributes.QueryPluginAttribute true "body params"
// @Success 200 {object} nil
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/node/attributes [delete]
func (n *Node) apiNodeRemoveAttributesValue(c *gin.Context) {
	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiNodeRemoveAttributesValue: user from context")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user", err, n.log)
		return
	}

	attrType, attributeID, err := attributes.PluginAttributeFromQuery(c, n)
	if err != nil {
		err := fmt.Errorf("node: apiNodeRemoveAttributesValue: failed to get plugin attribute: %w", err)
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_attribute", err, n.log)
		return
	}

	allowed, err := auth.CheckAttributePermissions(
		c, *attrType.GetEntry(), n.GetNodeAttributes(), attributeID, userID,
		auth.WriteOperation)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiNodeRemoveAttributesValue: permissions check")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_permissions_check", err, n.log)
		return
	} else if !allowed {
		err := fmt.Errorf("operation not permitted")
		api.AbortRequest(c, http.StatusForbidden, "operation_not_permitted", err, n.log)
		return
	}

	if _, err := n.GetNodeAttributes().UpdateValue(
		attributeID, modify.ReplaceWith[entry.AttributeValue](nil), true,
	); err != nil {
		err = errors.WithMessage(err, "Node: apiNodeRemoveAttributesValue: failed to update object attribute")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_update", err, n.log)
		return
	}

	c.JSON(http.StatusOK, nil)
}
