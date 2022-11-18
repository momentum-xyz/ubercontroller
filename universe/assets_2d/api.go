package assets_2d

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/momentum-xyz/ubercontroller"
	"github.com/momentum-xyz/ubercontroller/universe/common/api/middleware"
)

func (a *Assets2d) RegisterAPI(r *gin.Engine) {
	a.log.Info("Registering api for assets 2d...")

	vx := r.Group(fmt.Sprintf("/api/v%d", ubercontroller.APIMajorVersion))
	{
		// with auth
		auth := vx.Group("", middleware.VerifyUser(a.log))

		authAssets2d := auth.Group("/assets-2d")
		{
			authAsset2d := authAssets2d.Group("/:asset2dID")
			{
				authAsset2d.GET("", a.apiGetAsset2d)
			}
		}
	}
}
