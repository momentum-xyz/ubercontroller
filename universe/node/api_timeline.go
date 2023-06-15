package node

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/dto"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

// @Summary Get timeline for object
// @Schemes
// @Description Returns a timeline for an object, collection of activities == timeline
// @Tags timeline
// @Accept json
// @Produce json
// @Param query query node.apiTimelineForObject.InQuery true "query params"
// @Success 200 {object} node.apiTimelineForObject.Out
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/objects/{object_id}/timeline [get]
func (n *Node) apiTimelineForObject(c *gin.Context) {
	type InQuery struct {
		StartIndex string `form:"startIndex" binding:"required"`
		PageSize   string `form:"pageSize" binding:"required"`
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

	startIndex, err := strconv.Atoi(inQuery.StartIndex)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiTimelineForObject: failed to convert startIndex to integer")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_start_index_number", err, n.log)
		return
	}

	pageSize, err := strconv.Atoi(inQuery.PageSize)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiTimelineForObject: failed to convert pageSize to integer")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_page_size", err, n.log)
		return
	}

	activities, activitiesTotalCount := n.activities.GetPaginatedActivitiesByObjectID(&objectID, startIndex, pageSize)
	dtoActivities := make([]dto.Activity, 0, len(activities))

	for _, activity := range activities {
		user, err := n.LoadUser(activity.GetUserID())
		if err != nil {
			err := errors.WithMessage(err, "Node: apiTimelineForObject: failed to load user")
			api.AbortRequest(c, http.StatusInternalServerError, "failed_to_load_user", err, n.log)
			return
		}

		var avatarHash, userName string
		profile := user.GetProfile()

		if profile != nil {
			avatarHash = ""
			if profile.AvatarHash != nil {
				avatarHash = *profile.AvatarHash
			}

			userName = ""
			if profile.Name != nil {
				userName = *profile.Name
			}
		}

		object, ok := n.GetObjectFromAllObjects(activity.GetObjectID())
		if !ok {
			err := errors.WithMessage(err, "Node: apiTimelineForObject: failed to get object from all objects")
			api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_object", err, n.log)
			return
		}

		act := dto.Activity{
			ActivityID: activity.GetID(),
			UserID:     activity.GetUserID(),
			ObjectID:   activity.GetObjectID(),
			Type:       activity.GetType(),
			Data:       activity.GetData(),
			AvatarHash: &avatarHash,
			WorldName:  object.GetName(),
			UserName:   &userName,
			CreatedAt:  activity.GetCreatedAt(),
		}

		dtoActivities = append(dtoActivities, act)
	}

	type Out struct {
		Activities []dto.Activity `json:"activities"`
		StartIndex int            `json:"startIndex"`
		PageSize   int            `json:"pageSize"`
		TotalCount int            `json:"totalCount"`
	}
	out := Out{
		Activities: dtoActivities,
		StartIndex: startIndex,
		PageSize:   pageSize,
		TotalCount: activitiesTotalCount,
	}

	c.JSON(http.StatusOK, out)
}

// @Summary Adds a post to a timeline
// @Schemes
// @Description Creates a new post for a timeline
// @Tags timeline
// @Accept json
// @Produce json
// @Param body body node.apiTimelineAddForObject.InBody true "body params"
// @Success 200 {object} node.apiTimelineForObject.Out
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/objects/{object_id}/timeline [post]
func (n *Node) apiTimelineAddForObject(c *gin.Context) {
	type InBody struct {
		Type        string `json:"type" binding:"required"`
		Hash        string `json:"hash" binding:"required"`
		Description string `json:"description"`
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

	user, err := n.LoadUser(userID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiTimelineAddForObject: failed to load user")
		api.AbortRequest(c, http.StatusBadRequest, "failed_to_load_user", err, n.log)
		return
	}

	position := user.GetPosition()
	newActivity, err := n.activities.CreateActivity(umid.New())
	if err != nil {
		err := errors.WithMessage(err, "Node: apiTimelineAddForObject: failed to create activity")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_create_activity", err, n.log)
		return
	}

	if err := newActivity.SetObjectID(objectID, true); err != nil {
		err := errors.WithMessage(err, "Node: apiTimelineAddForObject: failed to set object ID")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_set_object_id", err, n.log)
		return
	}

	if err := newActivity.SetUserID(userID, true); err != nil {
		err := errors.WithMessage(err, "Node: apiTimelineAddForObject: failed to set user ID")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_set_user_id", err, n.log)
		return
	}

	if err := newActivity.SetCreatedAt(time.Now(), true); err != nil {
		err := errors.WithMessage(err, "Node: apiTimelineAddForObject: failed to set created_at")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_set_created_at", err, n.log)
		return
	}

	if !IsValidActivityType(inBody.Type) {
		err := errors.New("Node: apiTimelineAddForObject: invalid activity type")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_type", err, n.log)
		return
	}

	if err := newActivity.SetType((*entry.ActivityType)(&inBody.Type), true); err != nil {
		err := errors.WithMessage(err, "Node: apiTimelineAddForObject: failed to set activity type")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_set_type", err, n.log)
		return
	}

	modifyFn := func(current *entry.ActivityData) (*entry.ActivityData, error) {
		if current == nil {
			current = &entry.ActivityData{}
		}

		current.Position = &position
		current.Hash = &inBody.Hash
		current.Description = &inBody.Description

		return current, nil
	}

	_, err = newActivity.SetData(modifyFn, true)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiTimelineAddForObject: failed to set activity data")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_set_data", err, n.log)
		return
	}

	if err := newActivity.GetCollector().Inject(newActivity); err != nil {
		err := errors.WithMessage(err, "Node: apiTimelineAddForObject: failed to inject activity")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_inject_activity", err, n.log)
		return
	}

	c.JSON(http.StatusCreated, true)
}

