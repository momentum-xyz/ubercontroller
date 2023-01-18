package node

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe/common/api"
)

// @Schemes
// @Description Sets a user user attribute based on UserID and TargetID
// @Tags users
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param target_id path string true "Target ID"
// @Param plugin_id path string true "Plugin ID"
// @Param attribute_name path string true "Attribute Name"
// @Param sub_attribute_key path string true "Sub Attribute Key"
// @Param body body node.apiSetUserUserSubAttributeValue.InBody true "body params"
// @Success 200 {object} entry.AttributeValue
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Router /api/v4/users/attributes/{user_id}/{target_id}/{plugin_id}/{attribute_name}/sub/{sub_attribute_key} [post]
func (n *Node) apiSetUserUserSubAttributeValue(c *gin.Context) {
	type InBody struct {
		SubAttributeValue any `json:"sub_attribute_value" binding:"required"`
	}

	inBody := InBody{}

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err := errors.WithMessage(err, "Node: apiSetUserUserSubAttributeValue: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	userID, err := uuid.Parse(c.Param("userID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetUserUserSubAttributeValue: failed to parse user id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	targetID, err := uuid.Parse(c.Param("targetID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetUserUserSubAttributeValue: failed to parse user id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_target_id", err, n.log)
		return
	}

	pluginID, err := uuid.Parse(c.Param("pluginID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetUserUserSubAttributeValue: failed to parse plugin id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_id", err, n.log)
		return
	}

	attributeName := c.Param("attributeName")
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetUserUserSubAttributeValue: failed to get attribute name from path parameters")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_attribute_name", err, n.log)
		return
	}

	subAttributeKey := c.Param("subAttributeKey")
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetUserUserSubAttributeValue: failed to get sub-attribute key from path parameters")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_sub_attribute_key", err, n.log)
		return
	}

	attributeID := entry.NewAttributeID(pluginID, attributeName)
	userUserAttributeID := entry.NewUserUserAttributeID(attributeID, userID, targetID)

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

	userUserAttribute, err := n.GetUserUserAttributes().Upsert(userUserAttributeID, modifyFn, true)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSetUserUserSubAttributeValue: failed to upsert user user attribute")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_upsert", err, n.log)
		return
	}

	c.JSON(http.StatusAccepted, userUserAttribute.Value)
}
