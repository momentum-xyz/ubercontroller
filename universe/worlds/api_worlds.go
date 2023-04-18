package worlds

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/momentum-xyz/ubercontroller/utils/umid"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/converters"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/dto"
	"github.com/momentum-xyz/ubercontroller/universe/logic/common"
	"github.com/momentum-xyz/ubercontroller/utils"
)

// @Summary Get world online users
// @Schemes
// @Description Returns a list of online users for specified world
// @Tags worlds
// @Accept json
// @Produce json
// @Param worldID path string true "World UMID"
// @Success 200 {array} dto.User
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/worlds/{object_id}/online-users [get]
func (w *Worlds) apiGetOnlineUsers(c *gin.Context) {
	worldID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Worlds: apiGetOnlineUsers: failed to parse world umid")
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
	userIDs := make([]umid.UMID, 0, len(users))
	for userID, _ := range users {
		userIDs = append(userIDs, userID)
	}

	userEntries, err := w.db.GetUsersDB().GetUsersByIDs(c, userIDs)
	if err != nil {
		err := errors.WithMessage(err, "Worlds: apiGetOnlineUsers: failed to get users")
		api.AbortRequest(c, http.StatusInternalServerError, "get_users_failed", err, w.log)
		return
	}

	guestUserTypeID, err := common.GetGuestUserTypeID()
	if err != nil {
		err := errors.New("Worlds: apiGetOnlineUsers: failed to GetGuestUserType")
		api.AbortRequest(c, http.StatusInternalServerError, "server_error", err, w.log)
		return
	}

	userDTOs := converters.ToUserDTOs(userEntries, guestUserTypeID, false)

	c.JSON(http.StatusOK, userDTOs)
}

// @Summary Get latest worlds
// @Schemes
// @Description Returns a list of six latest created worlds
// @Tags worlds
// @Accept json
// @Produce json
// @Success 200 {array} dto.RecentWorld
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/worlds [get]
func (w *Worlds) apiWorldsGet(c *gin.Context) {
	type InQuery struct {
		Sort  string `form:"sort"`
		Limit string `form:"limit"`
	}
	var inQuery InQuery

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Worlds: apiWorldsGet: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, w.log)
		return
	}

	if inQuery.Limit == "" {
		inQuery.Limit = "100"
	}

	var sortType universe.SortType
	switch inQuery.Sort {
	case "ASC":
		sortType = universe.ASC
	case "DESC":
		sortType = universe.DESC
	default:
		sortType = universe.DESC
	}

	recentWorldIDs, err := w.db.GetWorldsDB().GetWorldIDs(w.ctx, sortType, inQuery.Limit)
	if err != nil {
		err := errors.WithMessage(err, "Worlds: apiWorldsGet: failed to get world ids")
		api.AbortRequest(c, http.StatusInternalServerError, "get_latest_worlds_failed", err, w.log)
		return
	}

	recents := make([]dto.RecentWorld, 0, len(recentWorldIDs))

	for _, worldID := range recentWorldIDs {
		world, _ := w.GetWorld(worldID)

		recent := dto.RecentWorld{
			ID:         world.GetID(),
			Name:       utils.GetPTR(world.GetName()),
			AvatarHash: nil,
		}

		recents = append(recents, recent)
	}

	c.JSON(http.StatusOK, recents)
}

// @Summary Returns objects and one level of children
// @Schemes
// @Description Returns object information and one level of children based on world_id (used in explore widget)
// @Tags worlds
// @Accept json
// @Produce json
// @Param world_id path string true "World UMID"
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

	objectID, err := umid.Parse(inQuery.ObjectID)
	if err != nil {
		err := errors.WithMessage(err, "Worlds: apiWorldsGetObjectsWithChildren: failed to parse object umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, w.log)
		return
	}

	worldID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Worlds: apiWorldsGetObjectsWithChildren: failed to parse world umid")
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
		err := errors.WithMessage(
			err, "Worlds: apiWorldsGetObjectsWithChildren: unable to get options for objects and subobjects",
		)
		api.AbortRequest(c, http.StatusNotFound, "options_not_found", err, w.log)
		return
	}

	c.JSON(http.StatusOK, options)
}

func (w *Worlds) apiWorldsGetRootOptions(root universe.Object) (dto.ExploreOption, error) {
	// objects := root.GetObjects(false)
	var option dto.ExploreOption

	name, description, err := w.apiWorldsResolveNameDescription(root)
	if err != nil {
		return dto.ExploreOption{}, errors.WithMessage(err, "failed to resolve name or description")
	}

	//foundSubObjects, err := w.apiWorldsGetChildrenOptions(objects, 0, 2)
	//if err != nil {
	//	return dto.ExploreOption{}, errors.WithMessage(err, "failed to get children")
	//}

	option = dto.ExploreOption{
		ID:          root.GetID(),
		Name:        utils.GetPTR(name),
		Description: utils.GetPTR(description),
	}

	return option, nil
}

