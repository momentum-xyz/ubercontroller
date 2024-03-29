package worlds

import (
	"fmt"
	"math/big"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/converters"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/dto"
	"github.com/momentum-xyz/ubercontroller/universe/logic/common"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

// @Summary Get world online users
// @Description Returns a list of online users for specified world
// @Tags worlds
// @Security Bearer
// @Param object_id path string true "World UMID"
// @Success 200 {array} dto.User
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

// @Summary Get world details
// @Description Returns a world by ID and its details
// @Tags worlds
// @Security Bearer
// @Param object_id path string true "World UMID"
// @Success 200 {array} dto.WorldDetails
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/worlds/{object_id} [get]
func (w *Worlds) apiWorldsGetDetails(c *gin.Context) {
	worldID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Worlds: apiWorldsGetDetails: failed to parse world umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_world_id", err, w.log)
		return
	}

	world, ok := w.GetWorld(worldID)
	if !ok || world == nil {
		err := errors.New("Worlds: apiWorldsGetDetails: world not found")
		api.AbortRequest(c, http.StatusNotFound, "world_not_found", err, w.log)
		return
	}

	node := universe.GetNode()
	ownerID := world.GetOwnerID()
	loadedUser, err := node.LoadUser(ownerID)
	if err != nil {
		err := errors.WithMessage(err, "Worlds: apiWorldsGet: failed to load user")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_load_user", err, w.log)
		return
	}

	currentUserID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "Worlds: apiWorldsGet: user from context")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user", err, w.log)
		return
	}

	isAdmin, err := w.db.GetUserObjectsDB().CheckIsIndirectAdminByID(w.ctx, entry.NewUserObjectID(currentUserID, worldID))
	if err != nil {
		err := errors.WithMessage(err, "Worlds: apiWorldsGet: failed to check is indirect admin")
		api.AbortRequest(c, http.StatusInternalServerError, "check_failed", err, w.log)
		return
	}

	var ownerName *string
	profile := loadedUser.GetProfile()
	if profile != nil {
		if profile.Name != nil {
			ownerName = profile.Name
		}
	}

	stakes, err := w.db.GetStakesDB().GetStakesByWorldID(c, worldID)
	if err != nil {
		err := errors.WithMessage(err, "Worlds: apiWorldsGet: failed to get stakes for world")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_stakes", err, w.log)
		return
	}

	stakeByUser := make(map[umid.UMID]*dto.WorldStaker)
	var totalStake big.Int
	if stakes != nil {
		for _, stake := range stakes {
			hexAddr := utils.AddressToHex(stake.WalletID)
			if len(hexAddr) != 42 && !strings.HasPrefix(hexAddr, "0x") {
				hexAddr = "0x" + hexAddr
			} else if len(hexAddr) != 66 && !strings.HasPrefix(hexAddr, "0x") {
				hexAddr = "0x" + hexAddr
			}
			user, _ := w.db.GetUsersDB().GetUserByWallet(w.ctx, hexAddr)

			if user != nil {
				loadedStaker, err := node.LoadUser(user.UserID)
				if err != nil {
					err := errors.WithMessage(err, "Worlds: apiWorldsGet: failed to load staker")
					api.AbortRequest(c, http.StatusInternalServerError, "failed_to_load_staker", err, w.log)
					return
				}

				stakerProfile := loadedStaker.GetProfile()
				var stakerName *string
				var avatarHash *string
				if stakerProfile != nil {
					if stakerProfile.Name != nil {
						stakerName = stakerProfile.Name
					}
					if stakerProfile.AvatarHash != nil {
						avatarHash = stakerProfile.AvatarHash
					}
				}

				stakeAmt := (*big.Int)(stake.Amount)
				stakeAmtStr := stakeAmt.String()
				totalStake.Add(&totalStake, stakeAmt)

				if worldStaker, ok := stakeByUser[user.UserID]; ok {
					oldStakeAmt := new(big.Int)
					oldStakeAmt.SetString(*worldStaker.Stake, 10)
					newStakeAmt := new(big.Int).Add(oldStakeAmt, stakeAmt)
					newStakeAmtStr := newStakeAmt.String()
					worldStaker.Stake = &newStakeAmtStr
				} else {
					stakeByUser[user.UserID] = &dto.WorldStaker{
						UserID:     user.UserID,
						Name:       stakerName,
						Stake:      &stakeAmtStr,
						AvatarHash: avatarHash,
					}
				}
			}
		}
	}

	worldStakers := make([]dto.WorldStaker, 0, len(stakeByUser))
	for _, staker := range stakeByUser {
		if staker != nil {
			worldStakers = append(worldStakers, *staker)
		}
	}

	sort.SliceStable(worldStakers, func(i, j int) bool {
		iStake := new(big.Int)
		jStake := new(big.Int)
		iStake.SetString(*worldStakers[i].Stake, 10)
		jStake.SetString(*worldStakers[j].Stake, 10)
		return iStake.Cmp(jStake) > 0
	})

	latestStakeComment := ""
	var lastStakeTimestamp time.Time
	for _, stake := range stakes {
		if stake.CreatedAt.After(lastStakeTimestamp) && stake.LastComment != "" {
			latestStakeComment = stake.LastComment
			lastStakeTimestamp = stake.CreatedAt
		}
	}

	totalStakeStr := totalStake.String()
	worldDetails := dto.WorldDetails{
		ID:                 world.GetID(),
		OwnerID:            ownerID,
		OwnerName:          ownerName,
		Name:               utils.GetPTR(world.GetName()),
		Description:        utils.GetPTR(world.GetDescription()),
		StakeTotal:         &totalStakeStr,
		CreatedAt:          world.GetCreatedAt().Format(time.RFC3339),
		UpdatedAt:          world.GetUpdatedAt().Format(time.RFC3339),
		AvatarHash:         utils.GetPTR(world.GetWorldAvatar()),
		WebsiteLink:        utils.GetPTR(world.GetWebsiteLink()),
		WorldStakers:       worldStakers,
		LastStakingComment: &latestStakeComment,
		IsAdmin:            &isAdmin,
	}

	c.JSON(http.StatusOK, worldDetails)
}

