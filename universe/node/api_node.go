package node

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

// @Summary Checks if an odyssey can be registered with a node
// @Description Checks if an odyssey can be registered with a node
// @Tags node
// @Security Bearer
// @Param body body node.apiNodeGetChallenge.Body true "body params"
// @Success 200 {object} nil
// @Failure 400 {object} api.HTTPError
// /api/v4/node/get-challenge [post]
func (n *Node) apiNodeGetChallenge(c *gin.Context) {
	if n.cfg.Common.HostingAllowAll {
		c.JSON(http.StatusOK, nil)
	}

	type Body struct {
		OdysseyID *umid.UMID `json:"odyssey_id" binding:"required"`
	}

	var inBody Body
	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiNodeGetChallenge: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	attributeID := entry.NewAttributeID(universe.GetSystemPluginID(), "hosting_allow_list")
	nodeAttributeValue, ok := n.GetNodeAttributes().GetValue(attributeID)
	if !ok {

	}

	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiNodeGetChallenge: failed to get user umid from context")
		api.AbortRequest(c, http.StatusInternalServerError, "get_user_id_failed", err, n.log)
		return
	}

	if !utils.Contains(nodeAttributeValue, userID) {
		err := errors.WithMessage(err, "Node: apiNodeGetChallenge: failed to get user umid from context")
		api.AbortRequest(c, http.StatusBadRequest, "get_user_id_failed", err, n.log)
		return
	}

	c.JSON(http.StatusOK, nil)
}
