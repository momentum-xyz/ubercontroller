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
		assets3d := vx.Group("", middleware.VerifyUser(a.log))
		{
			authAssets3d := assets3d.Group("/assets-3d")
			{
				authAssets3d.POST("", a.apiAddAssets3d)
				authAssets3d.POST("/upload", a.apiUploadAsset3d)
				authAssets3d.DELETE("", a.apiRemoveAssets3dByIDs)
			}

			assets3d.GET("", a.apiGetAssets3d)
			assets3d.GET("/meta", a.apiGetAssets3dMeta)
			assets3d.GET("/options", a.apiGetAssets3dOptions)
		}
	}
}
