package node

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/momentum-xyz/ubercontroller"
	"github.com/momentum-xyz/ubercontroller/universe/api/middleware"
)

// @title        Momentum API
// @version      4.0
// @description  Momentum REST API
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
			authSpace := authSpaces.Group("/:spaceID")
			{
				authSpace.POST("/options/sub", n.apiSpacesSetSpaceSubOption)
				authSpace.DELETE("/options/sub", n.apiSpacesRemoveSpaceSubOption)

				authSpace.GET("/effective-options", n.apiSpacesGetSpaceEffectiveOptions)
				authSpace.GET("/effective-options/sub", n.apiSpacesGetSpaceEffectiveSubOption)

				authSpace.GET("/attributes", n.apiGetSpaceAttributesValue)
				authSpace.GET("/attributes-with-children", n.apiGetSpaceWithChildrenAttributeValues)
				authSpace.POST("/attributes", n.apiSetSpaceAttributesValue)
				authSpace.GET("/attributes/sub", n.apiGetSpaceAttributeSubValue)
				authSpace.POST("/attributes/sub", n.apiSetSpaceAttributeSubValue)
				authSpace.DELETE("/attributes/sub", n.apiRemoveSpaceAttributeSubValue)
			}
		}
	}
}

// @Summary Version of running controller app
// @Schemes
// @Description Version of running controller app
// @Tags config
// @Accept json
// @Produce json
// @Success 200 {object} any
// @Success 500 {object} api.HTTPError
// @Router /version [get]
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
