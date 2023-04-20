package node

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/momentum-xyz/ubercontroller"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/middleware"
)

// @title        Momentum API
// @version      4.0
// @description  Momentum REST API

// @BasePath /

func (n *Node) RegisterAPI(r *gin.Engine) {
	n.log.Infof("Registering api for node: %s...", n.GetID())

	if n.cfg.Common.PProfAPI {
		registerPProfAPI(r.Group("/debug"))
	}

	if n.cfg.Common.AllowCORS {
		r.Use(cors.New(cors.Config{
			AllowOrigins: []string{"*"},
		}))
	}

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

			auth.POST("/attach-account", n.apiAttachAccount)
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
			userMe := verifiedUsers.Group("/me")
			{
				userMe.GET("", n.apiUsersGetMe)
				userMe.GET("/attributes", n.apiGetMeUserAttributeValue)

				userMe.POST("/attach-account", n.apiAttachAccount)

				userMe.GET("/stakes", n.apiGetMyStakes)
			}

			verifiedUsers.POST("/mutual-docks", n.apiUsersCreateMutualDocks)
			verifiedUsers.DELETE("/mutual-docks", n.apiUsersRemoveMutualDocks)

			verifiedUsers.GET("", n.apiUsersGet)
			verifiedUsers.GET("/search", n.apiUsersSearchUsers)

			user := verifiedUsers.Group("/:userID")
			{
				user.GET("", n.apiUsersGetByID)
				user.GET("/worlds", n.apiUsersGetOwnedWorlds)

				userAttributesGroup := user.Group("/attributes")
				{
					userAttributesGroup.GET("", n.apiGetUserAttributeValue)
					userAttributesGroup.GET("/sub", n.apiGetUserAttributeSubValue)

					userAttributesGroup.POST("", n.apiSetUserAttributeValue)
					userAttributesGroup.DELETE("", n.apiRemoveUserAttributeValue)

					userAttributesGroup.POST("/sub", n.apiSetUserAttributeSubValue)
					userAttributesGroup.DELETE("/sub", n.apiRemoveUserAttributeSubValue)
				}
			}

			userUserAttributesGroup := verifiedUsers.Group("/attributes")
			{
				userUserAttributesGroup.POST("/sub/:userID/:targetID", n.apiSetUserUserSubAttributeValue)
			}
		}

		verifiedProfile := verified.Group("/profile")
		{
			verifiedProfile.PATCH("", n.apiProfileUpdate)
			verifiedProfile.GET("/check-job/:jobID", n.apiProfileUpdateCheckJob)
		}

		verifiedObjects := verified.Group("/objects")
		{
			verifiedObjects.POST("", n.apiObjectsCreateObject)

			object := verifiedObjects.Group("/:objectID")
			{
				objectAdmin := object.Group("", middleware.AuthorizeAdmin(n.log))
				{
					objectAdmin.POST("/options/sub", n.apiObjectsSetObjectSubOption)
					objectAdmin.DELETE("/options/sub", n.apiObjectsRemoveObjectSubOption)

					objectAdmin.DELETE("", n.apiRemoveObject)
					objectAdmin.PATCH("", n.apiUpdateObject)
				}

				object.POST("/attributes", n.apiSetObjectAttributesValue)
				object.DELETE("/attributes", n.apiRemoveObjectAttributeValue)

				object.POST("/attributes/sub", n.apiSetObjectAttributeSubValue)
				object.DELETE("/attributes/sub", n.apiRemoveObjectAttributeSubValue)

				object.POST("/agora/token", n.apiGenAgoraToken)

				object.GET("", n.apiGetObject)

				object.GET("/options", n.apiObjectsGetObjectOptions)
				object.GET("/options/sub", n.apiObjectsGetObjectSubOptions)

				object.GET("/attributes", n.apiGetObjectAttributesValue)
				object.GET("/attributes-with-children", n.apiGetObjectWithChildrenAttributeValues)

				object.GET("/attributes/sub", n.apiGetObjectAttributeSubValue)

				object.GET("/all-users/attributes", n.apiGetObjectAllUsersAttributeValuesList)
			}

			objectUser := verifiedObjects.Group("/:objectID/:userID")
			{
				objectUser.POST("/attributes", n.apiSetObjectUserAttributesValue)
				objectUser.DELETE("/attributes", n.apiRemoveObjectUserAttributeValue)

				objectUser.POST("/attributes/sub", n.apiSetObjectUserAttributeSubValue)
				objectUser.DELETE("/attributes/sub", n.apiRemoveObjectUserAttributeSubValue)

				objectUser.GET("/attributes", n.apiGetObjectUserAttributesValue)

				objectUser.GET("/attributes/sub", n.apiGetObjectUserAttributeSubValue)
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
