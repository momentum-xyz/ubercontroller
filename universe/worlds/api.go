package worlds

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/momentum-xyz/ubercontroller"
	"github.com/momentum-xyz/ubercontroller/universe/common/api/middleware"
)

func (w *Worlds) RegisterAPI(r *gin.Engine) {
	w.log.Info("Registering api for worlds...")

	vx := r.Group(fmt.Sprintf("/api/v%d", ubercontroller.APIMajorVersion))
	{
		verified := vx.Group("", middleware.VerifyUser(w.log))
		{
			// with admin rights
			authorizedAdmin := verified.Group("", middleware.AuthorizeAdmin(w.log))
			{
				// Todo: implement
			}

			// with regular rights
			authorizedUser := verified.Group("", middleware.AuthorizeUser(w.log))
			{
				// Todo: implement
			}

			worlds := verified.Group("/worlds")
			{
				world := worlds.Group("/:worldID")
				{
					world.GET("/explore", w.apiWorldsGetSpacesWithChildren)
					world.GET("/explore/search", w.apiWorldsSearchSpaces)

					world.POST("/teleport-user", w.apiWorldsTeleportUser)

					world.POST("/fly-to-me", w.apiWorldsFlyToMe)
				}
			}
		}
	}
}
