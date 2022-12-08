package node

import (
	"github.com/gin-gonic/gin"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/common/api"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/pkg/errors"
	"net/http"
)

var newsFeedAttributeName = "news_feed"
var newsFeedAttributeKey = "items"
var newsFeedLimit = 100

// @Summary Add newsfeed items
// @Schemes
// @Description Adds provided data to newsfeed
// @Tags newsfeed
// @Accept json
// @Produce json
// @Param body body node.apiNewsFeedAddItem.InBody true "body params"
// @Success 200 {object} nil
// @Failure 400 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /api/v4/newsfeed [post]
func (n *Node) apiNewsFeedAddItem(c *gin.Context) {
	type InBody struct {
		Items []any `json:"items" binding:"required"`
	}

	var inBody InBody
	if err := c.ShouldBindJSON(&inBody); err != nil {
		err := errors.WithMessage(err, "Node: apiNewsFeedAddItem: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	modifyFn := func(current *entry.AttributePayload) (*entry.AttributePayload, error) {
		if current == nil {
			current = entry.NewAttributePayload(nil, nil)
		}
		if current.Value == nil {
			current.Value = entry.NewAttributeValue()
		}

		items := utils.GetFromAnyMap(*current.Value, newsFeedAttributeKey, []any(nil))

		items = append(inBody.Items, items...)

		if len(items) > newsFeedLimit {
			items = items[:newsFeedLimit]
		}

		(*current.Value)[newsFeedAttributeKey] = items

		return current, nil
	}

	if _, err := n.UpsertSpaceAttribute(
		entry.NewAttributeID(universe.GetSystemPluginID(), newsFeedAttributeName), modifyFn, true,
	); err != nil {
		err := errors.WithMessage(err, "Node: apiNewsFeedAddItem: failed to upsert node space attribute")
		api.AbortRequest(c, http.StatusInternalServerError, "upsert_attribute_failed", err, n.log)
		return
	}

	c.JSON(http.StatusCreated, nil)
}

// @Summary Get newsfeed
// @Schemes
// @Description Returns a stored newsfeed data
// @Tags newsfeed
// @Accept json
// @Produce json
// @Success 200 {object} node.apiNewsFeedGetAll.Out
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/newsfeed [get]
func (n *Node) apiNewsFeedGetAll(c *gin.Context) {
	value, ok := n.GetSpaceAttributeValue(entry.NewAttributeID(universe.GetSystemPluginID(), newsFeedAttributeName))
	if !ok || value == nil {
		err := errors.Errorf("Node: apiNewsFeedGetAll: failed to get node space attribute value")
		api.AbortRequest(c, http.StatusNotFound, "attribute_not_found", err, n.log)
		return
	}

	type Out struct {
		Items []any `json:"items"`
	}
	out := Out{
		Items: utils.GetFromAnyMap(*value, newsFeedAttributeKey, []any(nil)),
	}

	c.JSON(http.StatusOK, out)
}
