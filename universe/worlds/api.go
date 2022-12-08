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
				world := worlds.Group("/:spaceID")
				{
					world.GET("/explore", w.apiWorldsGetSpacesWithChildren)
					world.GET("/explore/search", w.apiWorldsSearchSpaces)
					world.GET("/online-users", w.apiGetOnlineUsers)

					authorizedAdmin := world.Group("", middleware.AuthorizeAdmin(w.log, w.db))
					{
						authorizedAdmin.POST("/fly-to-me", w.apiWorldsFlyToMe)
						authorizedAdmin.POST("/teleport-user", w.apiWorldsTeleportUser)
					}
				}
			}
		}
	}
}
