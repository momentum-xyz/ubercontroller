package worlds

import (
	"net/http"

	"github.com/momentum-xyz/ubercontroller/pkg/posbus"
	"github.com/momentum-xyz/ubercontroller/utils/umid"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
)

// @Summary Starts a fly to me session
// @Schemes
// @Description Initiates a forced fly to me session
// @Tags worlds
// @Accept json
// @Produce json
// @Param world_id path string true "World UMID"
// @Success 200 {object} nil
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/worlds/{object_id}/fly-to-me [post]
func (w *Worlds) apiWorldsFlyToMe(c *gin.Context) {
	objectID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Worlds: apiWorldsFlyToMe: failed to parse world umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_world_id", err, w.log)
		return
	}

	world, ok := w.GetWorld(objectID)
	if !ok {
		err := errors.Errorf("Worlds: apiWorldsFlyToMe: world not found: %s", objectID)
		api.AbortRequest(c, http.StatusNotFound, "world_not_found", err, w.log)
		return
	}

	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err = errors.WithMessage(err, "Worlds: apiWorldsFlyToMe: failed to get user umid from context")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_user_id", err, w.log)
		return
	}

	user, ok := world.GetUser(userID, true)
	if !ok {
		err := errors.Errorf("Worlds: apiWorldsFlyToMe: user not present in world: %s", objectID)
		api.AbortRequest(c, http.StatusNotFound, "user_not_found", err, w.log)
		return
	}

	userProfile, err := w.db.GetUsersDB().GetUserProfileByUserID(c, user.GetID())
	if err != nil {
		err = errors.WithMessage(err, "Worlds: apiWorldsFlyToMe: failed to get user profile by user umid")
		api.AbortRequest(c, http.StatusNotFound, "profile_not_found", err, w.log)
		return
	}

	userName := ""

	if userProfile != nil {
		if userProfile.Name != nil {
			userName = *userProfile.Name
		}
	}

	fwm := posbus.FlyToMe{
		Pilot:     user.GetID(),
		PilotName: userName,
		ObjectID:  world.GetID(),
	}

	//data, err := json.Marshal(&fwmDto)
	//if err != nil {
	//	err = errors.WithMessage(err, "Worlds: apiWorldsFlyToMe: failed to marshal dto")
	//	api.AbortRequest(c, http.StatusInternalServerError, "failed_to_marshal", err, w.log)
	//	return
	//}

	msg := posbus.WSMessage(&fwm)

	if err := world.Send(msg, false); err != nil {
		err = errors.WithMessage(err, "Worlds: apiWorldsFlyToMe: failed to dispatch event")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_dispatch_event", err, w.log)
		return
	}

	c.JSON(http.StatusOK, nil)
}
