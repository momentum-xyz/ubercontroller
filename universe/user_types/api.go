package user_types

import "github.com/gin-gonic/gin"

func (ut *UserTypes) RegisterAPI(r *gin.Engine) {
	ut.log.Info("Registering api for space types...")
}
