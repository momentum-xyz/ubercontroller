package assets_3d

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/momentum-xyz/ubercontroller"
)

func (a *Assets3d) RegisterAPI(r *gin.Engine) {
	a.log.Info("Registering api for assets 3d...")

	vx := r.Group(fmt.Sprintf("/api/v%d", ubercontroller.APIMajorVersion))
	{
		assets3d := vx.Group("/assets-3d")
		{
			asset3d := assets3d.Group("/:asset3dID")
			{
				asset3d.GET("", a.apiGetAsset3d)

				asset3d.GET("/options", a.apiGetAsset3dOptions)

				asset3d.GET("/meta", a.apiGetAsset3dMeta)
			}

			assets3d.GET("", a.apiGetAssets3d)

			assets3d.DELETE("", a.apiRemoveAssets3dByIDs)

			assets3d.POST("", a.apiAddAssets3d)
		}
	}
}
