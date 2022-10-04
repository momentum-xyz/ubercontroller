package plugins

import "github.com/gin-gonic/gin"

func (a *Plugins) RegisterAPI(r *gin.Engine) {
	a.log.Info("Registering api for plugins...")
}
