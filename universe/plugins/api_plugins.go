package plugins

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe/api"
	"github.com/momentum-xyz/ubercontroller/universe/api/dto"
	"github.com/momentum-xyz/ubercontroller/utils"
)

func (p *Plugins) apiGetPlugins(c *gin.Context) {
	p.plugins.Mu.RLock()
	defer p.plugins.Mu.RUnlock()

	out := make(dto.Plugins, len(p.plugins.Data))

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

func (p *Plugins) apiGetPluginsMeta(c *gin.Context) {
	inQuery := struct {
		PluginUUIDs []string `form:"plugin_uuids[]" binding:"required"`
	}{}

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Plugins: apiGetPluginsMeta: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, p.log)
		return
	}

	out := make(map[uuid.UUID]*dto.PluginMeta)

	for _, id := range inQuery.PluginUUIDs {
		pluginID, err := uuid.Parse(id)
		if err != nil {
			err := errors.WithMessagef(err, "Plugins: apiGetPluginsMeta: failed to parse uuid: %s", id)
			api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_uuid", err, p.log)
			return
		}

		plugin, ok := p.GetPlugin(pluginID)
		if !ok {
			err := errors.Errorf("Plugins: apiGetPluginsMeta: failed to get plugin by id: %s", pluginID)
			api.AbortRequest(c, http.StatusNotFound, "plugin_not_found", err, p.log)
			return
		}

		out[pluginID] = (*dto.PluginMeta)(plugin.GetMeta())
	}

	c.JSON(http.StatusOK, out)
}

func (p *Plugins) apiGetPluginsOptions(c *gin.Context) {
	panic("implement me")
}
