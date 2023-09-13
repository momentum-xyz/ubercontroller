package node

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

// @Summary Get plugins list
// @Description Returns plugins list
// @Tags plugins
// @Security Bearer
// @Param query query node.apiGetPluginsList.InQuery true "query params"
// @Success 200 {array} entry.Plugin
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/plugins [get]
func (n *Node) apiGetPluginsList(c *gin.Context) {
	type InQuery struct {
		PluginID *string `form:"pluginId" json:"pluginId"`
		Text     *string `form:"text" json:"text"`
	}

	var inQuery InQuery

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err = errors.WithMessage(err, "Node: apiGetPluginsList: failed to bind query parameters")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, n.log)
		return
	}

	var pluginID *umid.UMID
	if inQuery.PluginID != nil {
		id, err := umid.Parse(*inQuery.PluginID)
		if err != nil {
			err = errors.WithMessage(err, "Node: apiGetPluginsList: failed to parse pluginId")
			api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, n.log)
			return
		}
		pluginID = &id
	}

	plugins, err := n.db.GetPluginsDB().GetPlugins(c)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiGetPluginsList: failed to get plugins")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}

	filtered := make([]*entry.Plugin, 0)

	id := pluginID
	text := inQuery.Text

	for _, plugin := range plugins {
		if text == nil {
			if id == nil {
				filtered = append(filtered, plugin)
			}
			if id != nil {
				if *id == plugin.PluginID {
					filtered = append(filtered, plugin)
				}
			}
		}

		if text != nil {
			name := utils.GetFromAnyMap(plugin.Meta, "name", "")
			if id != nil {
				if *id == plugin.PluginID && strings.Contains(name, *text) {
					filtered = append(filtered, plugin)
				}
			}
			if id == nil {
				if strings.Contains(name, *text) {
					filtered = append(filtered, plugin)
				}
			}
		}
	}

	c.JSON(http.StatusOK, filtered)
}
