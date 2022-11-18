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

		authAssets3d := auth.Group("/assets-3d")
		{
			authAssets3d.GET("", a.apiGetAssets3d)
			authAssets3d.GET("/meta", a.apiGetAssets3dMeta)
			authAssets3d.GET("/options", a.apiGetAssets3dOptions)
			authAssets3d.POST("", a.apiAddAssets3d)
			authAssets3d.DELETE("", a.apiRemoveAssets3dByIDs)
		}
	}
}
