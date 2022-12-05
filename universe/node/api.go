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
		drive := vx.Group("/drive")
		{
			drive.POST("/mint-odyssey", n.apiDriveMintOdyssey)
			drive.GET("/mint-odyssey/check-job/:jobID", n.apiDriveMintOdysseyCheckJob)
		}

		config := vx.Group("/config")
		{
			config.GET("/ui-client", n.apiGetUIClientConfig)
		}

		auth := vx.Group("/auth")
		{
			auth.GET("/challenge", n.apiGenChallenge)
			auth.POST("/token", n.apiGenToken)

			auth.POST("/guest-token", n.apiGuestToken)
		}

		// with verified user
		verified := vx.Group("", middleware.VerifyUser(n.log))

		verifiedMedia := verified.Group("/media")
		{
			verifiedMedia.POST("/upload/image", n.apiMediaUploadImage)
		}

		verifiedUsers := verified.Group("/users")
		{
			verifiedUsers.GET("/me", n.apiUsersGetMe)

			verifiedUser := verifiedUsers.Group("/:userID")
			{
				verifiedUser.GET("", n.apiUsersGetById)
			}
		}

		verifiedProfile := verified.Group("/profile")
		{
			verifiedProfile.PATCH("", n.apiProfileUpdate)
		}

		verifiedSpaces := verified.Group("/spaces")
		{
			verifiedSpaces.POST("", n.apiCreateSpace)

			verifiedSpace := verifiedSpaces.Group("/:spaceID")
			{
				verifiedSpace.GET("", n.apiGetSpace)
				verifiedSpace.DELETE("", n.apiRemoveSpace)

				verifiedSpace.GET("/options", n.apiSpacesGetSpaceOptions)
				verifiedSpace.GET("/options/sub", n.apiSpacesGetSpaceSubOptions)
				verifiedSpace.POST("/options/sub", n.apiSpacesSetSpaceSubOption)
				verifiedSpace.DELETE("/options/sub", n.apiSpacesRemoveSpaceSubOption)

				verifiedSpace.GET("/attributes", n.apiGetSpaceAttributesValue)
				verifiedSpace.GET("/attributes-with-children", n.apiGetSpaceWithChildrenAttributeValues)
				verifiedSpace.POST("/attributes", n.apiSetSpaceAttributesValue)
				verifiedSpace.DELETE("/attributes", n.apiRemoveSpaceAttributeValue)
				verifiedSpace.GET("/attributes/sub", n.apiGetSpaceAttributeSubValue)
				verifiedSpace.POST("/attributes/sub", n.apiSetSpaceAttributeSubValue)
				verifiedSpace.DELETE("/attributes/sub", n.apiRemoveSpaceAttributeSubValue)

				verifiedAgora := verifiedSpace.Group("/agora")
				{
					verifiedAgora.POST("/token", n.apiGenAgoraToken)
				}
			}

			verifiedSpace.GET("/all-users/attributes", n.apiGetSpaceAllUsersAttributeValuesList)

			verifiedSpaceUser := verifiedSpaces.Group("/:spaceID/:userID")
			{
				verifiedSpaceUser.GET("/attributes", n.apiGetSpaceUserAttributesValue)
				verifiedSpaceUser.POST("/attributes", n.apiSetSpaceUserAttributesValue)
				verifiedSpaceUser.DELETE("/attributes", n.apiRemoveSpaceUserAttributeValue)
				verifiedSpaceUser.GET("/attributes/sub", n.apiGetSpaceUserAttributeSubValue)
				verifiedSpaceUser.POST("/attributes/sub", n.apiSetSpaceUserAttributeSubValue)
				verifiedSpaceUser.DELETE("/attributes/sub", n.apiRemoveSpaceUserAttributeSubValue)
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
