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
		auth := vx.Group("", middleware.VerifyUser(w.log))

		authWorlds := auth.Group("/worlds")
		{
			authWorld := authWorlds.Group("/:worldID")
			{
				authWorld.GET("/explore", w.apiWorldsGetSpacesWithChildren)
				authWorld.GET("/explore/search", w.apiWorldsSearchSpaces)

				authWorld.POST("/teleport-user", w.apiWorldsTeleportUser)

				authWorld.POST("/fly-to-me", w.apiWorldsFlyToMe)
			}
		}
	}
}
