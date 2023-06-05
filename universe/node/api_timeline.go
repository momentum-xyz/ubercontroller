package node

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

// @Param page query int false "Page number"
// @Param pageSize query int false "Number of items per page"
// @Summary Get timeline for object
// @Schemes
// @Description Returns a timeline for an object
// @Tags timeline
// @Accept json
// @Produce json
// @Success 200 {object} node.apiTimelineForObject.Out
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/objects/{object_id}/timeline [get]
func (n *Node) apiTimelineForObject(c *gin.Context) {
	type InQuery struct {
		Page     string `form:"page" binding:"required"`
		PageSize string `form:"pageSize" binding:"required"`
	}
	var inQuery InQuery

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Node: apiTimelineForObject: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, n.log)
		return
	}

	objectID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiTimelineForObject: failed to parse object umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	page, err := strconv.Atoi(inQuery.Page)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiTimelineForObject: failed to convert page to integer")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_page_number", err, n.log)
		return
	}

	pageSize, err := strconv.Atoi(inQuery.PageSize)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiTimelineForObject: failed to convert pageSize to integer")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_page_size", err, n.log)
		return
	}

	activities := n.activities.GetPaginatedActivitiesByObjectID(&objectID, page, pageSize)
	type Out struct {
		Activities []universe.Activity `json:"activities"`
		Page       int                 `json:"page"`
		PageSize   int                 `json:"pageSize"`
		TotalCount int                 `json:"totalCount"`
	}
	out := Out{
		Activities: activities,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: len(activities),
	}

	c.JSON(http.StatusOK, out)
}

// @Summary Get timeline for user
// @Schemes
// @Description Returns a timeline for a user
// @Tags timeline
// @Accept json
// @Produce json
// @Success 200 {object} node.apiTimelineForObject.Out
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/{user_id}/timeline [get]
func (n *Node) apiTimelineForUser(c *gin.Context) {
	userID, err := umid.Parse(c.Param("userID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiTimelineForUser: failed to parse object umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	activities := n.activities.GetActivitiesByUserID(&userID)
	type Out struct {
		Activities map[umid.UMID]universe.Activity `json:"activities"`
	}
	out := Out{
		Activities: activities,
	}

	c.JSON(http.StatusOK, out)
}
