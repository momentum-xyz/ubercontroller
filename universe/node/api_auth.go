package node

import (
	"github.com/gin-gonic/gin"
	"github.com/momentum-xyz/ubercontroller/universe/common/api"
	"github.com/pkg/errors"
	"net/http"
)

// @Summary Generate auth challenge
// @Schemes
// @Description Returns a new generated challenge based on params
// @Tags auth
// @Accept json
// @Produce json
// @Param query query node.apiGenChallenge.InQuery true "query params"
// @Success 200 {object} node.apiGenChallenge.Out
// @Failure 400 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /api/v4/auth/challenge [get]
func (n *Node) apiGenChallenge(c *gin.Context) {
	type InQuery struct {
		Wallet string `form:"wallet" binding:"required"`
	}
	var inQuery InQuery

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Node: apiGenChallenge: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, n.log)
		return
	}

	type Out struct {
		Challenge string `json:"challenge"`
	}
	out := Out{
		Challenge: "my super secret challenge",
	}

	c.JSON(http.StatusOK, out)
}

// @Summary Generate auth token
// @Schemes
// @Description Returns a new generated token based on params
// @Tags auth
// @Accept json
// @Produce json
// @Param body body node.apiGenToken.InBody true "body params"
// @Success 200 {object} node.apiGenToken.Out
// @Failure 400 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /api/v4/auth/token [post]
func (n *Node) apiGenToken(c *gin.Context) {
	type InBody struct {
		Wallet          string `json:"wallet" binding:"required"`
		SignedChallenge string `json:"signedChallenge" binding:"required"`
	}
	var inBody InBody

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiGenToken: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	type Out struct {
		Token string `json:"token"`
	}

	out := Out{
		Token: "my super secret token",
	}

	c.JSON(http.StatusOK, out)
}
