package worlds

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/common/api"
	"github.com/momentum-xyz/ubercontroller/universe/common/api/dto"
	"github.com/momentum-xyz/ubercontroller/utils"
)

// @Summary Get world online users
// @Schemes
// @Description Returns a list of online users for specified world
// @Tags worlds
// @Accept json
// @Produce json
// @Param worldID path string true "World ID"
// @Success 200 {array} dto.User
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/worlds/{world_id}/online-users [get]
func (w *Worlds) apiGetOnlineUsers(c *gin.Context) {
	worldID, err := uuid.Parse(c.Param("worldID"))
	if err != nil {
		err := errors.WithMessage(err, "Worlds: apiGetOnlineUsers: failed to parse world id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_world_id", err, w.log)
		return
	}

	world, ok := w.GetWorld(worldID)
	if !ok || world == nil {
		err := errors.WithMessage(err, "Worlds: apiGetOnlineUsers: world not found")
		api.AbortRequest(c, http.StatusNotFound, "world_not_found", err, w.log)
		return
	}

	users := world.GetUsers(true)
	userIDs := make([]uuid.UUID, 0, len(users))
	for userID, _ := range users {
		userIDs = append(userIDs, userID)
	}

	userEntries, err := w.db.UsersGetUsersByIDs(w.ctx, userIDs)
	if err != nil {
		err := errors.WithMessage(err, "Worlds: apiGetOnlineUsers: failed to get users")
		api.AbortRequest(c, http.StatusInternalServerError, "get_users_failed", err, w.log)
		return
	}

	userDTOs := api.ToUserDTOs(userEntries)

	c.JSON(http.StatusOK, userDTOs)
}

// @Summary Returns spaces and one level of children
// @Schemes
// @Description Returns space information and one level of children based on world_id (used in explore widget)
// @Tags worlds
// @Accept json
// @Produce json
// @Param world_id path string true "World ID"
// @Param query query worlds.apiWorldsGetSpacesWithChildren.Query true "query params"
// @Success 200 {object} dto.ExploreOption
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/worlds/{world_id}/explore [get]
func (w *Worlds) apiWorldsGetSpacesWithChildren(c *gin.Context) {
	type Query struct {
		SpaceID string `form:"space_id" binding:"required"`
	}

	inQuery := Query{}

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Worlds: apiWorldsGetSpacesWithChildren: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, w.log)
		return
	}

	spaceID, err := uuid.Parse(inQuery.SpaceID)
	if err != nil {
		err := errors.WithMessage(err, "Worlds: apiWorldsGetSpacesWithChildren: failed to parse space id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_space_id", err, w.log)
		return
	}

	worldID, err := uuid.Parse(c.Param("worldID"))
	if err != nil {
		err := errors.WithMessage(err, "Worlds: apiWorldsGetSpacesWithChildren: failed to parse world id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_world_id", err, w.log)
		return
	}

	world, ok := w.GetWorld(worldID)
	if !ok {
		err := errors.Errorf("Worlds: apiWorldsGetSpacesWithChildren: space not found: %s", worldID)
		api.AbortRequest(c, http.StatusNotFound, "world_not_found", err, w.log)
		return
	}

	root, ok := world.GetSpaceFromAllSpaces(spaceID)
	if !ok {
		err := errors.Errorf("Worlds: apiWorldsGetSpacesWithChildren: failed to get space: %s", spaceID)
		api.AbortRequest(c, http.StatusNotFound, "space_not_found", err, w.log)
		return
	}

	options, err := w.apiWorldsGetRootOptions(root)
	if err != nil {
		err := errors.WithMessage(err, "Worlds: apiWorldsGetSpacesWithChildren: unable to get options for spaces and subspaces")
		api.AbortRequest(c, http.StatusNotFound, "options_not_found", err, w.log)
		return
	}

	c.JSON(http.StatusOK, options)
}

func (w *Worlds) apiWorldsGetRootOptions(root universe.Space) (dto.ExploreOption, error) {
	spaces := root.GetSpaces(false)
	var option dto.ExploreOption

	name, description, err := w.apiWorldsResolveNameDescription(root)
	if err != nil {
		return dto.ExploreOption{}, errors.WithMessage(err, "failed to resolve name or description")
	}

	foundSubSpaces, err := w.apiWorldsGetChildrenOptions(spaces, 0, 2)
	if err != nil {
		return dto.ExploreOption{}, errors.WithMessage(err, "failed to get children")
	}

	option = dto.ExploreOption{
		ID:          root.GetID(),
		Name:        name,
		Description: description,
		SubSpaces:   foundSubSpaces,
	}

	return option, nil
}

