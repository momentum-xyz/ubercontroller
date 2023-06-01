package activities

import (
	"github.com/gin-gonic/gin"
)

func (a *Activities) RegisterAPI(r *gin.Engine) {
	a.log.Info("Registering api for activities...")
}
