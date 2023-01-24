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
	"github.com/momentum-xyz/ubercontroller/universe/common/helper"
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
// @Router /api/v4/worlds/{object_id}/online-users [get]
func (w *Worlds) apiGetOnlineUsers(c *gin.Context) {
	worldID, err := uuid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Worlds: apiGetOnlineUsers: failed to parse world id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_world_id", err, w.log)
		return
	}

	world, ok := w.GetWorld(worldID)
	if !ok || world == nil {
		err := errors.New("Worlds: apiGetOnlineUsers: world not found")
		api.AbortRequest(c, http.StatusNotFound, "world_not_found", err, w.log)
		return
	}

	users := world.GetUsers(true)
	userIDs := make([]uuid.UUID, 0, len(users))
	for userID, _ := range users {
		userIDs = append(userIDs, userID)
	}

	userEntries, err := w.db.GetUsersDB().GetUsersByIDs(c, userIDs)
	if err != nil {
		err := errors.WithMessage(err, "Worlds: apiGetOnlineUsers: failed to get users")
		api.AbortRequest(c, http.StatusInternalServerError, "get_users_failed", err, w.log)
		return
	}

	guestUserTypeID, err := helper.GetGuestUserTypeID()
	if err != nil {
		err := errors.New("Worlds: apiGetOnlineUsers: failed to GetGuestUserType")
		api.AbortRequest(c, http.StatusInternalServerError, "server_error", err, w.log)
		return
	}

	userDTOs := api.ToUserDTOs(userEntries, guestUserTypeID, false)

	c.JSON(http.StatusOK, userDTOs)
}

// @Summary Returns objects and one level of children
// @Schemes
// @Description Returns object information and one level of children based on world_id (used in explore widget)
// @Tags worlds
// @Accept json
// @Produce json
// @Param world_id path string true "World ID"
// @Param query query worlds.apiWorldsGetObjectsWithChildren.Query true "query params"
// @Success 200 {object} dto.ExploreOption
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/worlds/{object_id}/explore [get]
func (w *Worlds) apiWorldsGetObjectsWithChildren(c *gin.Context) {
	type Query struct {
		ObjectID string `form:"object_id" binding:"required"`
	}

	inQuery := Query{}

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Worlds: apiWorldsGetObjectsWithChildren: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, w.log)
		return
	}

	objectID, err := uuid.Parse(inQuery.ObjectID)
	if err != nil {
		err := errors.WithMessage(err, "Worlds: apiWorldsGetObjectsWithChildren: failed to parse object id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, w.log)
		return
	}

	worldID, err := uuid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Worlds: apiWorldsGetObjectsWithChildren: failed to parse world id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_world_id", err, w.log)
		return
	}

	world, ok := w.GetWorld(worldID)
	if !ok {
		err := errors.Errorf("Worlds: apiWorldsGetObjectsWithChildren: object not found: %s", objectID)
		api.AbortRequest(c, http.StatusNotFound, "world_not_found", err, w.log)
		return
	}

	root, ok := world.GetObjectFromAllObjects(objectID)
	if !ok {
		err := errors.Errorf("Worlds: apiWorldsGetObjectsWithChildren: failed to get object: %s", objectID)
		api.AbortRequest(c, http.StatusNotFound, "object_not_found", err, w.log)
		return
	}

	options, err := w.apiWorldsGetRootOptions(root)
	if err != nil {
		err := errors.WithMessage(err, "Worlds: apiWorldsGetObjectsWithChildren: unable to get options for objects and subobjects")
		api.AbortRequest(c, http.StatusNotFound, "options_not_found", err, w.log)
		return
	}

	c.JSON(http.StatusOK, options)
}

func (w *Worlds) apiWorldsGetRootOptions(root universe.Object) (dto.ExploreOption, error) {
	objects := root.GetObjects(false)
	var option dto.ExploreOption

	name, description, err := w.apiWorldsResolveNameDescription(root)
	if err != nil {
		return dto.ExploreOption{}, errors.WithMessage(err, "failed to resolve name or description")
	}

	foundSubObjects, err := w.apiWorldsGetChildrenOptions(objects, 0, 2)
	if err != nil {
		return dto.ExploreOption{}, errors.WithMessage(err, "failed to get children")
	}

	option = dto.ExploreOption{
		ID:          root.GetID(),
		Name:        name,
		Description: description,
		SubObjects:  foundSubObjects,
	}

	return option, nil
}

