package node

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/converters"
	"github.com/momentum-xyz/ubercontroller/universe/logic/common"
	"github.com/momentum-xyz/ubercontroller/utils/merge"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
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
		err := errors.WithMessage(err, "Node: apiUsersGetMe: failed to get user id from context")
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
	userID, err := uuid.Parse(c.Param("userID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiUsersGetByID: failed to parse user id")
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

func (n *Node) apiCreateGuestUserByName(ctx context.Context, name string) (*entry.User, error) {
	ue := &entry.User{
		UserID: uuid.New(),
		Profile: entry.UserProfile{
			Name: &name,
		},
	}

	guestUserTypeID, err := common.GetGuestUserTypeID()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to GetGuestUserTypeID")
	}

	ue.UserTypeID = guestUserTypeID

	err = n.db.GetUsersDB().UpsertUser(ctx, ue)

	n.log.Infof("Node: apiCreateGuestUserByName: guest created: %s", ue.UserID)

	return ue, err
}

func (n *Node) apiGetOrCreateUserFromWallet(ctx context.Context, wallet string) (*entry.User, int, error) {
	userEntry, err := n.db.GetUsersDB().GetUserByWallet(ctx, wallet)
	if err == nil {
		return userEntry, 0, nil
	}

	walletMeta, err := n.getWalletMetadata(wallet)
	if err != nil {
		return nil, http.StatusForbidden, errors.WithMessage(err, "failed to get wallet meta")
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
		return nil, errors.Errorf("failed to get normal user type id")
	}
	userEntry.UserTypeID = normUserTypeID

	if err := n.db.GetUsersDB().UpsertUser(ctx, userEntry); err != nil {
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
		)); err != nil {
		// TODO: think about rollback
		return nil, errors.WithMessagef(
			err, "failed to upsert user attribute for user: %s", userEntry.UserID,
		)
	}

	return userEntry, nil
}
