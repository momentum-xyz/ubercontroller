package streamchat

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/momentum-xyz/ubercontroller"
	"github.com/momentum-xyz/ubercontroller/universe/common/api/middleware"
)

func (s *StreamChat) RegisterAPI(r *gin.Engine) {
	s.log.Debug("Registering api for streamchat...")
	vx := r.Group(fmt.Sprintf("/api/v%d", ubercontroller.APIMajorVersion))
	{
		auth := vx.Group("", middleware.VerifyUser(s.log))

		authPlugins := auth.Group("/streamchat")
		{
			authPlugins.POST("/:objectID/token", s.apiChannelToken)
			authPlugins.POST("/:objectID/join", s.apiChannelJoin)
			authPlugins.POST("/:objectID/leave", s.apiChannelLeave)
		}
	}

}
