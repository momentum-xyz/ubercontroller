package plugins

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/momentum-xyz/ubercontroller"
	"github.com/momentum-xyz/ubercontroller/universe/common/api/middleware"
)

func (p *Plugins) RegisterAPI(r *gin.Engine) {
	p.log.Info("Registering api for plugins...")

	vx := r.Group(fmt.Sprintf("/api/v%d", ubercontroller.APIMajorVersion))
	{
		verified := vx.Group("", middleware.VerifyUser(p.log))
		{
			// with admin rights
			authorizedAdmin := verified.Group("", middleware.AuthorizeAdmin(p.log))
			{
				// Todo: implement
			}

			// with regular rights
			authorizedUser := verified.Group("", middleware.AuthorizeUser(p.log))
			{
				// Todo: implement
			}

			plugins := verified.Group("/plugins")
			{
				plugins.GET("", p.apiGetPlugins)
				plugins.GET("/search", p.apiSearchPlugins)
				plugins.GET("/meta", p.apiGetPluginsMeta)
				plugins.GET("/options", p.apiGetPluginsOptions)
			}
		}
	}
}