func (w *Worlds) apiWorldsGetChildrenOptions(spaces map[uuid.UUID]universe.Space, currentLevel int, maxLevel int) ([]dto.ExploreOption, error) {
	options := make([]dto.ExploreOption, 0, len(spaces))
	if currentLevel == maxLevel {
		return options, nil
	}

	for _, space := range spaces {
		name, description, err := w.apiWorldsResolveNameDescription(space)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to resolve name or description")
		}

		subSpaces := space.GetSpaces(false)
		foundSubSpaces, err := w.apiWorldsGetChildrenOptions(subSpaces, currentLevel+1, maxLevel)
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
		name = utils.GetFromAnyMap(*nameValue, universe.Attributes.Space.Name.Key, "")
	}

	descriptionAttributeID := entry.NewAttributeID(universe.GetSystemPluginID(), universe.Attributes.Space.Description.Name)
	descriptionValue, _ := space.GetSpaceAttributeValue(descriptionAttributeID)

	if descriptionValue != nil {
		description = utils.GetFromAnyMap(*descriptionValue, universe.Attributes.Space.Description.Name, "")
	}

	return name, description, nil
}

// @Summary Search spaces
// @Schemes
// @Description Returns spaces information based on a search query and categorizes the results
// @Tags worlds
// @Accept json
// @Produce json
// @Param world_id path string true "World ID"
// @Param query query worlds.apiWorldsSearchSpaces.Query true "query params"
// @Success 200 {object} dto.SearchOptions
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/worlds/{world_id}/explore/search [get]
func (w *Worlds) apiWorldsSearchSpaces(c *gin.Context) {
	type Query struct {
		SearchQuery string `form:"query" binding:"required"`
	}

	inQuery := Query{}

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Worlds: apiWorldsSearchSpaces: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, w.log)
		return
	}

	worldID, err := uuid.Parse(c.Param("worldID"))
	if err != nil {
		err := errors.WithMessage(err, "Worlds: apiWorldsSearchSpaces: failed to parse world id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_world_id", err, w.log)
		return
	}

	world, ok := w.GetWorld(worldID)
	if !ok {
		err := errors.Errorf("Worlds: apiWorldsSearchSpaces: space not found: %s", worldID)
		api.AbortRequest(c, http.StatusNotFound, "world_not_found", err, w.log)
		return
	}

	spaces, err := w.apiWorldsFilterSpaces(inQuery.SearchQuery, world)
	if err != nil {
		err := errors.WithMessage(err, "Worlds: apiWorldsSearchSpaces: failed to filter spaces")
		api.AbortRequest(c, http.StatusBadRequest, "failed_to_filter", err, w.log)
		return
	}

	c.JSON(http.StatusOK, spaces)
}

func (w *Worlds) apiWorldsFilterSpaces(searchQuery string, world universe.World) (dto.SearchOptions, error) {
	predicateFn := func(spaceID uuid.UUID, space universe.Space) bool {
		name, _, err := w.apiWorldsResolveNameDescription(space)
		if err != nil {
			return false
		}

		name = strings.ToLower(name)
		searchQuery = strings.ToLower(searchQuery)
		return strings.Contains(name, searchQuery)
	}

	spaces := world.FilterAllSpaces(predicateFn)

	options, err := w.apiWorldsGetChildrenOptions(spaces, 0, 1)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get options")
	}

	group := make(dto.SearchOptions)
	for _, option := range options {
		space, ok := world.GetSpaceFromAllSpaces(option.ID)
		if !ok {
			return nil, errors.Errorf("failed to get space: %T", option.ID)
		}

		spaceType := space.GetSpaceType()
		categoryName := spaceType.GetCategoryName()

		group[categoryName] = append(group[categoryName], option)
	}

	return group, nil
}

// @Summary Teleports user from token to another world
// @Schemes
// @Description Teleports user from token to another world
// @Tags worlds
// @Accept json
// @Produce json
// @Param world_id path string true "World ID"
// @Success 200 {object} nil
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/worlds/{world_id}/teleport-user [post]
func (w *Worlds) apiWorldsTeleportUser(c *gin.Context) {
	worldID, err := uuid.Parse(c.Param("worldID"))
	if err != nil {
		err := errors.WithMessage(err, "Worlds: apiWorldsTeleportUser: failed to parse world id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_world_id", err, w.log)
		return
	}

	world, ok := w.GetWorld(worldID)
	if !ok {
		err := errors.Errorf("Worlds: apiWorldsTeleportUser: world not found: %s", worldID)
		api.AbortRequest(c, http.StatusNotFound, "world_not_found", err, w.log)
		return
	}

	fmt.Sprintln(world)

	token, err := api.GetTokenFromContext(c)
	if err != nil {
		err = errors.WithMessage(err, "Worlds: apiWorldsTeleportUser: failed to get token from context")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_token", err, w.log)
		return
	}

	userID, err := api.GetUserIDFromToken(token)
	if err != nil {
		err = errors.WithMessage(err, "Worlds: apiWorldsTeleportUser: failed to get user id from token")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_user_id", err, w.log)
		return
	}

	userEntry, err := w.db.UsersGetUserByID(c, userID)
	if err != nil {
		err = errors.WithMessage(err, "Worlds: apiWorldsTeleportUser: failed to get user by id")
		api.AbortRequest(c, http.StatusNotFound, "user_not_found", err, w.log)
		return
	}

	fmt.Sprintln(userEntry)

	c.JSON(http.StatusOK, nil)
}
