package node

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/momentum-xyz/ubercontroller"
	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/universe/api/middleware"
)

func (n *Node) RegisterAPI(r *gin.Engine) {
	n.log.Infof("Registering api for node: %s...", n.GetID())

	r.GET("/version", n.apiGetVersion)
	r.GET("/health", n.apiHealthCheck)
	r.GET("/posbus", n.apiPosBusHandler)

	vx := r.Group(fmt.Sprintf("/api/v%d", ubercontroller.APIMajorVersion))
	{
		users := vx.Group("/users")
		{
			users.GET("/check", n.apiUsersCheck)
		}

		auth := vx.Group("", middleware.VerifyUser(n.log))

		profile := auth.Group("/profile")
		{
			profile.PUT("/:userID", n.apiProfileEdit)
		}
	}
}

func (n *Node) apiGetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
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
	})
}

func (n *Node) apiHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (n *Node) apiGetUIConfig(c *gin.Context) {
	unityClientURL := n.cfg.UIClient.FrontendURL + "/unity"

	cfg := struct {
		config.UIClient
		UnityClientURL          string `json:"UNITY_CLIENT_URL"`
		UnityClientLoaderURL    string `json:"UNITY_CLIENT_LOADER_URL"`
		UnityClientDataURL      string `json:"UNITY_CLIENT_DATA_URL"`
		UnityClientFrameworkURL string `json:"UNITY_CLIENT_FRAMEWORK_URL"`
		UnityClientCodeURL      string `json:"UNITY_CLIENT_CODE_URL"`
		RenderServiceURL        string `json:"RENDER_SERVICE_URL"`
		BackendEndpointURL      string `json:"BACKEND_ENDPOINT_URL"`
	}{
		UIClient:                n.cfg.UIClient,
		UnityClientURL:          unityClientURL,
		UnityClientLoaderURL:    unityClientURL + "/WebGL.loader.js",
		UnityClientDataURL:      unityClientURL + "/WebGL.data.gz",
		UnityClientFrameworkURL: unityClientURL + "/WebGL.framework.js.gz",
		UnityClientCodeURL:      unityClientURL + "/WebGL.wasm.gz",
		RenderServiceURL:        n.cfg.UIClient.FrontendURL + "/api/v3/render",
		BackendEndpointURL:      n.cfg.UIClient.FrontendURL + "/api/v3/backend",
	}

	c.JSON(http.StatusOK, cfg)
}
