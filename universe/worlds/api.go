package worlds

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/momentum-xyz/ubercontroller"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/middleware"
)

func (w *Worlds) RegisterAPI(r *gin.Engine) {
	w.log.Info("Registering api for worlds...")

	vx := r.Group(fmt.Sprintf("/api/v%d", ubercontroller.APIMajorVersion))
	{
		verified := vx.Group("", middleware.VerifyUser(w.log))
		{
			worlds := verified.Group("/worlds")
			{
				worlds.GET("/latest", w.apiWorldsGetLatest)
				worlds.GET("/explore/search", w.apiWorldsSearchWorlds)

				world := worlds.Group("/:objectID")
				{
					world.GET("/explore", w.apiWorldsGetObjectsWithChildren)

					world.GET("/online-users", w.apiGetOnlineUsers)

					authorizedAdmin := world.Group("", middleware.AuthorizeAdmin(w.log))
					{
						authorizedAdmin.POST("/fly-to-me", w.apiWorldsFlyToMe)
						authorizedAdmin.POST("/teleport-user", w.apiWorldsTeleportUser)
					}
				}
			}
		}
	}
}
