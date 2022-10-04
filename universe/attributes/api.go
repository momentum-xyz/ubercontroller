package attributes

import "github.com/gin-gonic/gin"

func (a *Attributes) RegisterAPI(r *gin.Engine) {
	a.log.Info("Registering api for Attributes...")
}
