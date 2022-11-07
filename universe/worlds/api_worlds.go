package worlds

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/api"
	"github.com/momentum-xyz/ubercontroller/universe/api/dto"
	"github.com/momentum-xyz/ubercontroller/utils"
)

// @Summary Returns spaces and one level of children based on world_id
// @Schemes
// @Description Returns space information and one level of children (used in explore widget)
// @Tags spaces
// @Accept json
// @Produce json
// @Param world_id path string true "World ID"
// @Success 200 {object} dto.ExploreOption
// @Success 500 {object} api.HTTPError
// @Success 400 {object} api.HTTPError
// @Success 404 {object} api.HTTPError
// @Router /api/v4/worlds/{world_id}/explore [get]
func (w *Worlds) apiWorldsGetSpacesWithChildren(c *gin.Context) {
	type Query struct {
		SpaceID string `form:"space_id" binding:"required"`
	}

	inQuery := Query{}

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Node: apiWorldsGetSpacesWithChildren: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, w.log)
		return
	}

	spaceID, err := uuid.Parse(inQuery.SpaceID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiWorldsGetSpacesWithChildren: failed to parse space id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_space_id", err, w.log)
		return
	}

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

	root, ok := world.GetSpaceFromAllSpaces(spaceID)
	if !ok {
		err := errors.Errorf("Node: apiWorldsGetSpacesWithChildren: failed to get space: %s", spaceID)
		api.AbortRequest(c, http.StatusNotFound, "space_not_found", err, w.log)
		return
	}

	options, err := w.apiWorldsGetRootOptions(root)
	if err != nil {
		err := errors.Errorf("Node: apiWorldsGetSpacesWithChildren: unable to get options for spaces and subspaces: %s", err)
		api.AbortRequest(c, http.StatusNotFound, "options_not_found", err, w.log)
		return
	}

	c.JSON(http.StatusOK, options)
}

func (w *Worlds) apiWorldsGetRootOptions(root universe.Space) ([]dto.ExploreOption, error) {
	spaces := root.GetSpaces(false)
	options := make([]dto.ExploreOption, 0, len(spaces))

	if len(options) == 0 {
		name, description, err := w.apiWorldsResolveNameDescription(root)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to resolve name or description")
		}

		foundSubSpaces, err := w.apiWorldsGetChildrenOptions(spaces, 0)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to get children")
		}

		option := dto.ExploreOption{
			ID:          root.GetID(),
			Name:        name,
			Description: description,
			SubSpaces:   foundSubSpaces,
		}

		options = append(options, option)
	}

	return options, nil
}

func (w *Worlds) apiWorldsGetChildrenOptions(spaces map[uuid.UUID]universe.Space, level int) ([]dto.ExploreOption, error) {
	options := make([]dto.ExploreOption, 0, len(spaces))
	if level == 2 {
		return options, nil
	}

	for _, space := range spaces {
		name, description, err := w.apiWorldsResolveNameDescription(space)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to resolve name or description")
		}

		subSpaces := space.GetSpaces(false)
		foundSubSpaces, err := w.apiWorldsGetChildrenOptions(subSpaces, level+1)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to get options")
		}

		option := dto.ExploreOption{
			ID:          space.GetID(),
			Name:        name,
			Description: description,
			SubSpaces:   foundSubSpaces,
		}

		options = append(options, option)
	}

	return options, nil
}

func (w *Worlds) apiWorldsResolveNameDescription(space universe.Space) (spaceName string, spaceDescription string, err error) {
	var name string
	var description string

	nameAttributeID := entry.NewAttributeID(universe.GetSystemPluginID(), universe.Attributes.Space.Name.Name)
	nameValue, ok := space.GetSpaceAttributeValue(nameAttributeID)
	if !ok {
		return "", "", errors.Errorf("invalid nameValue: %T", nameAttributeID)
	}

	if nameValue != nil {
		name = utils.GetFromAnyMap(*nameValue, universe.Attributes.Space.Name.Name, "")
	}

	descriptionAttributeID := entry.NewAttributeID(universe.GetSystemPluginID(), universe.Attributes.Space.Description.Name)
	descriptionValue, _ := space.GetSpaceAttributeValue(descriptionAttributeID)

	if descriptionValue != nil {
		description = utils.GetFromAnyMap(*descriptionValue, universe.Attributes.Space.Description.Name, "")
	}

	return name, description, nil
}
