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
			worlds := verified.Group("/worlds")
			{
				world := worlds.Group("/:worldID")
				{
					world.GET("/explore", w.apiWorldsGetSpacesWithChildren)
					world.GET("/explore/search", w.apiWorldsSearchSpaces)

					// with special rights
					authorizedSpecial := world.Group("", middleware.AuthorizeSpecial(w.log, w.db))
					{
						authorizedSpecial.POST("/teleport-user", w.apiWorldsTeleportUser)

						authorizedSpecial.POST("/fly-to-me", w.apiWorldsFlyToMe)
					}
				}
			}
		}
	}
}