func (w *Worlds) apiWorldsGetChildrenOptions(objects map[uuid.UUID]universe.Object, currentLevel int, maxLevel int) ([]dto.ExploreOption, error) {
	options := make([]dto.ExploreOption, 0, len(objects))
	if currentLevel == maxLevel {
		return options, nil
	}

	for _, object := range objects {
		name, description, err := w.apiWorldsResolveNameDescription(object)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to resolve name or description")
		}

		subObjects := object.GetObjects(false)
		foundSubObjects, err := w.apiWorldsGetChildrenOptions(subObjects, currentLevel+1, maxLevel)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to get options")
		}

		option := dto.ExploreOption{
			ID:          object.GetID(),
			Name:        name,
			Description: description,
			SubObjects:  foundSubObjects,
		}

		options = append(options, option)
	}

	return options, nil
}

func (w *Worlds) apiWorldsResolveNameDescription(object universe.Object) (objectName string, objectDescription string, err error) {
	var description string
	descriptionAttributeID := entry.NewAttributeID(universe.GetSystemPluginID(), universe.ReservedAttributes.Object.Description.Name)
	descriptionValue, _ := object.GetObjectAttributes().GetValue(descriptionAttributeID)
	if descriptionValue != nil {
		description = utils.GetFromAnyMap(*descriptionValue, universe.ReservedAttributes.Object.Description.Name, "")
	}

	return object.GetName(), description, nil
}

// @Summary Search objects
// @Schemes
// @Description Returns objects information based on a search query and categorizes the results
// @Tags worlds
// @Accept json
// @Produce json
// @Param world_id path string true "World ID"
// @Param query query worlds.apiWorldsSearchObjects.Query true "query params"
// @Success 200 {object} dto.SearchOptions
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/worlds/{object_id}/explore/search [get]
func (w *Worlds) apiWorldsSearchObjects(c *gin.Context) {
	type Query struct {
		SearchQuery string `form:"query" binding:"required"`
	}

	inQuery := Query{}

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Worlds: apiWorldsSearchObjects: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, w.log)
		return
	}

	objectID, err := uuid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Worlds: apiWorldsSearchObjects: failed to parse world id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_world_id", err, w.log)
		return
	}

	world, ok := w.GetWorld(objectID)
	if !ok {
		err := errors.Errorf("Worlds: apiWorldsSearchObjects: object not found: %s", objectID)
		api.AbortRequest(c, http.StatusNotFound, "world_not_found", err, w.log)
		return
	}

	objects, err := w.apiWorldsFilterObjects(inQuery.SearchQuery, world)
	if err != nil {
		err := errors.WithMessage(err, "Worlds: apiWorldsSearchObjects: failed to filter objects")
		api.AbortRequest(c, http.StatusBadRequest, "failed_to_filter", err, w.log)
		return
	}

	c.JSON(http.StatusOK, objects)
}

func (w *Worlds) apiWorldsFilterObjects(searchQuery string, world universe.World) (dto.SearchOptions, error) {
	predicateFn := func(objectID uuid.UUID, object universe.Object) bool {
		name, _, err := w.apiWorldsResolveNameDescription(object)
		if err != nil {
			return false
		}

		name = strings.ToLower(name)
		searchQuery = strings.ToLower(searchQuery)
		return strings.Contains(name, searchQuery)
	}

	objects := world.FilterAllObjects(predicateFn)

	options, err := w.apiWorldsGetChildrenOptions(objects, 0, 1)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get options")
	}

	group := make(dto.SearchOptions)
	for _, option := range options {
		object, ok := world.GetObjectFromAllObjects(option.ID)
		if !ok {
			return nil, errors.Errorf("failed to get object: %T", option.ID)
		}

		objectType := object.GetObjectType()
		categoryName := objectType.GetCategoryName()

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
// @Router /api/v4/worlds/{object_id}/teleport-user [post]
func (w *Worlds) apiWorldsTeleportUser(c *gin.Context) {
	objectID, err := uuid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Worlds: apiWorldsTeleportUser: failed to parse world id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_world_id", err, w.log)
		return
	}

	world, ok := w.GetWorld(objectID)
	if !ok {
		err := errors.Errorf("Worlds: apiWorldsTeleportUser: world not found: %s", objectID)
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

	userEntry, err := w.db.GetUsersDB().GetUserByID(c, userID)
	if err != nil {
		err = errors.WithMessage(err, "Worlds: apiWorldsTeleportUser: failed to get user by id")
		api.AbortRequest(c, http.StatusNotFound, "user_not_found", err, w.log)
		return
	}

	fmt.Sprintln(userEntry)

	c.JSON(http.StatusOK, nil)
}
