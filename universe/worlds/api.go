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

		verifiedWorlds := verified.Group("/worlds")
		{
			verifiedWorld := verifiedWorlds.Group("/:worldID")
			{
				verifiedWorld.GET("/online-users", w.apiGetOnlineUsers)

				verifiedWorld.GET("/explore", w.apiWorldsGetSpacesWithChildren)
				verifiedWorld.GET("/explore/search", w.apiWorldsSearchSpaces)

				verifiedWorld.POST("/teleport-user", w.apiWorldsTeleportUser)

				verifiedWorld.POST("/fly-to-me", w.apiWorldsFlyToMe)
			}
		}
	}
}
