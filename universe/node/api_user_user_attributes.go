package node

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

// @Schemes
// @Description Sets a user user attribute based on UserID and TargetID
// @Tags users
// @Accept json
// @Produce json
// @Param body body node.apiSetUserUserSubAttributeValue.InBody true "body params"
// @Success 200 {object} entry.AttributeValue
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Router /api/v4/users/attributes/sub/{user_id}/{target_id} [post]
func (n *Node) apiSetUserUserSubAttributeValue(c *gin.Context) {
	type InBody struct {
		PluginID          string `json:"plugin_id" binding:"required"`
		AttributeName     string `json:"attribute_name" binding:"required"`
		SubAttributeKey   string `json:"sub_attribute_key" binding:"required"`
		SubAttributeValue any    `json:"sub_attribute_value" binding:"required"`
	}

	inBody := InBody{}

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err := errors.WithMessage(err, "Node: apiSetUserUserSubAttributeValue: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	sourceID, err := umid.Parse(c.Param("userID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetUserUserSubAttributeValue: failed to parse user umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	targetID, err := umid.Parse(c.Param("targetID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetUserUserSubAttributeValue: failed to parse user umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_target_id", err, n.log)
		return
	}

	pluginID, err := umid.Parse(inBody.PluginID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetUserUserSubAttributeValue: failed to parse plugin umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}

	result, err := n.AssessPermissions(c, pluginID, inBody.AttributeName, sourceID, WriteOperation, UserUserAttribute)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetUserUserSubAttributeValue: failed to assess permissions")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_assess_permissions", err, n.log)
		return
	}

	if !result {
		err := errors.WithMessage(err, "Node: apiSetUserUserSubAttributeValue: operation not permitted")
		api.AbortRequest(c, http.StatusForbidden, "operation_not_permitted", err, n.log)
		return
	}

	attributeID := entry.NewAttributeID(pluginID, inBody.AttributeName)
	userUserAttributeID := entry.NewUserUserAttributeID(attributeID, sourceID, targetID)

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

	userUserAttribute, err := n.GetUserUserAttributes().Upsert(userUserAttributeID, modifyFn, true)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetUserUserSubAttributeValue: failed to upsert user user attribute")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_upsert", err, n.log)
		return
	}

	c.JSON(http.StatusAccepted, userUserAttribute.Value)
}
