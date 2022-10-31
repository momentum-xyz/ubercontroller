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
			assets3d.GET("/:asset3dID", a.apiGetAsset3d)

			assets3d.GET("", a.apiGetAssets3d)

			assets3d.POST("/create-asset-3d", a.apiCreateAsset3d)

			assets3d.POST("/add-asset-3d", a.apiAddAssets3d)

			assets3d.DELETE("/remove-asset-3d/:asset3dID", a.apiRemoveAsset3d)

			assets3d.DELETE("/remove-assets-3d", a.apiRemoveAssets3d)

			assets3d.DELETE("/remove-assets-3d-ids", a.apiRemoveAssets3dByIDs)

			assets3d.POST("/add-assets-3d", a.apiAddAssets3d)
		}
	}
}
