package node

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
)

type AIProvidersFlags struct {
	Leonardo     bool `json:"leonardo"`
	Blockadelabs bool `json:"blockadelabs"`
	Chatgpt      bool `json:"chatgpt"`
}

// @Summary Config for UI client
// @Schemes
// @Description Returns config for UI client
// @Tags config
// @Success 200 {object} node.apiGetUIClientConfig.Response
// @Router /api/v4/config/ui-client [get]
func (n *Node) apiGetUIClientConfig(c *gin.Context) {
	type Response struct {
		config.UIClient
		NodeID           string           `json:"NODE_ID"`
		RenderServiceURL string           `json:"RENDER_SERVICE_URL"`
		AIProvidersFlags AIProvidersFlags `json:"AI_PROVIDERS"`
	}

	blockadelabsApiKey := n.getApiKeyParameterFromNodeAttribute("blockadelabs")
	leonardoApiKey := n.getApiKeyParameterFromNodeAttribute("leonardo")
	chatgptApiKey := n.getApiKeyParameterFromNodeAttribute("open_ai")

	out := Response{
		UIClient:         n.cfg.UIClient,
		NodeID:           n.GetID().String(),
		RenderServiceURL: n.cfg.Settings.FrontendURL + "/api/v4/media/render",
		AIProvidersFlags: AIProvidersFlags{
			Leonardo:     len(leonardoApiKey) > 0,
			Blockadelabs: len(blockadelabsApiKey) > 0,
			Chatgpt:      len(chatgptApiKey) > 0 && chatgptApiKey != "set_your_api_key_here",
		},
	}

	c.JSON(http.StatusOK, out)
}

func (n *Node) getApiKeyParameterFromNodeAttribute(nodeAttributeName string) string {
	attribute, ok := n.nodeAttributes.GetValue(entry.NewAttributeID(universe.GetSystemPluginID(), nodeAttributeName))
	if attribute != nil && ok {
		apiKey := utils.GetFromAnyMap(*attribute, "api_key", "")
		return apiKey
	}
	return ""
}
