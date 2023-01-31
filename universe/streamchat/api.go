package streamchat

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/momentum-xyz/ubercontroller"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/middleware"
)

func (s *StreamChat) RegisterAPI(r *gin.Engine) {
	s.log.Debug("Registering api for streamchat...")
	vx := r.Group(fmt.Sprintf("/api/v%d", ubercontroller.APIMajorVersion))
	{
		verified := vx.Group("", middleware.VerifyUser(s.log))
		{
			streamChat := verified.Group("/streamchat")
			{
				streamChat.POST("/:spaceID/token", s.apiChannelToken)
				streamChat.POST("/:spaceID/join", s.apiChannelJoin)
				streamChat.POST("/:spaceID/leave", s.apiChannelLeave)
			}
		}
	}
}
