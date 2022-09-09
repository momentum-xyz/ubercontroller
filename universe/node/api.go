package node

import "github.com/gin-gonic/gin"

func (n *Node) RegisterAPI(r *gin.Engine) {
	n.log.Infof("Registering api for node: %s...", n.GetID())
}
