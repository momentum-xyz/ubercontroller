package node

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/momentum-xyz/ubercontroller/config"
)

// @Summary Config for UI client
// @Schemes
// @Description Returns config for UI client
// @Tags config
// @Accept json
// @Produce json
// @Success 200 {object} node.apiGetUIClientConfig.Response
// @Success 500 {object} api.HTTPError
// @Router /api/v4/config/ui-client [get]
func (n *Node) apiGetUIClientConfig(c *gin.Context) {
	unityClientURL := n.cfg.UIClient.FrontendURL + "/unity"

	type Response struct {
		config.UIClient
		UnityClientURL          string `json:"UNITY_CLIENT_URL"`
		UnityClientLoaderURL    string `json:"UNITY_CLIENT_LOADER_URL"`
		UnityClientDataURL      string `json:"UNITY_CLIENT_DATA_URL"`
		UnityClientFrameworkURL string `json:"UNITY_CLIENT_FRAMEWORK_URL"`
		UnityClientCodeURL      string `json:"UNITY_CLIENT_CODE_URL"`
		RenderServiceURL        string `json:"RENDER_SERVICE_URL"`
		BackendEndpointURL      string `json:"BACKEND_ENDPOINT_URL"`
	}

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
