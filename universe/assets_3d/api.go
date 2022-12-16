package assets_3d

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/momentum-xyz/ubercontroller"
	"github.com/momentum-xyz/ubercontroller/universe/common/api/middleware"
)

func (a *Assets3d) RegisterAPI(r *gin.Engine) {
	a.log.Info("Registering api for assets 3d...")

	vx := r.Group(fmt.Sprintf("/api/v%d", ubercontroller.APIMajorVersion))
	{
		assets3d := vx.Group("/assets-3d/:spaceID", middleware.VerifyUser(a.log))
		{
			authorizedAdmin := assets3d.Group("", middleware.AuthorizeAdmin(a.log, a.db))
			{
				authorizedAdmin.POST("", a.apiAddAssets3d)
				authorizedAdmin.POST("/upload", a.apiUploadAsset3d)
				authorizedAdmin.DELETE("", a.apiRemoveAssets3dByIDs)
			}

			assets3d.GET("", a.apiGetAssets3d)
			assets3d.GET("/meta", a.apiGetAssets3dMeta)
			assets3d.GET("/options", a.apiGetAssets3dOptions)
			assets3d.DELETE("/:asset3dID", a.apiRemoveAsset3dByID)
		}
	}
}
