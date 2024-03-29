package assets_3d

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/momentum-xyz/ubercontroller"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/middleware"
)

func (a *Assets3d) RegisterAPI(r *gin.Engine) {
	a.log.Info("Registering api for assets 3d...")

	vx := r.Group(fmt.Sprintf("/api/v%d", ubercontroller.APIMajorVersion))
	{
		assets3d := vx.Group("/assets-3d", middleware.VerifyUser(a.log))

		assets3d.GET("", a.apiGetAssets3d)
		assets3d.GET("/options", a.apiGetAssets3dOptions)
		assets3d.POST("/upload", a.apiUploadAsset3d)
		assets3d.DELETE("/:asset3dID", a.apiRemoveAsset3dByID)
		assets3d.PATCH("/:asset3dID", a.apiUpdateAsset3dByID)
	}
}
