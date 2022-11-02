package worlds

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/momentum-xyz/ubercontroller/universe/api"
	"github.com/pkg/errors"
)

// @Summary Returns spaces and one level of children based on world_id
// @Schemes
// @Description Returns space information and one level of children (used in explore widget)
// @Tags spaces
// @Accept json
// @Produce json
// @Param world_id path string true "World ID"
// @Success 200 {object} dto.SpaceEffectiveOptions
// @Success 500 {object} api.HTTPError
// @Success 400 {object} api.HTTPError
// @Success 404 {object} api.HTTPError
// @Router /api/v4/worlds/{world_id}/explore [get]
func (w *Worlds) apiWorldsGetSpacesWithChildren(c *gin.Context) {
	worldID, err := uuid.Parse(c.Param("worldID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiWorldsGetSpacesWithChildren: failed to parse world id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_world_id", err, w.log)
		return
	}

	world, ok := w.GetWorld(worldID)
	if !ok {
		err := errors.Errorf("Node: apiWorldsGetSpacesWithChildren: space not found: %s", worldID)
		api.AbortRequest(c, http.StatusNotFound, "world_not_found", err, w.log)
		return
	}

	spaces := world.GetSpaces(false)
	options, err := w.GetOptions(spaces)
	if err != nil {
		err := errors.Errorf("Node: apiWorldsGetSpacesWithChildren: unable to get options for spaces and subspaces: %s", err)
		api.AbortRequest(c, http.StatusNotFound, "options_not_found", err, w.log)
		return
	}

	c.JSON(http.StatusOK, options)
}

// @Summary Returns categorized list of spaces based on a search query
// @Schemes
// @Description Returns categorized list of spaces based on a search query for use with the Explore widget
// @Tags spaces
// @Accept json
// @Produce json
// @Param world_id path string true "World ID"
// @Success 200 {object} dto.SpaceEffectiveOptions
// @Success 500 {object} api.HTTPError
// @Success 400 {object} api.HTTPError
// @Success 404 {object} api.HTTPError
// @Router /api/v4/worlds/{world_id}/search [get]
func (w *Worlds) apiWorldsSearch(c *gin.Context) {
	inQuery := struct {
		Query   string `form:"query"`
		SpaceID string `form:"space_id" binding:"required"`
	}{}

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Plugins: apiWorldsSearch: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, w.log)
		return
	}

	worldID, err := uuid.Parse(c.Param("worldID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiWorldsSearch: failed to parse world id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_world_id", err, w.log)
		return
	}

	world, ok := w.GetWorld(worldID)
	if !ok {
		err := errors.Errorf("Node: apiWorldsSearch: world not found: %s", worldID)
		api.AbortRequest(c, http.StatusNotFound, "world_not_found", err, w.log)
		return
	}

	fmt.Sprintln(world)

	// world.GetSpace()
}