func (w *Worlds) apiWorldsGetChildrenOptions(
	objects map[umid.UMID]universe.Object, currentLevel int, maxLevel int,
) ([]dto.ExploreOption, error) {
	options := make([]dto.ExploreOption, 0, len(objects))
	if currentLevel == maxLevel {
		return options, nil
	}

	for _, object := range objects {
		name, description, err := w.apiWorldsResolveNameDescription(object)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to resolve name or description")
		}

		// subObjects := object.GetObjects(false)
		// foundSubObjects, err := w.apiWorldsGetChildrenOptions(subObjects, currentLevel+1, maxLevel)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to get options")
		}

		option := dto.ExploreOption{
			ID:          object.GetID(),
			Name:        utils.GetPTR(name),
			Description: utils.GetPTR(description),
		}

		options = append(options, option)
	}

	return options, nil
}

func (w *Worlds) apiWorldsResolveNameDescription(object universe.Object) (
	objectName string, objectDescription string, err error,
) {
	var description string
	descriptionAttributeID := entry.NewAttributeID(
		universe.GetSystemPluginID(), universe.ReservedAttributes.Object.Description.Name,
	)
	descriptionValue, _ := object.GetObjectAttributes().GetValue(descriptionAttributeID)
	if descriptionValue != nil {
		description = utils.GetFromAnyMap(*descriptionValue, universe.ReservedAttributes.Object.Description.Name, "")
	}

	return object.GetName(), description, nil
}

// @Summary Search available worlds
// @Schemes
// @Description Returns world information based on a search query and categorizes the results
// @Tags worlds
// @Accept json
// @Produce json
// @Success 200 {object} dto.SearchOptions
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/worlds/explore/search [get]
func (w *Worlds) apiWorldsSearchWorlds(c *gin.Context) {
	type Query struct {
		SearchQuery string `form:"query" binding:"required"`
	}

	inQuery := Query{}

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Worlds: apiWorldsSearchWorlds: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, w.log)
		return
	}

	worlds, err := w.apiWorldsFilterWorlds(inQuery.SearchQuery)
	if err != nil {
		err := errors.WithMessage(err, "Worlds: apiWorldsSearchWorlds: failed to filter objects")
		api.AbortRequest(c, http.StatusBadRequest, "failed_to_filter", err, w.log)
		return
	}

	c.JSON(http.StatusOK, worlds)
}

func (w *Worlds) apiWorldsFilterWorlds(searchQuery string) (dto.SearchOptions, error) {
	node := universe.GetNode()
	worlds := node.GetWorlds()

	predicateFn := func(worldID umid.UMID, world universe.World) bool {
		name, _, err := w.apiWorldsResolveNameDescription(world)
		if err != nil {
			return false
		}

		name = strings.ToLower(name)
		searchQuery = strings.ToLower(searchQuery)
		return strings.Contains(name, searchQuery)
	}

	filteredWorlds := worlds.FilterWorlds(predicateFn)
	options := make([]dto.ExploreOption, 0, len(filteredWorlds))

	for _, filteredWorld := range filteredWorlds {
		name, description, err := w.apiWorldsResolveNameDescription(filteredWorld)
		if err != nil {
			return nil, errors.WithMessage(err, "Worlds: apiWorldsFilterWorlds: failed to get name description")
		}

		option := dto.ExploreOption{
			ID:          filteredWorld.GetID(),
			Name:        utils.GetPTR(name),
			Description: utils.GetPTR(description),
		}

		options = append(options, option)
	}

	return options, nil
}

// @Summary Teleports user from token to another world
// @Schemes
// @Description Teleports user from token to another world
// @Tags worlds
// @Accept json
// @Produce json
// @Param world_id path string true "World UMID"
// @Success 200 {object} nil
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/worlds/{object_id}/teleport-user [post]
func (w *Worlds) apiWorldsTeleportUser(c *gin.Context) {
	objectID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Worlds: apiWorldsTeleportUser: failed to parse world umid")
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
		err = errors.WithMessage(err, "Worlds: apiWorldsTeleportUser: failed to get user umid from token")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_user_id", err, w.log)
		return
	}

	userEntry, err := w.db.GetUsersDB().GetUserByID(c, userID)
	if err != nil {
		err = errors.WithMessage(err, "Worlds: apiWorldsTeleportUser: failed to get user by umid")
		api.AbortRequest(c, http.StatusNotFound, "user_not_found", err, w.log)
		return
	}

	fmt.Sprintln(userEntry)

	c.JSON(http.StatusOK, nil)
}
