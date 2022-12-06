package assets_3d

import (
	"fmt"
	"github.com/momentum-xyz/ubercontroller/universe/common/api/middleware"

	"github.com/gin-gonic/gin"

	"github.com/momentum-xyz/ubercontroller"
)

func (a *Assets3d) RegisterAPI(r *gin.Engine) {
	a.log.Info("Registering api for assets 3d...")

	vx := r.Group(fmt.Sprintf("/api/v%d", ubercontroller.APIMajorVersion))
	{
		// with auth
		auth := vx.Group("", middleware.VerifyUser(a.log))
		{
			// with admin rights
			authorizedSpecial := auth.Group("", middleware.AuthorizeSpecial(a.log))
			{
				authAssets3d := authorizedSpecial.Group("/assets-3d")
				{
					authAssets3d.POST("", a.apiAddAssets3d)
					authAssets3d.POST("/upload", a.apiUploadAsset3d)
					authAssets3d.DELETE("", a.apiRemoveAssets3dByIDs)
				}
			}

			auth.GET("", a.apiGetAssets3d)
			auth.GET("/meta", a.apiGetAssets3dMeta)
			auth.GET("/options", a.apiGetAssets3dOptions)
		}
	}
}
