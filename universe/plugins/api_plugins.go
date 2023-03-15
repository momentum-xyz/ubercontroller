package plugins

import (
	"github.com/momentum-xyz/ubercontroller/utils/mid"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/dto"
	"github.com/momentum-xyz/ubercontroller/utils"
)

// @Summary Get plugins
// @Schemes
// @Description Returns a list of plugins filtered by parameters
// @Tags plugins
// @Accept json
// @Produce json
// @Param query query plugins.apiGetPlugins.Query false "query params"
// @Success 200 {object} dto.Plugins
// @Failure 400 {object} api.HTTPError
// @Router /api/v4/plugins [get]
func (p *Plugins) apiGetPlugins(c *gin.Context) {
	type Query struct {
		IDs  []string `form:"ids[]"`
		Type string   `form:"type"`
	}

	var inQuery Query
	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err = errors.WithMessage(err, "Plugins: apiGetPlugins: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, p.log)
		return
	}

	ids := make([]mid.ID, len(inQuery.IDs))
	for i := range inQuery.IDs {
		id, err := mid.Parse(inQuery.IDs[i])
		if err != nil {
			err = errors.WithMessagef(err, "Plugins: apiGetPlugins: failed to parse mid: %s", inQuery.IDs[i])
			api.AbortRequest(c, http.StatusBadRequest, "failed_to_parse_id", err, p.log)
			return
		}
		ids[i] = id
	}

	filterFn := func(pluginID mid.ID, plugin universe.Plugin) bool {
		if len(ids) > 0 {
			var found bool
			for i := range ids {
				if pluginID == ids[i] {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}

		if inQuery.Type != "" {
			metaType := utils.GetFromAnyMap(plugin.GetMeta(), "type", "")
			if metaType != inQuery.Type {
				return false
			}
		}

		return true
	}

	var plugins map[mid.ID]universe.Plugin
	if len(inQuery.IDs) == 0 && inQuery.Type == "" {
		plugins = p.GetPlugins()
	} else {
		plugins = p.FilterPlugins(filterFn)
	}

	out := make(dto.Plugins, len(plugins))
	for _, plugin := range plugins {
		name := utils.GetFromAnyMap(plugin.GetMeta(), "name", "")
		out[plugin.GetID()] = name
	}

	c.JSON(http.StatusOK, out)
}

// @Summary Search for plugins
// @Schemes
// @Description Returns a list of plugins filtered by parameters
// @Tags plugins
// @Accept json
// @Produce json
// @Param query query plugins.apiSearchPlugins.InQuery false "query params"
// @Success 200 {object} dto.Plugins
// @Failure 400 {object} api.HTTPError
// @Router /api/v4/plugins/search [get]
func (p *Plugins) apiSearchPlugins(c *gin.Context) {
	type InQuery struct {
		Name        string `form:"name"`
		Type        string `form:"type"`
		Description string `form:"description"`
	}

	var inQuery InQuery
	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err = errors.WithMessage(err, "Plugins: apiSearchPlugins: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, p.log)
		return
	}

	filterFn := func(pluginID mid.ID, plugin universe.Plugin) bool {
		meta := plugin.GetMeta()
		if meta == nil {
			return false
		}

		if inQuery.Type != "" {
			metaType := utils.GetFromAnyMap(meta, "type", "")
			if metaType != inQuery.Type {
				return false
			}
		}

		if inQuery.Name != "" {
			metaName := utils.GetFromAnyMap(meta, "name", "")
			if metaName != "" &&
				strings.Contains(strings.ToLower(metaName), strings.ToLower(inQuery.Name)) {
				return true
			}
		}

		if inQuery.Description != "" {
			metaDescription := utils.GetFromAnyMap(meta, "description", "")
			if metaDescription != "" &&
				strings.Contains(strings.ToLower(metaDescription), strings.ToLower(inQuery.Description)) {
				return true
			}
		}

		return inQuery.Name == "" && inQuery.Description == ""
	}

	var plugins map[mid.ID]universe.Plugin
	if inQuery.Name == "" && inQuery.Type == "" && inQuery.Description == "" {
		plugins = p.GetPlugins()
	} else {
		plugins = p.FilterPlugins(filterFn)
	}

	out := make(dto.Plugins, len(plugins))
	for _, plugin := range plugins {
		name := utils.GetFromAnyMap(plugin.GetMeta(), "name", "")
		out[plugin.GetID()] = name
	}

	c.JSON(http.StatusOK, out)
}

// @Summary Get plugins meta
// @Schemes
// @Description Returns a list of plugins meta filtered by parameters
// @Tags plugins
// @Accept json
// @Produce json
// @Param query query plugins.apiGetPluginsMeta.InQuery true "query params"
// @Success 200 {object} dto.PluginsMeta
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/plugins/meta [get]
func (p *Plugins) apiGetPluginsMeta(c *gin.Context) {
	type InQuery struct {
		IDs []string `form:"ids[]" binding:"required"`
	}
	var inQuery InQuery

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err = errors.WithMessage(err, "Plugins: apiGetPluginsMeta: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, p.log)
		return
	}

	out := make(dto.PluginsMeta, len(inQuery.IDs))

	for _, id := range inQuery.IDs {
		pluginID, err := mid.Parse(id)
		if err != nil {
			err = errors.WithMessagef(err, "Plugins: apiGetPluginsMeta: failed to parse uuid: %s", id)
			api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_uuid", err, p.log)
			return
		}

		plugin, ok := p.GetPlugin(pluginID)
		if !ok {
			err = errors.Errorf("Plugins: apiGetPluginsMeta: failed to get plugin by mid: %s", pluginID)
			api.AbortRequest(c, http.StatusNotFound, "plugin_not_found", err, p.log)
			return
		}

		out[pluginID] = dto.PluginMeta(plugin.GetMeta())
	}

	c.JSON(http.StatusOK, out)
}

// @Summary Get plugins options
// @Schemes
// @Description Returns a list of plugins options filtered by parameters
// @Tags plugins
// @Accept json
// @Produce json
// @Param query query plugins.apiGetPluginsOptions.InQuery true "query params"
// @Success 200 {object} dto.PluginsOptions
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/plugins/options [get]
func (p *Plugins) apiGetPluginsOptions(c *gin.Context) {
	type InQuery struct {
		IDs []string `form:"ids[]" binding:"required"`
	}
	var inQuery InQuery

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err = errors.WithMessage(err, "Plugins: apiGetPluginsOptions: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, p.log)
		return
	}

	out := make(dto.PluginsOptions, len(inQuery.IDs))

	for _, id := range inQuery.IDs {
		pluginID, err := mid.Parse(id)
		if err != nil {
			err := errors.WithMessagef(err, "Plugins: apiGetPluginsOptions: failed to parse uuid: %s", id)
			api.AbortRequest(c, http.StatusBadRequest, "invalid_plugin_uuid", err, p.log)
			return
		}

		plugin, ok := p.GetPlugin(pluginID)
		if !ok {
			err = errors.Errorf("Plugins: apiGetPluginsOptions: failed to get plugin by mid: %s", pluginID)
			api.AbortRequest(c, http.StatusNotFound, "plugin_not_found", err, p.log)
			return
		}

		out[pluginID] = plugin.GetOptions()
	}

	c.JSON(http.StatusOK, out)
}
