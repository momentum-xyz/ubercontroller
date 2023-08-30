package node

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
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
	type Body struct {
		OdysseyID *umid.UMID `json:"odyssey_id" binding:"required"`
	}

	var inBody Body
	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiNodeGetChallenge: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	if n.cfg.Common.HostingAllowAll {

	}

	id := entry.NewUserObjectID(userID, objectID)

	_, err = n.userObjects.Upsert(id, modifyFunc, true)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiPostMemberForObject: failed to upsert user_object")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}

	c.JSON(http.StatusOK, nil)
}
