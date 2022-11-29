package node

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (n *Node) apiGetChallenge(c *gin.Context) {
	type Out struct {
		Challenge string `json:"challenge"`
	}
	out := Out{
		Challenge: "my super secret challenge",
	}

	c.JSON(http.StatusOK, out)
}

func (n *Node) apiGenToken(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}
