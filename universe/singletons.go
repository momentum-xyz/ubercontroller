package universe

import (
	"github.com/momentum-xyz/ubercontroller/logger"
)

var (
	log           = logger.L()
	nodeSingleton Node
)

func InitializeNode(node Node) {
	nodeSingleton = node
}

func GetNode() Node {
	if nodeSingleton == nil {
		log.Fatal("GetNode: node singleton is nil")
	}
	return nodeSingleton
}
