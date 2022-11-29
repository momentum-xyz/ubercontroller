package worlds

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/posbus-protocol/posbus"
	"github.com/momentum-xyz/ubercontroller/universe/common/api"
	"github.com/momentum-xyz/ubercontroller/universe/common/api/dto"
)

// @Summary Starts a fly to me session
// @Schemes
// @Description Initiates a forced fly to me session
// @Tags worlds
// @Accept json
// @Produce json
// @Param world_id path string true "World ID"
// @Success 200 {object} nil
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/worlds/{world_id}/fly-to-me [post]
func (w *Worlds) apiWorldsFlyToMe(c *gin.Context) {
	worldID, err := uuid.Parse(c.Param("worldID"))
	if err != nil {
		err := errors.WithMessage(err, "Worlds: apiWorldsFlyToMe: failed to parse world id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_world_id", err, w.log)
		return
	}

	world, ok := w.GetWorld(worldID)
	if !ok {
		err := errors.Errorf("Worlds: apiWorldsFlyToMe: world not found: %s", worldID)
		api.AbortRequest(c, http.StatusNotFound, "world_not_found", err, w.log)
		return
	}

	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err = errors.WithMessage(err, "Worlds: apiWorldsFlyToMe: failed to get user id from context")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_user_id", err, w.log)
		return
	}

	user, ok := world.GetUser(userID, true)
	if !ok {
		err := errors.Errorf("Worlds: apiWorldsFlyToMe: user not present in world: %s", worldID)
		api.AbortRequest(c, http.StatusNotFound, "user_not_found", err, w.log)
		return
	}

	userProfile, err := w.db.UsersGetUserProfileByUserID(c, user.GetID())
	if err != nil {
		err = errors.WithMessage(err, "Worlds: apiWorldsFlyToMe: failed to get user profile by user id")
		api.AbortRequest(c, http.StatusNotFound, "profile_not_found", err, w.log)
		return
	}

	userName := ""

	if userProfile != nil {
		if userProfile.Name != nil {
			userName = *userProfile.Name
		}
	}

	fwmDto := dto.FlyToMe{
		Pilot:     user.GetID(),
		PilotName: userName,
		SpaceID:   world.GetID(),
	}

	data, err := json.Marshal(&fwmDto)
	if err != nil {
		err = errors.WithMessage(err, "Worlds: apiWorldsFlyToMe: failed to marshal dto")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_marshal", err, w.log)
		return
	}
	
	msg := posbus.NewRelayToReactMsg(string(dto.FlyToMeTrigger), data).WebsocketMessage()

	if err := world.Send(msg, false); err != nil {
		err = errors.WithMessage(err, "Worlds: apiWorldsFlyToMe: failed to dispatch event")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_dispatch_event", err, w.log)
		return
	}

	c.JSON(http.StatusOK, nil)
}
