package node

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/dto"
	"github.com/momentum-xyz/ubercontroller/utils"
)

// @Summary Get the current newsfeed
// @Description Returns a newsfeed, with activities from all timelines
// @Tags newsfeed
// @Security Bearer
// @Param query query node.apiNewsfeedOverview.InQuery true "query params"
// @Success 200 {object} node.apiNewsfeedOverview.Out
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/objects/newsfeed [get]
func (n *Node) apiNewsfeedOverview(c *gin.Context) {
	type InQuery struct {
		StartIndex string `form:"startIndex" json:"startIndex" binding:"required"`
		PageSize   string `form:"pageSize" json:"pageSize" binding:"required"`
	}
	var inQuery InQuery

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Node: apiNewsfeedOverview: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, n.log)
		return
	}

	startIndex, err := strconv.Atoi(inQuery.StartIndex)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiNewsfeedOverview: failed to convert startIndex to integer")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_start_index_number", err, n.log)
		return
	}

	pageSize, err := strconv.Atoi(inQuery.PageSize)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiNewsfeedOverview: failed to convert pageSize to integer")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_page_size", err, n.log)
		return
	}

	activities, activitiesTotalCount := n.activities.GetPaginatedActivities(startIndex, pageSize)
	dtoActivities := make([]dto.Activity, 0, len(activities))

	for _, activity := range activities {
		user, err := n.LoadUser(activity.GetUserID())
		if err != nil {
			err := errors.WithMessagef(err, "Node: apiNewsfeedOverview: %s: failed to load user", activity.GetID())
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
			err := errors.WithMessage(err, "Node: apiNewsfeedOverview: failed to get object from all objects")
			api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_object", err, n.log)
			return
		}

		var worldAvatarHash string
		attributeID := entry.NewAttributeID(universe.GetSystemPluginID(), universe.ReservedAttributes.Object.WorldAvatar.Name)
		objectAttributeValue, ok := object.GetObjectAttributes().GetValue(attributeID)
		if !ok || objectAttributeValue == nil {
			worldAvatarHash = ""
		}
		if objectAttributeValue != nil {
			worldAvatarHash = utils.GetFromAnyMap(*objectAttributeValue, universe.ReservedAttributes.Object.WorldAvatar.Key, "")
		}

		act := dto.Activity{
			ActivityID:      activity.GetID(),
			UserID:          activity.GetUserID(),
			ObjectID:        activity.GetObjectID(),
			Type:            activity.GetType(),
			Data:            activity.GetData(),
			AvatarHash:      &avatarHash,
			WorldName:       object.GetName(),
			WorldAvatarHash: &worldAvatarHash,
			UserName:        &userName,
			CreatedAt:       activity.GetCreatedAt(),
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
