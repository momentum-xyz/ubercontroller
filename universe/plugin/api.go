package plugin

import "github.com/gin-gonic/gin"

func (p *Plugin) RegisterAPI(r *gin.Engine) {
	p.log.Info("Registering api for plugin...")
}
