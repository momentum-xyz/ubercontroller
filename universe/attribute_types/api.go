package attribute_types

import "github.com/gin-gonic/gin"

func (a *AttributeTypes) RegisterAPI(r *gin.Engine) {
	a.log.Info("Registering api for AttributeTypes...")
}
