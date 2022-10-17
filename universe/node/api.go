package node

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/momentum-xyz/ubercontroller"
	"github.com/momentum-xyz/ubercontroller/universe/api/middleware"
)

func (n *Node) sendConfig(c *gin.Context) {
	s := "{\"KEYCLOAK_OPENID_CONNECT_URL\":\"https://dev-x5u42do.odyssey.ninja/auth/realms/Momentum\",\"KEYCLOAK_OPENID_CLIENT_ID\":\"react-client\",\"KEYCLOAK_OPENID_SCOPE\":\"openid offline_access\",\"HYDRA_OPENID_CONNECT_URL\":\"https://oidc.dev.odyssey.ninja/\",\"HYDRA_OPENID_CLIENT_ID\":\"8ad3d327-e2cf-4828-80c2-d218cf6a547d\",\"HYDRA_OPENID_GUEST_CLIENT_ID\":\"93f4f607-a56d-4689-947a-0529630167ad\",\"HYDRA_OPENID_SCOPE\":\"openid offline\",\"WEB3_IDENTITY_PROVIDER_URL\":\"https://dev.odyssey.ninja/web3-idp\",\"GUEST_IDENTITY_PROVIDER_URL\":\"https://dev.odyssey.ninja/guest-idp\",\"SENTRY_DSN\":\"https://c80e1c06d0cc495a8ead3bf782003f0b@o572058.ingest.sentry.io/6389462\",\"AGORA_APP_ID\":\"94eb4a78b1c44f488767bba62a7bda74\",\"AUTH_SERVICE_URL\":\"https://dev-c3thnss.odyssey.ninja/auth\",\"GOOGLE_API_CLIENT_ID\":\"361695877651-ugpcm6qnet4r2ub72sff2e1atmipt2mm.apps.googleusercontent.com\",\"GOOGLE_API_DEVELOPER_KEY\":\"AIzaSyCu4HkmYN2Ehf3hUHfao4hYdi-AiWLQ0m0\",\"MIRO_APP_ID\":\"3074457355834782314\",\"YOUTUBE_KEY\":\"AIzaSyBsruHhy84M2natgb-WofqGIe4sKQT8PxY\",\"STREAMCHAT_KEY\":\"bjhq75kp5xum\",\"UNITY_CLIENT_STREAMING_ASSETS_URL\":\"StreamingAssets\",\"UNITY_CLIENT_COMPANY_NAME\":\"Odyssey\",\"UNITY_CLIENT_PRODUCT_NAME\":\"Odyssey Momentum\",\"UNITY_CLIENT_PRODUCT_VERSION\":\"0.1\",\"UNITY_CLIENT_URL\":\"https://dev.odyssey.ninja/unity\",\"UNITY_CLIENT_LOADER_URL\":\"https://dev.odyssey.ninja/unity/WebGL.loader.js\",\"UNITY_CLIENT_DATA_URL\":\"https://dev.odyssey.ninja/unity/WebGL.data.gz\",\"UNITY_CLIENT_FRAMEWORK_URL\":\"https://dev.odyssey.ninja/unity/WebGL.framework.js.gz\",\"UNITY_CLIENT_CODE_URL\":\"https://dev.odyssey.ninja/unity/WebGL.wasm.gz\",\"RENDER_SERVICE_URL\":\"https://dev.odyssey.ninja/api/v3/render\",\"BACKEND_ENDPOINT_URL\":\"https://dev.odyssey.ninja/api/v3/backend\",\"BACKEND_V4_ENDPOINT_URL\":\"https://dev.odyssey.ninja/api/v4\"}"

	c.String(http.StatusOK, s)

}

func (n *Node) RegisterAPI(r *gin.Engine) {
	n.log.Infof("Registering api for node: %s...", n.GetID())

	r.GET("/version", n.apiGetVersion)
	r.GET("/health", n.apiHealthCheck)
	r.GET("/posbus", n.apiPosBusHandler)
	r.GET("/config/ui-client", n.sendConfig)

	vx := r.Group(fmt.Sprintf("/api/v%d", ubercontroller.APIMajorVersion))
	{
		config := vx.Group("/config")
		{
			config.GET("/ui-client", n.apiGetUIClientConfig)
		}

		users := vx.Group("/users")
		{
			users.GET("/check", n.apiUsersCheck)
		}

		// with auth
		auth := vx.Group("", middleware.VerifyUser(n.log))

		authUsers := auth.Group("/users")
		{
			authUsers.GET("/me", n.apiUsersGetMe)
		}

		authProfile := auth.Group("/profile")
		{
			authProfile.PUT("", n.apiProfileEdit)
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
