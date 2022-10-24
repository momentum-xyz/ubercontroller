package plugins

import (
	"fmt"
	
	"github.com/gin-gonic/gin"

	"github.com/momentum-xyz/ubercontroller"
	"github.com/momentum-xyz/ubercontroller/universe/api/middleware"
)

func (p *Plugins) RegisterAPI(r *gin.Engine) {
	p.log.Info("Registering api for plugins...")

	vx := r.Group(fmt.Sprintf("/api/v%d", ubercontroller.APIMajorVersion))
	{
		auth := vx.Group("", middleware.VerifyUser(p.log))

		authPlugins := auth.Group("/plugins")
		{
			authPlugins.GET("", p.apiGetPlugins)
			authPlugins.GET("/meta", p.apiGetPluginsMeta)
			authPlugins.GET("/options", p.apiGetPluginsOptions)
		}
	}
}
