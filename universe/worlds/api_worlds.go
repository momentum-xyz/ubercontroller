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
	options := w.apiWorldsGetOptions(c, spaces)
	if err != nil {
		err := errors.Errorf("Node: apiWorldsGetSpacesWithChildren: unable to get options for spaces and subspaces: %s", err)
		api.AbortRequest(c, http.StatusNotFound, "options_not_found", err, w.log)
		return
	}

	c.JSON(http.StatusOK, options)
}

func (w *Worlds) apiWorldsGetOptions(c *gin.Context, spaces map[uuid.UUID]universe.Space) []dto.ExploreOption {
	options := make([]dto.ExploreOption, 0, len(spaces))

	for _, space := range spaces {
		var description any

		nameAttributeID := entry.NewAttributeID(universe.GetSystemPluginID(), universe.SpaceAttributeNameName)
		nameValue, ok := space.GetSpaceAttributeValue(nameAttributeID)
		if !ok {
			err := errors.Errorf("Node: apiWorldsGetOptions: could not get name value: %s", nameAttributeID)
			api.AbortRequest(c, http.StatusInternalServerError, "invalid_attribute_id", err, w.log)
			return nil
		}

		name := (*nameValue)[universe.SpaceAttributeNameName]

		if nameValue == nil {
			err := errors.Errorf("Node: apiWorldsGetOptions: could not get nameValue: %s", nameValue)
			api.AbortRequest(c, http.StatusNotFound, "nameValue_not_found", err, w.log)
			return nil
		}

		descriptionAttributeID := entry.NewAttributeID(universe.GetSystemPluginID(), universe.SpaceAttributeDescriptionName)
		descriptionValue, _ := space.GetSpaceAttributeValue(descriptionAttributeID)

		if descriptionValue != nil {
			description = (*descriptionValue)[universe.SpaceAttributeDescriptionName]
		}

		subSpaces := space.GetSpaces(false)
		subOptions := w.apiWorldsGetSubOptions(c, subSpaces)

		option := dto.ExploreOption{
			ID:          space.GetID(),
			Name:        name,
			Description: description,
			SubSpaces:   subOptions,
		}

		options = append(options, option)
	}

	return options
}

func (w *Worlds) apiWorldsGetSubOptions(c *gin.Context, subSpaces map[uuid.UUID]universe.Space) []dto.SubSpace {
	subSpacesOptions := make([]dto.SubSpace, 0, len(subSpaces))

	for _, subSpace := range subSpaces {
		nameAttributeID := entry.NewAttributeID(universe.GetSystemPluginID(), universe.SpaceAttributeNameName)
		subSpaceValue, subOk := subSpace.GetSpaceAttributeValue(nameAttributeID)
		if !subOk {
			err := errors.Errorf("Node: apiWorldsGetSubOptions: could not get name value: %s", nameAttributeID)
			api.AbortRequest(c, http.StatusInternalServerError, "invalid_attribute_id", err, w.log)
			return nil
		}

		subSpaceName := (*subSpaceValue)[universe.SpaceAttributeNameName]

		if subSpaceValue == nil {
			err := errors.Errorf("Node: apiWorldsGetSubOptions: subSpaceValue not found: %s", subSpaceValue)
			api.AbortRequest(c, http.StatusNotFound, "nameValue_not_found", err, w.log)
			return nil
		}

		subSpacesOption := dto.SubSpace{
			ID:   subSpace.GetID(),
			Name: subSpaceName,
		}

		subSpacesOptions = append(subSpacesOptions, subSpacesOption)
	}

	return subSpacesOptions
}
