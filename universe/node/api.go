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
			drive.GET("/wallet-meta", n.apiGetWalletMeta)

			drive.POST("/mint-odyssey", n.apiDriveMintOdyssey)
			drive.GET("/mint-odyssey/check-job/:jobID", n.apiDriveMintOdysseyCheckJob)
			drive.GET("/resolve-node", n.apiResolveNode)
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

		verified := vx.Group("", middleware.VerifyUser(n.log))

		//verifiedMedia := verified.Group("/media")
		//{
		//	verifiedMedia.POST("/upload/image", n.apiMediaUploadImage)
		//}

		vx.POST("/media/upload/image", n.apiMediaUploadImage)

		verifiedUsers := verified.Group("/users")
		{
			verifiedUsers.GET("/me", n.apiUsersGetMe)
			verifiedUsers.POST("/mutual-docks", n.apiUsersMutualDocks)

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

			space := verifiedSpaces.Group("/:spaceID")
			{
				authorizedAdmin := auth.Group("", middleware.AuthorizeAdmin(n.log, n.db))
				{
					space.POST("/options/sub", n.apiSpacesSetSpaceSubOption)
					space.DELETE("/options/sub", n.apiSpacesRemoveSpaceSubOption)

					authorizedAdmin.DELETE("", n.apiRemoveSpace)
				}

				space.POST("/attributes", n.apiSetSpaceAttributesValue)
				space.DELETE("/attributes", n.apiRemoveSpaceAttributeValue)

				space.POST("/attributes/sub", n.apiSetSpaceAttributeSubValue)
				space.DELETE("/attributes/sub", n.apiRemoveSpaceAttributeSubValue)

				space.POST("/agora/token", n.apiGenAgoraToken)

				space.GET("", n.apiGetSpace)

				space.GET("/options", n.apiSpacesGetSpaceOptions)
				space.GET("/options/sub", n.apiSpacesGetSpaceSubOptions)

				space.GET("/attributes", n.apiGetSpaceAttributesValue)
				space.GET("/attributes-with-children", n.apiGetSpaceWithChildrenAttributeValues)

				space.GET("/attributes/sub", n.apiGetSpaceAttributeSubValue)

				space.GET("/all-users/attributes", n.apiGetSpaceAllUsersAttributeValuesList)
			}

			spaceUser := verifiedSpaces.Group("/:spaceID/:userID")
			{
				spaceUser.POST("/attributes", n.apiSetSpaceUserAttributesValue)
				spaceUser.DELETE("/attributes", n.apiRemoveSpaceUserAttributeValue)

				spaceUser.POST("/attributes/sub", n.apiSetSpaceUserAttributeSubValue)
				spaceUser.DELETE("/attributes/sub", n.apiRemoveSpaceUserAttributeSubValue)

				spaceUser.GET("/attributes", n.apiGetSpaceUserAttributesValue)

				spaceUser.GET("/attributes/sub", n.apiGetSpaceUserAttributeSubValue)
			}
		}

		newsfeed := vx.Group("/newsfeed")
		{
			newsfeed.GET("", n.apiNewsFeed)
		}

		notifications := vx.Group("/notifications")
		{
			notifications.GET("", n.apiNotifications)
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