// @Summary Get latest worlds
// @Description Returns a list of six latest created worlds
// @Tags worlds
// @Security Bearer
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

	node := universe.GetNode()
	recentWorldIDs, err := w.db.GetWorldsDB().GetWorldIDs(w.ctx, sortType, inQuery.Limit)
	if err != nil {
		err := errors.WithMessage(err, "Worlds: apiWorldsGet: failed to get world ids")
		api.AbortRequest(c, http.StatusInternalServerError, "get_latest_worlds_failed", err, w.log)
		return
	}

	recents := make([]dto.RecentWorld, 0, len(recentWorldIDs))

	for _, worldID := range recentWorldIDs {
		world, ok := w.GetWorld(worldID)
		if !ok {
			err := errors.WithMessage(err, "Worlds: apiWorldsGet: failed to get world by id")
			api.AbortRequest(c, http.StatusInternalServerError, "get_world_by_id_failed", err, w.log)
			return
		}

		ownerID := world.GetOwnerID()
		loadedUser, err := node.LoadUser(ownerID)
		if err != nil {
			err := errors.WithMessage(err, "Worlds: apiWorldsGet: failed to load user")
			api.AbortRequest(c, http.StatusInternalServerError, "failed_to_load_user", err, w.log)
			return
		}

		var ownerName *string
		profile := loadedUser.GetProfile()
		if profile != nil {
			if profile.Name != nil {
				ownerName = profile.Name
			}
		}

		stakes, err := w.db.GetStakesDB().GetStakesByWorldID(c, worldID)
		if err != nil {
			err := errors.WithMessage(err, "Worlds: apiWorldsGet: failed to get stakes for world")
			api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_stakes", err, w.log)
			return
		}

		var totalStake big.Int
		if stakes != nil {
			for _, stake := range stakes {
				s := (*big.Int)(stake.Amount)
				totalStake.Add(&totalStake, s)
			}
		}

		totalStakeStr := totalStake.String()
		recent := dto.RecentWorld{
			ID:          world.GetID(),
			OwnerID:     ownerID,
			OwnerName:   ownerName,
			Name:        utils.GetPTR(world.GetName()),
			Description: utils.GetPTR(world.GetDescription()),
			StakeTotal:  &totalStakeStr,
			AvatarHash:  utils.GetPTR(world.GetWorldAvatar()),
			WebsiteLink: utils.GetPTR(world.GetWebsiteLink()),
		}

		recents = append(recents, recent)
	}

	c.JSON(http.StatusOK, recents)
}

// @Summary Returns objects and one level of children
// @Description Returns object information and one level of children based on world_id (used in explore widget)
// @Tags worlds
// @Security Bearer
// @Param world_id path string true "World UMID"
// @Param query query worlds.apiWorldsGetObjectsWithChildren.Query true "query params"
// @Success 200 {object} dto.ExploreOption
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/worlds/{world_id}/explore [get]
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
// @Description Returns world information based on a search query and categorizes the results
// @Tags worlds
// @Security Bearer
// @Success 200 {object} dto.SearchOptions
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
// @Description Teleports user from token to another world
// @Tags worlds
// @Security Bearer
// @Param object_id path string true "World UMID"
// @Success 200 {object} nil
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

// @Summary Updates world data
// @Description Returns updates world with new data
// @Tags worlds
// @Security Bearer
// @Param object_id path string true "World UMID"
// @Success 200 {array} dto.User
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/worlds/{object_id} [patch]
func (w *Worlds) apiWorldsUpdateByID(c *gin.Context) {
	type InBody struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		WebsiteLink string `json:"website_link"`
	}

	inBody := InBody{}

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err := errors.WithMessage(err, "Worlds: apiWorldsUpdateByID: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, w.log)
		return
	}

	worldID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Worlds: apiWorldsUpdateByID: failed to parse world umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_world_id", err, w.log)
		return
	}

	world, ok := w.GetWorld(worldID)
	if !ok || world == nil {
		err := errors.New("Worlds: apiWorldsUpdateByID: world not found")
		api.AbortRequest(c, http.StatusNotFound, "world_not_found", err, w.log)
		return
	}

	worldEntry := world.GetEntry()
	if err := w.db.GetObjectsDB().UpsertObject(c, worldEntry); err != nil {
		err := errors.New("Worlds: apiWorldsUpdateByID: failed to upsert world")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_upsert", err, w.log)
		return
	}

	c.JSON(http.StatusOK, nil)
}
