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
// @Description Returns a timeline for an object, collection of activities == timeline
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

// @Summary Adds an activity to a timeline
// @Schemes
// @Description Creates a new activity for a timeline
// @Tags timeline
// @Accept json
// @Produce json
// @Success 200 {object} node.apiTimelineForObject.Out
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/objects/{object_id}/timeline [post]
func (n *Node) apiTimelineAddForObject(c *gin.Context) {
	type InBody struct {
		Type string `json:"type" binding:"required"`
		Hash string `json:"hash" binding:"required"`
	}
	var inBody InBody

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err := errors.WithMessage(err, "Node: apiTimelineAddForObject: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, n.log)
		return
	}

	objectID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiTimelineAddForObject: failed to parse object umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiTimelineAddForObject: user from context")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user", err, n.log)
		return
	}

	user, ok := n.GetUser(userID, true)
	if !ok {
		err := errors.Errorf("Node: apiTimelineAddForObject: failed to find user in node: %s", userID)
		api.AbortRequest(c, http.StatusNotFound, "user_not_found", err, n.log)
		return
	}

	position := user.GetPosition()
	activity, err := n.activities.CreateActivity(umid.New())
	if err != nil {
		err := errors.WithMessage(err, "Node: apiTimelineAddForObject: failed to create activity")
		api.AbortRequest(c, http.StatusInternalServerError, "invalid_user", err, n.log)
		return
	}

	//Todo: implement setters
	activity.SetObjectID()
	activity.SetUserID()
	activity.SetType()
	activity.SetData()

	c.JSON(http.StatusOK, nil)
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
