package user_types

import "github.com/gin-gonic/gin"

func (u *UserTypes) RegisterAPI(r *gin.Engine) {
	u.log.Info("Registering api for object types...")
}
