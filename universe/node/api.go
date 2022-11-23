package node

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/momentum-xyz/ubercontroller"
	"github.com/momentum-xyz/ubercontroller/universe/common/api/middleware"
)

// @title        Momentum API
// @version      4.0
// @description  Momentum REST API

// @BasePath /

func (n *Node) RegisterAPI(r *gin.Engine) {
	n.log.Infof("Registering api for node: %s...", n.GetID())

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
			authProfile.POST("/avatar", n.apiProfileUploadAvatar)
		}

		authSpaces := auth.Group("/spaces")
		{
			authSpaces.POST("", n.apiCreateSpace)

			authSpace := authSpaces.Group("/:spaceID")
			{
				authSpace.GET("", n.apiGetSpace)
				authSpace.DELETE("", n.apiRemoveSpace)

				authSpace.GET("/options", n.apiSpacesGetSpaceOptions)
				authSpace.GET("/options/sub", n.apiSpacesGetSpaceSubOptions)
				authSpace.POST("/options/sub", n.apiSpacesSetSpaceSubOption)
				authSpace.DELETE("/options/sub", n.apiSpacesRemoveSpaceSubOption)

				authSpace.GET("/attributes", n.apiGetSpaceAttributesValue)
				authSpace.GET("/attributes-with-children", n.apiGetSpaceWithChildrenAttributeValues)
				authSpace.POST("/attributes", n.apiSetSpaceAttributesValue)
				authSpace.DELETE("/attributes", n.apiRemoveSpaceAttributeValue)
				authSpace.GET("/attributes/sub", n.apiGetSpaceAttributeSubValue)
				authSpace.POST("/attributes/sub", n.apiSetSpaceAttributeSubValue)
				authSpace.DELETE("/attributes/sub", n.apiRemoveSpaceAttributeSubValue)

				authAgora := authSpace.Group("/agora")
				{
					authAgora.POST("/token", n.apiGenAgoraToken)
				}
			}
		}
	}
}

// @Summary Get application version
// @Schemes
// @Description Returns version of running controller app
// @Tags config
// @Accept json
// @Produce json
// @Success 200 {object} node.apiGetVersion.Out
// @Router /version [get]
func (n *Node) apiGetVersion(c *gin.Context) {
	type Out struct {
		API struct {
			Major int `json:"major"`
			Minor int `json:"minor"`
			Path  int `json:"patch"`
		} `json:"api"`
		Controller struct {
			Major int    `json:"major"`
			Minor int    `json:"minor"`
			Path  int    `json:"patch"`
			Git   string `json:"git"`
		} `json:"controller"`
	}

	out := Out{
		API: struct {
			Major int `json:"major"`
			Minor int `json:"minor"`
			Path  int `json:"patch"`
		}{
			Major: ubercontroller.APIMajorVersion,
			Minor: ubercontroller.APIMinorVersion,
			Path:  ubercontroller.APIPatchVersion,
		},
		Controller: struct {
			Major int    `json:"major"`
			Minor int    `json:"minor"`
			Path  int    `json:"patch"`
			Git   string `json:"git"`
		}{
			Major: ubercontroller.ControllerMajorVersion,
			Minor: ubercontroller.ControllerMinorVersion,
			Path:  ubercontroller.ControllerPathVersion,
			Git:   ubercontroller.ControllerGitVersion,
		},
	}

	c.JSON(http.StatusOK, out)
}

// @Summary Application health check
// @Schemes
// @Description Controller application health check
// @Tags app
// @Accept json
// @Produce json
// @Success 200 {object} any
// @Router /health [get]
func (n *Node) apiHealthCheck(c *gin.Context) {
	c.JSON(
		http.StatusOK, gin.H{
			"status": "ok",
		},
	)
}