// @Summary Edits an activity to a timeline
// @Schemes
// @Description Edits an existing activity for a timeline
// @Tags timeline
// @Accept json
// @Produce json
// @Param body body node.apiTimelineEditForObject.InBody true "body params"
// @Success 200 {object} node.apiTimelineForObject.Out
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/objects/{object_id}/timeline/{activity_id} [patch]
func (n *Node) apiTimelineEditForObject(c *gin.Context) {
	type InBody struct {
		Type        string `json:"type"`
		Hash        string `json:"hash"`
		Description string `json:"description"`
	}
	var inBody InBody

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err := errors.WithMessage(err, "Node: apiTimelineEditForObject: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, n.log)
		return
	}

	activityID, err := umid.Parse(c.Param("activityID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiTimelineEditForObject: failed to parse activity umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_activity_id", err, n.log)
		return
	}

	existingActivity, ok := n.activities.GetActivity(activityID)
	if !ok {
		err := errors.New("Node: apiTimelineEditForObject: failed to find existing activity")
		api.AbortRequest(c, http.StatusNotFound, "activity_not_found", err, n.log)
		return
	}

	if inBody.Type != "" {
		if !IsValidActivityType(inBody.Type) {
			err := errors.New("Node: apiTimelineEditForObject: invalid activity type")
			api.AbortRequest(c, http.StatusBadRequest, "invalid_type", err, n.log)
			return
		}

		err = existingActivity.SetType((*entry.ActivityType)(&inBody.Type), true)
		if err != nil {
			err := errors.WithMessage(err, "Node: apiTimelineEditForObject: failed to set activity type")
			api.AbortRequest(c, http.StatusInternalServerError, "failed_to_set_type", err, n.log)
			return
		}
	}

	modifyFn := func(current *entry.ActivityData) (*entry.ActivityData, error) {
		if current == nil {
			current = &entry.ActivityData{}
		}
		if inBody.Hash != "" {
			current.Hash = &inBody.Hash
		}
		if inBody.Description != "" {
			current.Description = &inBody.Description
		}

		return current, nil
	}

	if err := existingActivity.GetCollector().Modify(existingActivity, modifyFn); err != nil {
		err := errors.WithMessage(err, "Node: apiTimelineEditForObject: failed to modify activity")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_modify_activity", err, n.log)
		return
	}

	c.JSON(http.StatusOK, true)
}

// @Summary Remove an item from a timeline
// @Schemes
// @Description Removes an item from the timeline of an object
// @Tags timeline
// @Accept json
// @Produce json
// @Success 200 {object} node.apiTimelineForObject.Out
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/objects/{object_id}/timeline/{activity_id} [delete]
func (n *Node) apiTimelineRemoveForObject(c *gin.Context) {
	activityID, err := umid.Parse(c.Param("activityID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiTimelineEditForObject: failed to parse activity umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_activity_id", err, n.log)
		return
	}

	activity, ok := n.activities.GetActivity(activityID)
	if !ok {
		err := errors.New("Node: apiTimelineRemoveForObject: activity not found")
		api.AbortRequest(c, http.StatusNotFound, "activity_not_found", err, n.log)
		return
	}

	if err := activity.GetCollector().Remove(activity); err != nil {
		err := errors.WithMessage(err, "Node: apiTimelineRemoveForObject: failed to remove activity")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_remove_activity", err, n.log)
		return
	}

	c.JSON(http.StatusOK, ok)
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

	activities := n.activities.GetActivitiesByUserID(userID)
	type Out struct {
		Activities map[umid.UMID]universe.Activity `json:"activities"`
	}
	out := Out{
		Activities: activities,
	}

	c.JSON(http.StatusOK, out)
}

func IsValidActivityType(t string) bool {
	switch entry.ActivityType(t) {
	case entry.ActivityTypeVideo, entry.ActivityTypeScreenshot:
		return true
	default:
		return false
	}
}
