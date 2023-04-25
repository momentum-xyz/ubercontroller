package node

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/converters"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/dto"
	"github.com/momentum-xyz/ubercontroller/universe/logic/common"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/merge"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

// @Summary Get user based on token
// @Schemes
// @Description Returns user information based on token
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} dto.User
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/users/me [get]
func (n *Node) apiUsersGetMe(c *gin.Context) {
	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiUsersGetMe: failed to get user umid from context")
		api.AbortRequest(c, http.StatusInternalServerError, "get_user_id_failed", err, n.log)
		return
	}

	userEntry, err := n.db.GetUsersDB().GetUserByID(c, userID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiUsersGetMe: user not found")
		api.AbortRequest(c, http.StatusNotFound, "user_not_found", err, n.log)
		return
	}

	guestUserTypeID, err := common.GetGuestUserTypeID()
	if err != nil {
		err := errors.New("Node: apiUsersGetMe: failed to GetGuestUserTypeID")
		api.AbortRequest(c, http.StatusInternalServerError, "server_error", err, n.log)
		return
	}

	userDTO := converters.ToUserDTO(userEntry, guestUserTypeID, true)

	c.JSON(http.StatusOK, userDTO)
}

// @Summary Get user profile based on UserID
// @Schemes
// @Description Returns user profile based on UserID
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} dto.User
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/users/{user_id} [get]
func (n *Node) apiUsersGetByID(c *gin.Context) {
	userID, err := umid.Parse(c.Param("userID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiUsersGetByID: failed to parse user umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	userEntry, err := n.db.GetUsersDB().GetUserByID(c, userID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiUsersGetByID: user not found")
		api.AbortRequest(c, http.StatusNotFound, "user_not_found", err, n.log)
		return
	}

	guestUserTypeID, err := common.GetGuestUserTypeID()
	if err != nil {
		err := errors.New("Node: apiUsersGetByID: failed to GetGuestUserTypeID")
		api.AbortRequest(c, http.StatusInternalServerError, "server_error", err, n.log)
		return
	}

	userDTO := converters.ToUserDTO(userEntry, guestUserTypeID, true)

	c.JSON(http.StatusOK, userDTO)
}

// @Summary Get latest users
// @Schemes
// @Description Returns a list of six newest users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} dto.RecentUser
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/users [get]
func (n *Node) apiUsersGet(c *gin.Context) {
	type InQuery struct {
		Sort  string `form:"sort"`
		Limit string `form:"limit"`
	}
	var inQuery InQuery

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Node: apiUsersGet: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, n.log)
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

	normUserTypeID, err := common.GetNormalUserTypeID()
	if err != nil {
		err := errors.WithMessage(err, "Node: apiUsersGet: failed to get normal user type id")
		api.AbortRequest(c, http.StatusInternalServerError, "get_user_type_failed", err, n.log)
		return
	}

	recentUserIDs, err := n.db.GetUsersDB().GetUserIDs(n.ctx, sortType, inQuery.Limit, normUserTypeID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiUsersGet: failed to get latest user ids")
		api.AbortRequest(c, http.StatusInternalServerError, "get_latest_users_failed", err, n.log)
		return
	}

	recents := make([]dto.RecentUser, 0, len(recentUserIDs))

	for _, userID := range recentUserIDs {
		user, err := n.LoadUser(userID)
		if err != nil {
			err := errors.WithMessage(err, "Node: apiUsersGet: failed to get load user by id")
			api.AbortRequest(c, http.StatusInternalServerError, "failed_to_load_user", err, n.log)
			return
		}

		profile := user.GetProfile()

		recent := dto.RecentUser{
			ID:   user.GetID(),
			Name: profile.Name,
			Profile: dto.Profile{
				Bio:         profile.Bio,
				Location:    profile.Location,
				AvatarHash:  profile.AvatarHash,
				ProfileLink: profile.ProfileLink,
			},
		}

		recents = append(recents, recent)
	}

	c.JSON(http.StatusOK, recents)
}

