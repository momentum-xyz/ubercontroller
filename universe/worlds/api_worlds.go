package worlds

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/api"
	"github.com/momentum-xyz/ubercontroller/universe/api/dto"
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

	options := make([]dto.ExploreOption, 0, len(spaces))

	for _, space := range spaces {
		nameAttributeID := entry.NewAttributeID(universe.GetSystemPluginID(), universe.SpaceAttributeNameName)
		nameValue, nameOk := space.GetSpaceAttributeValue(nameAttributeID)
		if !nameOk {
			err := errors.Errorf("Node: apiWorldsGetSpacesWithChildren: name attribute value not found: %s", nameAttributeID)
			api.AbortRequest(c, http.StatusNotFound, "attribute_value_name_not_found", err, w.log)
			return
		}

		descriptionAttributeID := entry.NewAttributeID(universe.GetSystemPluginID(), universe.SpaceAttributeNameDescription)
		descriptionValue, descriptionOk := space.GetSpaceAttributeValue(nameAttributeID)
		if !descriptionOk {
			err := errors.Errorf("Node: apiWorldsGetSpacesWithChildren: attribute value not found: %s", descriptionAttributeID)
			api.AbortRequest(c, http.StatusNotFound, "attribute_value_description_not_found", err, w.log)
			return
		}

		name := (*nameValue)[universe.SpaceAttributeNameName]
		description := (*descriptionValue)[universe.SpaceAttributeNameDescription]

		subSpaces := space.GetSpaces(false)
		subSpacesOptions := make([]dto.SubSpace, 0, len(subSpaces))

		for _, subSpace := range subSpaces {
			subSpaceValue, subOk := subSpace.GetSpaceAttributeValue(nameAttributeID)
			if !subOk {
				err := errors.Errorf("Node: apiWorldsGetSpacesWithChildren: attribute value not found: %s", nameAttributeID)
				api.AbortRequest(c, http.StatusNotFound, "attribute_value_not_found", err, w.log)
				return
			}

			if subSpaceValue == nil {
				err := errors.Errorf("Node: apiWorldsGetSpacesWithChildren: subSpaceValue value not found")
				api.AbortRequest(c, http.StatusNotFound, "attribute_value_not_found", err, w.log)
				return
			}

			subSpaceName := (*subSpaceValue)[universe.SpaceAttributeNameName]

			subSpacesOption := dto.SubSpace{
				ID:   subSpace.GetID(),
				Name: subSpaceName,
			}

			subSpacesOptions = append(subSpacesOptions, subSpacesOption)
		}

		option := dto.ExploreOption{
			ID:          space.GetID(),
			Name:        name,
			Description: description,
			SubSpaces:   subSpacesOptions,
		}

		options = append(options, option)
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
func (w *Worlds) apiWorldsExplore(c *gin.Context) {
	inQuery := struct {
		Query   string `form:"query"`
		WorldID string `form:"world_id" binding:"required"`
		SpaceID string `form:"space_id" binding:"required"`
	}{}

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Plugins: apiSpacesExplore: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, w.log)
		return
	}

}
