package space

import (
	"github.com/gin-gonic/gin"
	"github.com/momentum-xyz/ubercontroller/universe/api/dto"
	"net/http"
)

func (s *Space) apiGetSpaceAttributes(c *gin.Context) {
	s.spaceAttributes.Mu.RLock()
	defer s.spaceAttributes.Mu.RUnlock()

	out := make(dto.Space, len(s.spaceAttributes.Data))

	for pluginID, plugin := range p.plugins.Data {
		meta := plugin.GetMeta()
		if meta == nil {
			p.log.Warnf("Plugins: apiGetPlugins: failed to get meta for plugin: %s", pluginID)
			continue
		}
		out[pluginID] = utils.GetFromAnyMap(*meta, "name", "")
	}

	c.JSON(http.StatusOK, out)
}