// @Summary Get owned worlds
// @Schemes
// @Description Returns a list of owned Worlds for a user
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} dto.OwnedWorld
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/users/{user_id}/worlds [get]
func (n *Node) apiUsersGetOwnedWorlds(c *gin.Context) {
	userID, err := umid.Parse(c.Param("userID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiUsersGetOwnedWorlds: failed to parse user umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	worlds := n.GetWorldsByOwnerID(userID)
	ownedWorlds := make([]dto.OwnedWorld, 0, len(worlds))
	for _, world := range worlds {
		name := world.GetName()

		ownedWorld := dto.OwnedWorld{
			ID:          world.GetID(),
			OwnerID:     world.GetOwnerID(),
			Name:        &name,
			Description: utils.GetPTR(world.GetDescription()),
			AvatarHash:  nil,
		}

		ownedWorlds = append(ownedWorlds, ownedWorld)
	}

	c.JSON(http.StatusOK, ownedWorlds)
}

// @Summary Get the worlds the user has staked in
// @Schemes
// @Description Returns a list of staked Worlds for a user
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} dto.StakedWorld
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/users/{user_id}/staked-worlds [get]
func (n *Node) apiUsersGetStakedWorlds(c *gin.Context) {
	userID, err := umid.Parse(c.Param("userID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiUsersGetStakedWorlds: failed to parse user umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	wallets, err := n.db.GetUsersDB().GetUserWalletsByUserID(c, userID)
	if err != nil {
		err := errors.WithMessagef(err, "Node: apiUsersGetStakedWorlds: wallets not found for given user_id:%s", userID)
		api.AbortRequest(c, http.StatusNotFound, "wallets_not_found", err, n.log)
		return
	}

	var stakedWorlds []dto.StakedWorld
	for _, wallet := range wallets {
		stakes, err := n.db.GetStakesDB().GetStakesByWalletID(c, *wallet)
		if err != nil {
			err := errors.WithMessage(err, "Node: apiUsersGetStakedWorlds: failed to get stakes for world")
			api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_stakes", err, n.log)
			return
		}

		for _, stake := range stakes {
			world, ok := n.GetObjectFromAllObjects(stake.ObjectID)
			if !ok {
				err := errors.Errorf("Node: apiUsersGetStakedWorlds: world not found: %s", stake.ObjectID)
				api.AbortRequest(c, http.StatusNotFound, "world_not_found", err, n.log)
				return
			}

			stakedWorld := dto.StakedWorld{
				ID:          world.GetID(),
				OwnerID:     world.GetOwnerID(),
				Name:        utils.GetPTR(world.GetName()),
				Description: utils.GetPTR(world.GetDescription()),
				AvatarHash:  nil,
			}
			stakedWorlds = append(stakedWorlds, stakedWorld)
		}
	}

	c.JSON(http.StatusOK, stakedWorlds)
}

//// @Summary Search available users
//// @Schemes
//// @Description Returns user information based on a search query
//// @Tags users
//// @Accept json
//// @Produce json
//// @Success 200 {object} dto.UserSearchResult
//// @Failure 500 {object} api.HTTPError
//// @Failure 400 {object} api.HTTPError
//// @Failure 404 {object} api.HTTPError
//// @Router /api/v4/users/top-stakers [get]
//func (n *Node) apiUsersTopStakers(c *gin.Context) {
//	stakes, err := n.db.GetStakesDB().GetStakesByWorldID(c)
//	if err != nil {
//		err := errors.WithMessage(err, "Node: apiUsersTopStakers: failed to get stakes for user")
//		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_stakes", err, w.log)
//		return
//	}
//
//	guestUserTypeID, err := common.GetGuestUserTypeID()
//	if err != nil {
//		err := errors.New("Node: apiUsersTopStakers: failed to GetGuestUserTypeID")
//		api.AbortRequest(c, http.StatusInternalServerError, "server_error", err, n.log)
//		return
//	}
//
//	userDTO := converters.ToUserDTOs(userEntry, guestUserTypeID, true)
//	c.JSON(http.StatusOK, userDTO)
//}

// @Summary Search available users
// @Schemes
// @Description Returns user information based on a search query
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} dto.UserSearchResult
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/users/search [get]
func (n *Node) apiUsersSearchUsers(c *gin.Context) {
	type Query struct {
		SearchQuery string `form:"query" binding:"required"`
	}

	inQuery := Query{}

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Node: apiUsersSearchUsers: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, n.log)
		return
	}

	users, err := n.apiUsersFilterUsers(inQuery.SearchQuery)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiUsersSearchUsers: failed to filter objects")
		api.AbortRequest(c, http.StatusBadRequest, "failed_to_filter", err, n.log)
		return
	}

	c.JSON(http.StatusOK, users)
}

