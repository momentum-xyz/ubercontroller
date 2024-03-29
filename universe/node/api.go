package node

import (
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"

	"github.com/gin-gonic/gin"
	servefiles "github.com/rickb777/servefiles/v3"

	"github.com/momentum-xyz/ubercontroller"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/middleware"
)

// @title        Momentum API
// @version      4.0
// @description  Momentum REST API

// @BasePath /
// @accept json
// @produce json
// @schemes https

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Authorization header with "Bearer" followed by a space and JWT token.

// @tag.name auth
// @tag.name users
// @tag.name profile
// @tag.name worlds
// @tag.name objects
// @tag.name members
// @tag.name media
// @tag.name assets2d
// @tag.name assets3d
// @tag.name newsfeed
// @tag.name timeline
// @tag.name plugins
// @tag.name config
// @tag.name app

func (n *Node) RegisterAPI(r *gin.Engine) {
	n.log.Infof("Registering api for node: %s...", n.GetID())

	if n.cfg.Common.PProfAPI {
		registerPProfAPI(r.Group("/debug"))
	}

	r.GET("/version", n.apiGetVersion)
	r.GET("/health", n.apiHealthCheck)
	r.GET("/posbus", n.apiPosBusHandler)
	r.GET("/iot", n.apiIOTHandler)

	vx := r.Group(fmt.Sprintf("/api/v%d", ubercontroller.APIMajorVersion))
	{
		vx.GET("/plugin", n.apiGetPluginsList)

		webhook := vx.Group("/webhook")
		{
			webhook.POST("/skybox-blockadelabs", n.apiPostSkyboxWebHook)
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

		media := vx.Group("/media")
		{
			media.GET("/render/get/:file", n.apiMediaGetImage)
			media.GET("/render/texture/:rsize/:file", n.apiMediaGetTexture)
			media.GET("/render/asset/:file", n.apiMediaGetAsset)

			media.POST("/upload/image", n.apiMediaUploadImage)

			// media.GET("/get/:file", n.apiMediaGetPlugin)
			media.Static("/render/plugin/", n.CFG.Media.Pluginpath)
			media.POST("/upload/plugin", middleware.VerifyUser(n.log), middleware.AuthorizeNodeAdmin(n.log), n.apiMediaUploadPlugin)

			media.GET("/render/video/:file", n.apiMediaGetVideo)
			media.POST("/upload/video", n.apiMediaUploadVideo)

			media.GET("/render/track/:file", n.apiMediaGetAudio)
			media.POST("/upload/audio", n.apiMediaUploadAudio)
			media.DELETE("/deltrack/:file", n.apiMediaDeleteAudio)
		}

		verified.GET("/skybox/styles", n.apiGetSkyboxStyles)
		verified.POST("/skybox/generate", n.apiPostSkyboxGenerate)
		verified.DELETE("/skybox/:skyboxID", n.apiRemoveSkyboxByID)
		verified.GET("/skybox/:skyboxID", n.apiGetSkyboxByID)

		verified.GET("/leonardo/generate/:leonardoID", n.apiGetImageGeneration)
		verified.POST("/leonardo/generate", n.apiPostImageGenerationID)

		verifiedUsers := verified.Group("/users")
		{
			userMe := verifiedUsers.Group("/me")
			{
				userMe.GET("", n.apiUsersGetMe)
				userMe.GET("/attributes", n.apiGetMeUserAttributeValue)

				userMe.POST("/attach-account", n.apiAttachAccount)
				userMe.DELETE("/remove-wallet", n.apiDeleteWallet)

				userMe.GET("/stakes", n.apiGetMyStakes)
				userMe.POST("/stakes", n.apiAddPendingStakeTransaction)

				userMe.GET("/wallets", n.apiGetMyWallets)
			}

			verifiedUsers.POST("/mutual-docks", n.apiUsersCreateMutualDocks)
			verifiedUsers.DELETE("/mutual-docks", n.apiUsersRemoveMutualDocks)

			verifiedUsers.GET("", n.apiUsersGet)
			verifiedUsers.GET("/search", n.apiUsersSearchUsers)
			verifiedUsers.GET("/top-stakers", n.apiUsersTopStakers)

			user := verifiedUsers.Group("/:userID")
			{
				user.GET("", n.apiUsersGetByID)
				user.GET("/worlds", n.apiUsersGetOwnedWorlds)
				user.GET("/staked-worlds", n.apiUsersGetStakedWorlds)

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
		}

		verifiedCanvas := verified.Group("/canvas")
		{
			verifiedCanvas.GET("/:objectID/user-contributions", n.apiCanvasGetUserContributions)
		}

		verifiedNode := verified.Group("/node")
		{
			verifiedNode.POST("/get-challenge", n.apiNodeGetChallenge)

			verifiedNode.GET("/attributes", n.apiNodeGetAttributesValue)

			verifiedNode.POST("/attributes", n.apiNodeSetAttributesValue)
			verifiedNode.DELETE("/attributes", n.apiNodeRemoveAttributesValue)

			verifiedNode.GET("/hosting-allow-list", middleware.AuthorizeNodeAdmin(n.log), n.apiGetHostingAllowList)
			verifiedNode.POST("/hosting-allow-list", middleware.AuthorizeNodeAdmin(n.log), n.apiPostItemForHostingAllowList)
			verifiedNode.DELETE("/hosting-allow-list/:userID", middleware.AuthorizeNodeAdmin(n.log), n.apiDeleteItemFromHostingAllowList)

			verifiedNode.POST("/activate-plugin", middleware.AuthorizeNodeAdmin(n.log), n.apiNodeActivatePlugin)
		}

		verifiedObjects := verified.Group("/objects")
		{
			verifiedObjects.POST("", n.apiObjectsCreateObject)

			newsfeed := verifiedObjects.Group("/newsfeed")
			{
				newsfeed.GET("", n.apiNewsfeedOverview)
			}

			object := verifiedObjects.Group("/:objectID")
			{
				object.POST("/claim-and-customise", n.apiClaimAndCustomise)
				object.POST("/unclaim-and-clear-customisation", n.apiUnclaimAndClearCustomisation)

				object.POST("/spawn-by-user", n.apiSpawnByUser)

				object.GET("/tree", n.apiGetObjectsTree)

				objectAdmin := object.Group("", middleware.AuthorizeAdmin(n.log))
				{
					objectAdmin.POST("/attributes/publicize", n.apiSetObjectAttributesPublic)

					objectAdmin.POST("/options", n.apiObjectsSetObjectOption)
					objectAdmin.POST("/options/sub", n.apiObjectsSetObjectSubOption)
					objectAdmin.DELETE("/options/sub", n.apiObjectsRemoveObjectSubOption)

					objectAdmin.DELETE("", n.apiRemoveObject)
					objectAdmin.PATCH("", n.apiUpdateObject)

					objectAdmin.POST("/clone", n.apiCloneObject)
				}

				object.POST("/attributes", n.apiSetObjectAttributesValue)
				object.DELETE("/attributes", n.apiRemoveObjectAttributeValue)

				object.POST("/attributes/sub", n.apiSetObjectAttributeSubValue)
				object.DELETE("/attributes/sub", n.apiRemoveObjectAttributeSubValue)

				if n.CFG.UIClient.AgoraAppID != "" {
					object.POST("/agora/token", n.apiGenAgoraToken)
				}

				object.GET("", n.apiGetObject)

				object.GET("/options", n.apiObjectsGetObjectOptions)
				object.GET("/options/sub", n.apiObjectsGetObjectSubOptions)

				object.GET("/attributes", n.apiGetObjectAttributesValue)
				object.GET("/attributes-with-children", n.apiGetObjectWithChildrenAttributeValues)

				object.GET("/attributes/sub", n.apiGetObjectAttributeSubValue)

				object.GET("/all-users/attributes", n.apiGetObjectAllUsersAttributeValuesList)
				object.GET("/all-users/count", n.apiGetObjectUserAttributeCount)
				object.GET("/all-users/attributes/:pluginID/:attrName/entries", n.apiObjectUserAttributeValueEntries)

				timeline := object.Group("/timeline")
				{
					timeline.GET("", n.apiTimelineForObject)
					timeline.POST("", n.apiTimelineAddForObject)

					pd := timeline.Group("/:activityID")
					{
						pd.GET("", n.apiTimelineForObjectById)
						pd.PATCH("", n.apiTimelineEditForObject)
						pd.DELETE("", n.apiTimelineRemoveForObject)
					}
				}

				members := object.Group("/members", middleware.AuthorizeAdmin(n.log))
				{
					members.GET("", n.apiMembersGetForObject)
					members.POST("", n.apiPostMemberForObject)
					members.DELETE(":userID", n.apiDeleteMemberFromObject)
				}
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
	}
	feSrvPath := n.cfg.Settings.FrontendServeDir
	if feSrvPath != "" {
		r.GET("/", func(c *gin.Context) {
			http.ServeFile(c.Writer, c.Request, filepath.Join(feSrvPath, "index.html"))
		})
		p := regexp.MustCompile("^/(api|static)")
		staticHandler := servefiles.NewAssetHandler(feSrvPath).
			WithNotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// 'SPA' fallback
				if !p.Match([]byte(r.RequestURI)) {
					http.ServeFile(w, r, filepath.Join(feSrvPath, "index.html"))
				} else {
					http.NotFound(w, r)
				}
			}))
		r.NoRoute(func(c *gin.Context) {
			staticHandler.ServeHTTP(c.Writer, c.Request)
		})
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
