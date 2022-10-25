package node

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/momentum-xyz/ubercontroller"
	_ "github.com/momentum-xyz/ubercontroller/docs"
	"github.com/momentum-xyz/ubercontroller/universe/api/middleware"
)

func (n *Node) RegisterAPI(r *gin.Engine) {
	n.log.Infof("Registering api for node: %s...", n.GetID())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/version", n.apiGetVersion)
	r.GET("/health", n.apiHealthCheck)
	r.GET("/posbus", n.apiPosBusHandler)

	vx := r.Group(fmt.Sprintf("/api/v%d", ubercontroller.APIMajorVersion))
	{
		config := vx.Group("/config")
		{
			config.GET("/ui-client", n.apiGetUIClientConfig)
		}

		users := vx.Group("/users")
		{
			users.POST("/check", n.apiUsersCheck)
		}

		// with auth
		auth := vx.Group("", middleware.VerifyUser(n.log))

		authUsers := auth.Group("/users")
		{
			authUsers.GET("/me", n.apiUsersGetMe)
		}

		authProfile := auth.Group("/profile")
		{
			authProfile.PATCH("", n.apiProfileUpdate)
		}
	}
}

func (n *Node) apiGetVersion(c *gin.Context) {
	c.JSON(
		http.StatusOK, gin.H{
			"api": gin.H{
				"major": ubercontroller.APIMajorVersion,
				"minor": ubercontroller.APIMinorVersion,
				"patch": ubercontroller.APIPatchVersion,
			},
			"controller": gin.H{
				"major": ubercontroller.ControllerMajorVersion,
				"minor": ubercontroller.ControllerMinorVersion,
				"patch": ubercontroller.ControllerPathVersion,
				"git":   ubercontroller.ControllerGitVersion,
			},
		},
	)
}

func (n *Node) apiHealthCheck(c *gin.Context) {
	c.JSON(
		http.StatusOK, gin.H{
			"status": "ok",
		},
	)
}