func (n *Node) apiUsersFilterUsers(searchQuery string) (dto.UserSearchResults, error) {
	predicateFn := func(userID umid.UMID, user universe.User) bool {
		var name string
		loadedUser, _ := n.LoadUser(userID)
		profile := loadedUser.GetProfile()

		if profile != nil && profile.Name != nil {
			name = *profile.Name
		}

		name = strings.ToLower(name)
		searchQuery = strings.ToLower(searchQuery)
		return strings.Contains(name, searchQuery)
	}

	filteredUsers, _ := n.Filter(predicateFn)
	options := make([]dto.UserSearchResult, 0, len(filteredUsers))

	for _, filteredUser := range filteredUsers {
		profile := filteredUser.GetProfile()

		option := dto.UserSearchResult{
			ID:     filteredUser.GetID(),
			Name:   profile.Name,
			Wallet: nil,
			Profile: dto.Profile{
				Bio:         profile.Bio,
				Location:    profile.Location,
				AvatarHash:  profile.AvatarHash,
				ProfileLink: profile.ProfileLink,
			},
		}

		options = append(options, option)
	}

	return options, nil
}

func (n *Node) apiCreateGuestUserByName(ctx context.Context, name string) (*entry.User, error) {
	ue := &entry.User{
		UserID: umid.New(),
		Profile: entry.UserProfile{
			Name: &name,
		},
	}

	guestUserTypeID, err := common.GetGuestUserTypeID()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to GetGuestUserTypeID")
	}

	ue.UserTypeID = guestUserTypeID

	err = n.CreateUsers(ctx, ue)

	n.log.Infof("Node: apiCreateGuestUserByName: guest created: %s", ue.UserID)

	return ue, err
}

func (n *Node) apiGetOrCreateUserFromWallet(ctx context.Context, wallet string) (*entry.User, int, error) {
	userEntry, err := n.db.GetUsersDB().GetUserByWallet(ctx, wallet)
	if err == nil {
		return userEntry, 0, nil
	}

	// walletMeta, err := n.getWalletMetadata(wallet)
	// if err != nil {
	// 	return nil, http.StatusForbidden, errors.WithMessage(err, "failed to get wallet meta")
	// }

	// Temp create empty user
	walletMeta := &WalletMeta{
		Wallet:   wallet,
		UserID:   umid.New(),
		Username: "",
		Avatar:   "",
	}

	userEntry, err = n.createUserFromWalletMeta(ctx, walletMeta)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.WithMessage(err, "failed to create user from wallet meta")
	}

	return userEntry, 0, nil
}

func (n *Node) createUserFromWalletMeta(ctx context.Context, walletMeta *WalletMeta) (*entry.User, error) {
	userEntry := &entry.User{
		UserID: walletMeta.UserID,
		Profile: entry.UserProfile{
			Name:       &walletMeta.Username,
			AvatarHash: &walletMeta.Avatar,
		},
	}

	normUserTypeID, err := common.GetNormalUserTypeID()
	if err != nil {
		return nil, errors.Errorf("failed to get normal user type umid")
	}
	userEntry.UserTypeID = normUserTypeID

	if err := n.CreateUsers(ctx, userEntry); err != nil {
		return nil, errors.WithMessagef(err, "failed to upsert user: %s", userEntry.UserID)
	}

	n.log.Infof("Node: createUserFromWalletMeta: user created: %s", userEntry.UserID)

	// adding wallet to user attributes
	userAttributeID := entry.NewUserAttributeID(
		entry.NewAttributeID(
			universe.GetKusamaPluginID(), universe.ReservedAttributes.Kusama.User.Wallet.Name,
		),
		userEntry.UserID,
	)

	walletAddressKey := universe.ReservedAttributes.Kusama.User.Wallet.Key
	newPayload := entry.NewAttributePayload(
		&entry.AttributeValue{
			walletAddressKey: []any{walletMeta.Wallet},
		},
		nil,
	)

	walletAddressKeyPath := ".Value." + walletAddressKey
	if _, err := n.db.GetUserAttributesDB().UpsertUserAttribute(
		n.ctx, userAttributeID,
		modify.MergeWith(
			newPayload,
			merge.NewTrigger(walletAddressKeyPath, merge.AppendTriggerFn),
			merge.NewTrigger(walletAddressKeyPath, merge.UniqueTriggerFn),
		),
	); err != nil {
		// TODO: think about rollback
		return nil, errors.WithMessagef(
			err, "failed to upsert user attribute for user: %s", userEntry.UserID,
		)
	}

	return userEntry, nil
}
