package node

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

// @Summary Get timeline for object
// @Schemes
// @Description Returns a timeline for an object
// @Tags timeline
// @Accept json
// @Produce json
// @Success 200 {object} node.apiTimelineForObject.Out
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/{object_id}/timeline [get]
func (n *Node) apiTimelineForObject(c *gin.Context) {
	objectID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiTimelineForObject: failed to parse object umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	activities := n.activities.GetActivitiesByObjectID(&objectID)
	type Out struct {
		Activities map[umid.UMID]universe.Activity `json:"activities"`
	}
	out := Out{
		Activities: activities,
	}

	c.JSON(http.StatusOK, out)
}
