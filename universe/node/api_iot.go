package node

import (
	"github.com/gin-gonic/gin"
	"github.com/momentum-xyz/ubercontroller/universe/iot"
	"github.com/pkg/errors"
)

func (n *Node) apiIOTHandler(c *gin.Context) {
	ws, err := websocketUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		n.log.Error(errors.WithMessage(err, "error: socket upgrade error, aborting connection"))
		return
	}

	iRunner := iot.NewIOTWorker(ws, n.ctx)
	if iRunner != nil {
		iRunner.Run()
	}
}
