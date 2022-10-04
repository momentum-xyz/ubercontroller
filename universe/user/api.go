package user

import "github.com/gin-gonic/gin"

func (u *User) RegisterAPI(r *gin.Engine) {
	u.log.Info("Registering api for space types...")
}
