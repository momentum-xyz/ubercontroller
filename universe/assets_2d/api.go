package assets_2d

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/momentum-xyz/ubercontroller"
)

func (a *Assets2d) RegisterAPI(r *gin.Engine) {
	a.log.Info("Registering api for assets 2d...")

	vx := r.Group(fmt.Sprintf("/api/v%d", ubercontroller.APIMajorVersion))
	{
		assets2d := vx.Group("/assets-2d")
		{
			asset2d := assets2d.Group("/:asset2dID")
			{
				asset2d.GET("", a.apiGetAsset2d)
			}
		}
	}
}
