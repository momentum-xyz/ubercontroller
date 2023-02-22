package node

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
)

// @Summary Config for UI client
// @Schemes
// @Description Returns config for UI client
// @Tags config
// @Accept json
// @Produce json
// @Success 200 {object} node.apiGetUIClientConfig.Response
// @Failure 500 {object} api.HTTPError
// @Router /api/v4/config/ui-client [get]
func (n *Node) apiGetUIClientConfig(c *gin.Context) {
	type Response struct {
		config.UIClient
		NodeID                  string `json:"NODE_ID"`
		UnityClientURL          string `json:"UNITY_CLIENT_URL"`
		UnityClientLoaderURL    string `json:"UNITY_CLIENT_LOADER_URL"`
		UnityClientDataURL      string `json:"UNITY_CLIENT_DATA_URL"`
		UnityClientFrameworkURL string `json:"UNITY_CLIENT_FRAMEWORK_URL"`
		UnityClientCodeURL      string `json:"UNITY_CLIENT_CODE_URL"`
		RenderServiceURL        string `json:"RENDER_SERVICE_URL"`
		BackendEndpointURL      string `json:"BACKEND_ENDPOINT_URL"`
	}

	var unityClientURLString string
	if n.cfg.UIClient.UnityClientURL != "" {
		unityClientURLString = n.cfg.UIClient.UnityClientURL
	} else {
		unityClientURLString = n.cfg.Settings.FrontendURL + "/unity"
	}
	unityClientURL, err := url.Parse(unityClientURLString)
	if err != nil {
		api.AbortRequest(c, http.StatusInternalServerError, "invalid_configuration", err, n.log)
		return
	}
	out := Response{
		UIClient:                n.cfg.UIClient,
		NodeID:                  n.GetID().String(),
		UnityClientURL:          unityClientURL.String(),
		UnityClientLoaderURL:    unityClientURL.JoinPath(n.cfg.UIClient.UnityLoaderFileName).String(),
		UnityClientDataURL:      unityClientURL.JoinPath(n.cfg.UIClient.UnityDataFileName).String(),
		UnityClientFrameworkURL: unityClientURL.JoinPath(n.cfg.UIClient.UnityFrameworkFileName).String(),
		UnityClientCodeURL:      unityClientURL.JoinPath(n.cfg.UIClient.UnityCodeFileName).String(),
		RenderServiceURL:        n.cfg.Settings.FrontendURL + "/api/v3/render",
		BackendEndpointURL:      n.cfg.Settings.FrontendURL + "/api/v3/backend",
	}

	c.JSON(http.StatusOK, out)
}
