package space

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/momentum-xyz/ubercontroller"
	"github.com/momentum-xyz/ubercontroller/universe/api/middleware"
)

func (s *Space) RegisterAPI(r *gin.Engine) {
	s.log.Info("Registering api for plugins...")

	vx := r.Group(fmt.Sprintf("/api/v%d", ubercontroller.APIMajorVersion))
	{
		auth := vx.Group("", middleware.VerifyUser(s.log))

		authPlugins := auth.Group("/spaces//attributes")
		{
			authPlugins.GET("", s.apiGetPlugins)
		}
	}
}
