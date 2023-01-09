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
	r.GET("/iot", n.apiIOTHandler)

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

			verifiedUsers.POST("/mutual-docks", n.apiUsersCreateMutualDocks)
			verifiedUsers.DELETE("/mutual-docks", n.apiUsersRemoveMutualDocks)

			verifiedUser := verifiedUsers.Group("/:userID")
			{
				verifiedUser.GET("", n.apiUsersGetByID)
			}
			attributes := verifiedUsers.Group("/attributes")
			{
				attributes.POST("/sub/:userID/:targetID", n.apiSetUserUserSubAttributeValue)
			}
		}

		verifiedProfile := verified.Group("/profile")
		{
			verifiedProfile.PATCH("", n.apiProfileUpdate)
		}

		verifiedSpaces := verified.Group("/spaces")
		{
			// TODO: it was created only for tests, fix or remove
			verifiedSpaces.POST("/template", n.apiSpacesCreateSpaceFromTemplate)

			verifiedSpaces.POST("", n.apiSpacesCreateSpace)

			space := verifiedSpaces.Group("/:spaceID")
			{
				spaceAttributes := space.Group("/attributes/:pluginID/:attributeName")
				{
					spaceAttributes.GET("", n.apiGetSpaceAttributesValue)
					spaceAttributes.POST("", n.apiSetSpaceAttributesValue)
					spaceAttributes.DELETE("", n.apiRemoveSpaceAttributeValue)

					spaceAttributes.GET("/with-children", n.apiGetSpaceWithChildrenAttributeValues)
					spaceAttributes.GET("/all-users", n.apiGetSpaceAllUsersAttributeValuesList)

					spaceSubAttributes := spaceAttributes.Group("/sub/:subAttributeKey")
					{
						spaceSubAttributes.GET("", n.apiGetSpaceAttributeSubValue)
						spaceSubAttributes.POST("", n.apiSetSpaceAttributeSubValue)
						spaceSubAttributes.DELETE("", n.apiRemoveSpaceAttributeSubValue)
					}
				}

				spaceAdmin := space.Group("", middleware.AuthorizeAdmin(n.log, n.db))
				{
					spaceAdmin.POST("/options/sub", n.apiSpacesSetSpaceSubOption)
					spaceAdmin.DELETE("/options/sub", n.apiSpacesRemoveSpaceSubOption)

					spaceAdmin.DELETE("", n.apiRemoveSpace)
					spaceAdmin.PATCH("", n.apiUpdateSpace)
				}

				space.POST("/agora/token", n.apiGenAgoraToken)

				space.GET("", n.apiGetSpace)

				space.GET("/options", n.apiSpacesGetSpaceOptions)
				space.GET("/options/sub", n.apiSpacesGetSpaceSubOptions)
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
			newsfeed.POST("", n.apiNewsFeedAddItem)
			newsfeed.GET("", n.apiNewsFeedGetAll)
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
