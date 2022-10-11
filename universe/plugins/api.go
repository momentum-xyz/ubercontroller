package plugins

import "github.com/gin-gonic/gin"

func (p *Plugins) RegisterAPI(r *gin.Engine) {
	p.log.Info("Registering api for plugins...")
}
