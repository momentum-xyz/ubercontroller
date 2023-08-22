package node

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe/attributes"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/dto"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

// @Summary Gets user contributions by object id
// @Description Returns an object with a nested items array of contributions
// @Tags canvas
// @Security Bearer
// @Param object_id path string true "ObjectID string"
// @Success 200 {object} node.apiCanvasGetUserContributions.Out
// @Failure 400 {object} api.HTTPError
// @Router /api/v4/canvas/{object_id}/user-contributions [get]
func (n *Node) apiCanvasGetUserContributions(c *gin.Context) {
	type InQuery struct {
		OrderBy string `form:"order" json:"order"`
		Limit   uint   `form:"limit,default=10" json:"limit"`
		Offset  uint   `form:"offset" json:"offset"`
		Search  string `form:"q" json:"q"`
	}
	var inQuery InQuery
	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Node: apiCanvasGetUserContributions: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_query", err, n.log)
	}
	var limit uint
	if inQuery.Limit > 100 { // TODO: go 1.21 max function
		limit = 100
	} else {
		limit = inQuery.Limit
	}

	_, attrID, err := attributes.PluginAttributeFromURL(c, n)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiCanvasGetUserContributions: plugin attribute")
		api.AbortRequest(c, http.StatusNotFound, "invalid_param", err, n.log)
		return
	}

	objectID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiCanvasGetUserContributions: failed to parse object umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	objAttrID := entry.NewObjectAttributeID(attrID, objectID)
	ouaDB := n.db.GetObjectUserAttributesDB()
	objectUserAttributes, err := ouaDB.GetObjectUserAttributesByObjectAttributeID(c, objAttrID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiCanvasGetUserContributions: failed to get count")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_count", err, n.log)
		return
	}

	out := dto.UserCanvasContributions{
		Count:  0,
		Limit:  0,
		Offset: 0,
		Items:  nil,
	}

	c.JSON(http.StatusOK, out)
}
