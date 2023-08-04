package media

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/momentum-xyz/ubercontroller"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/middleware"
)

func (m *Media) RegisterAPI(r *gin.Engine) {
	m.log.Debug("Registering api for media...")
	vx := r.Group(fmt.Sprintf("/api/v%d", ubercontroller.APIMajorVersion))
	{
		verified := vx.Group("", middleware.VerifyUser(m.log))
		{
			media := verified.Group("/media")
			{
				//media.GET("/render/get/{file:[a-zA-Z0-9]+}",)
				//media.GET("/render/texture/{rsize:s[0-9]}/{file:[a-zA-Z0-9]+}", )
				//media.GET("/render/track/{file:[a-zA-Z0-9]+}", )
				//media.GET("/render/video/{file:[a-zA-Z0-9]+}", )
				//media.GET("/render/asset/{file:[a-zA-Z0-9]+}", )
				//
				//media.POST("/render/addimage", )
				//media.POST("/render/addframe", )
				//media.POST("/render/addtube", )
				//
				//media.POST("/addvideo", )
				//media.POST("/addtrack", )
				//media.POST("/addasset", )
				//
				//media.DELETE("/deltrack/{file:[a-zA-Z0-9]+}", )
			}
			println(media)
		}
	}
}
