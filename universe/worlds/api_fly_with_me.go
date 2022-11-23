package worlds

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe/common/api"
)

// @Summary Starts a fly with me session for a certain world
// @Schemes
// @Description Initiates a forced fly with me session for all users in a given world
// @Tags worlds
// @Accept json
// @Produce json
// @Param world_id path string true "World ID"
// @Success 200 {object} dto.ExploreOption
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/worlds/{world_id}/fly-with-me/start [post]
func (w *Worlds) apiWorldsFlyWithMeStart(c *gin.Context) {
	worldID, err := uuid.Parse(c.Param("worldID"))
	if err != nil {
		err := errors.WithMessage(err, "Worlds: apiWorldsFlyWithMeStart: failed to parse world id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_world_id", err, w.log)
		return
	}

}

// @Summary Stops a fly with me session for a certain world
// @Schemes
// @Description Stops a forced fly with me session for all users in a given world
// @Tags worlds
// @Accept json
// @Produce json
// @Param world_id path string true "World ID"
// @Success 200 {object} dto.ExploreOption
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/worlds/{world_id}/fly-with-me/stop [post]
func (w *Worlds) apiWorldsFlyWithMeStop(c *gin.Context) {
	worldID, err := uuid.Parse(c.Param("worldID"))
	if err != nil {
		err := errors.WithMessage(err, "Worlds: apiWorldsFlyWithMeStop: failed to parse world id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_world_id", err, w.log)
		return
	}

}
