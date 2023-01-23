package object_types

import "github.com/gin-gonic/gin"

func (ot *ObjectTypes) RegisterAPI(r *gin.Engine) {
	ot.log.Info("Registering api for object types...")
}
