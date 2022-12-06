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
		verified := vx.Group("", middleware.VerifyUser(a.log))
		{
			// with admin rights
			authorizedAdmin := verified.Group("", middleware.AuthorizeAdmin(a.log))
			{
				// Todo: implement
			}

			// with regular rights
			authorizedUser := verified.Group("", middleware.AuthorizeUser(a.log))
			{
				// Todo: implement
			}

			assets2d := verified.Group("/assets-2d")
			{
				asset2d := assets2d.Group("/:asset2dID")
				{
					asset2d.GET("", a.apiGetAsset2d)
				}
			}
		}

	}
}
