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
// @Failure 500 {object} api.HTTPError
// @Router /api/v4/config/ui-client [get]
func (n *Node) apiGetUIClientConfig(c *gin.Context) {
	type Response struct {
		config.UIClient
		NodeID             string `json:"NODE_ID"`
		RenderServiceURL   string `json:"RENDER_SERVICE_URL"`
		BackendEndpointURL string `json:"BACKEND_ENDPOINT_URL"`
	}

	out := Response{
		UIClient:           n.cfg.UIClient,
		NodeID:             n.GetID().String(),
		RenderServiceURL:   n.cfg.Settings.FrontendURL + "/api/v3/render",
		BackendEndpointURL: n.cfg.Settings.FrontendURL + "/api/v3/backend",
	}

	c.JSON(http.StatusOK, out)
}
