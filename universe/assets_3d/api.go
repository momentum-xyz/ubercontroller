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

				asset3d.DELETE("", a.ApiRemoveAsset3d)

				asset3d.POST("", a.apiCreateAsset3d)
			}

			assets3d.GET("", a.apiGetAssets3d)

			assets3d.POST("/add", a.apiAddAssets3d)

			assets3d.DELETE("/remove-assets-3d-ids", a.apiRemoveAssets3dByIDs)

			assets3d.POST("/add-assets-3d", a.apiAddAssets3d)
		}
	}
}
