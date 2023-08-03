package node

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/dto"
	"github.com/momentum-xyz/ubercontroller/utils"
)

// @Summary Generate auth challenge
// @Description Returns a new generated challenge based on params
// @Tags auth
// @Param query query node.apiGenChallenge.InQuery true "query params"
// @Success 200 {object} node.apiGenChallenge.Out
// @Failure 400 {object} api.HTTPError
// @Router /api/v4/auth/challenge [get]
func (n *Node) apiGenChallenge(c *gin.Context) {
	type InQuery struct {
		Wallet string `form:"wallet" json:"wallet" binding:"required"`
	}
	var inQuery InQuery

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Node: apiGenChallenge: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, n.log)
		return
	}

	challenge, err := api.GenerateChallenge(inQuery.Wallet)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGenChallenge: failed to generate challenge")
		api.AbortRequest(c, http.StatusInternalServerError, "challenge_generation_failed", err, n.log)
		return
	}

	challengesKey := universe.ReservedAttributes.Kusama.Challenges.Key
	modifyFn := func(current *entry.AttributeValue) (*entry.AttributeValue, error) {
		if current == nil {
			current = entry.NewAttributeValue()
		}

		challenges := utils.GetFromAnyMap(*current, challengesKey, map[string]any(nil))
		if challenges == nil {
			challenges = make(map[string]any)
		}
		challenges[inQuery.Wallet] = challenge

		// store challenges because we don't know where we got it from
		(*current)[challengesKey] = challenges

		return current, nil
	}

	if _, err := n.GetNodeAttributes().UpdateValue(
		entry.NewAttributeID(universe.GetKusamaPluginID(), universe.ReservedAttributes.Kusama.Challenges.Name),
		modifyFn, true,
	); err != nil {
		err := errors.WithMessage(err, "Node: apiGenChallenge: failed to update node attribute value")
		api.AbortRequest(c, http.StatusInternalServerError, "attribute_update_failed", err, n.log)
		return
	}

	type Out struct {
		Challenge string `json:"challenge"`
	}
	out := Out{
		Challenge: challenge,
	}

	c.JSON(http.StatusOK, out)
}

// @Summary Verifies a signed challenge
// @Description Returns OK when a signature has been validated
// @Tags auth
// @Security Bearer
// @Param body body node.apiAttachAccount.InBody true "body params"
// @Success 200 {object} nil
// @Failure 400 {object} api.HTTPError
// @Router /api/v4/auth/attach-account [post]
func (n *Node) apiAttachAccount(ctx *gin.Context) {
	type InBody struct {
		Wallet          string `json:"wallet" binding:"required"`
		Network         string `json:"network"`
		SignedChallenge string `json:"signedChallenge" binding:"required"`
	}
	var inBody InBody

	if err := ctx.ShouldBindJSON(&inBody); err != nil {
		err := errors.WithMessage(err, "Node: apiAttachAccount: failed to bind json")
		api.AbortRequest(ctx, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	exists, err := n.db.GetUsersDB().CheckIsUserExistsByWallet(ctx, inBody.Wallet)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiAttachAccount: unable to check if wallet exists")
		api.AbortRequest(ctx, http.StatusInternalServerError, "invalid_wallet_query", err, n.log)
		return
	}

	if exists {
		err := errors.Errorf("Node: apiAttachAccount: user with wallet already exists")
		api.AbortRequest(ctx, http.StatusBadRequest, "wallet_already_exists", err, n.log)
		return
	}

	challengeAttributeID := entry.NewAttributeID(universe.GetKusamaPluginID(), universe.ReservedAttributes.Kusama.Challenges.Name)

	challengesAttributeValue, ok := n.GetNodeAttributes().GetValue(challengeAttributeID)
	if !ok || challengesAttributeValue == nil {
		err := errors.Errorf("Node: apiAttachAccount: node attribute not found")
		api.AbortRequest(ctx, http.StatusInternalServerError, "attribute_not_found", err, n.log)
		return
	}

	var challenge string
	if store := utils.GetFromAnyMap(
		*challengesAttributeValue, universe.ReservedAttributes.Kusama.Challenges.Key, (map[string]any)(nil),
	); store != nil {
		challenge = utils.GetFromAnyMap(store, inBody.Wallet, "")
	}
	if challenge == "" {
		err := errors.Errorf("Node: apiAttachAccount: challenge not found")
		api.AbortRequest(ctx, http.StatusNotFound, "challenge_not_found", err, n.log)
		return
	}

	valid, err := func() (bool, error) {
		switch inBody.Network {
		case "ethereum":
			return api.VerifyEthereumSignature(inBody.Wallet, challenge, inBody.SignedChallenge)
		default:
			return api.VerifyPolkadotSignature(inBody.Wallet, challenge, inBody.SignedChallenge)
		}
	}()
	if err != nil {
		err := errors.WithMessage(err, "Node: apiAttachAccount: failed to update node attribute value")
		api.AbortRequest(ctx, http.StatusInternalServerError, "attribute_update_failed", err, n.log)
		return
	}

	if !valid {
		err := errors.Errorf("Node: apiAttachAccount: ethereum signature appears to be invalid")
		api.AbortRequest(ctx, http.StatusNotFound, "invalid_signature", err, n.log)
		return
	}

	userID, err := api.GetUserIDFromContext(ctx)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiAttachAccount: failed to parse user umid")
		api.AbortRequest(ctx, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	walletKey := universe.ReservedAttributes.Kusama.User.Wallet.Key
	walletAttributeID := entry.NewAttributeID(universe.GetKusamaPluginID(), walletKey)
	userAttributeID := entry.NewUserAttributeID(walletAttributeID, userID)

	modifyFn := func(current *entry.AttributePayload) (*entry.AttributePayload, error) {
		newValue := func() *entry.AttributeValue {
			value := entry.NewAttributeValue()
			walletSlice := make([]string, 0)
			walletSlice = append(walletSlice, inBody.Wallet)
			(*value)[walletKey] = walletSlice
			return value
		}

		if current == nil {
			return entry.NewAttributePayload(newValue(), nil), nil
		}

		if current.Value == nil {
			current.Value = newValue()
			return current, nil
		}

		walletSlice := utils.GetFromAny((*current.Value)[walletKey], []any{})
		if walletSlice == nil {
			return current, nil
		}

		walletSlice = append(walletSlice, inBody.Wallet)
		(*current.Value)[walletKey] = walletSlice

		return current, nil
	}

	payload, err := n.GetUserAttributes().Upsert(userAttributeID, modifyFn, true)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiAttachAccount: failed to upsert user attribute")
		api.AbortRequest(ctx, http.StatusInternalServerError, "failed_to_upsert", err, n.log)
		return
	}
	n.checkNFTWorld(ctx, userID, inBody.Wallet)

	ctx.JSON(http.StatusAccepted, payload.Value)
}

// @Summary Generate auth token
// @Description Returns a new generated token based on params
// @Tags auth
// @Param body body node.apiGenToken.InBody true "body params"
// @Success 200 {object} node.apiGenToken.Out
// @Failure 400 {object} api.HTTPError
// @Failure 403 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/auth/token [post]
func (n *Node) apiGenToken(c *gin.Context) {
	type InBody struct {
		Wallet          string `json:"wallet" binding:"required"`
		Network         string `json:"network"`
		SignedChallenge string `json:"signedChallenge" binding:"required"`
	}
	var inBody InBody

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err := errors.WithMessage(err, "Node: apiGenToken: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	attributeID := entry.NewAttributeID(universe.GetKusamaPluginID(), universe.ReservedAttributes.Kusama.Challenges.Name)

	challengesAttributeValue, ok := n.GetNodeAttributes().GetValue(attributeID)
	if !ok || challengesAttributeValue == nil {
		err := errors.Errorf("Node: apiGenToken: node attribute not found")
		api.AbortRequest(c, http.StatusInternalServerError, "attribute_not_found", err, n.log)
		return
	}

	var challenge string
	if store := utils.GetFromAnyMap(
		*challengesAttributeValue, universe.ReservedAttributes.Kusama.Challenges.Key, (map[string]any)(nil),
	); store != nil {
		challenge = utils.GetFromAnyMap(store, inBody.Wallet, "")
	}
	if challenge == "" {
		err := errors.Errorf("Node: apiGenToken: challenge not found")
		api.AbortRequest(c, http.StatusNotFound, "challenge_not_found", err, n.log)
		return
	}

	valid, err := func() (bool, error) {
		switch inBody.Network {
		case "ethereum":
			return api.VerifyEthereumSignature(inBody.Wallet, challenge, inBody.SignedChallenge)
		default:
			return api.VerifyPolkadotSignature(inBody.Wallet, challenge, inBody.SignedChallenge)
		}
	}()

	if err != nil {
		err := errors.WithMessage(err, "Node: apiGenToken: failed to verify wallet signature")
		api.AbortRequest(c, http.StatusBadRequest, "signature_validation_failed", err, n.log)
		return
	}

	if !valid {
		err := errors.Errorf("Node: apiGenToken: invalid signature")
		api.AbortRequest(c, http.StatusForbidden, "invalid_signature", err, n.log)
		return
	}

	modifyFn := func(current *entry.AttributeValue) (*entry.AttributeValue, error) {
		if current == nil {
			return current, nil
		}

		challenges := utils.GetFromAnyMap(*current, universe.ReservedAttributes.Kusama.Challenges.Key, (map[string]any)(nil))
		if challenges == nil {
			return current, nil
		}

		delete(challenges, inBody.Wallet)

		return current, nil
	}

	if _, err := n.GetNodeAttributes().UpdateValue(attributeID, modifyFn, true); err != nil {
		err := errors.WithMessage(err, "Node: apiGenToken: failed to update node attribute value")
		api.AbortRequest(c, http.StatusInternalServerError, "attribute_update_failed", err, n.log)
		return
	}

	userEntry, httpCode, err := n.apiGetOrCreateUserFromWallet(c, inBody.Wallet)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGenToken: failed to get or create user from meta")
		api.AbortRequest(c, httpCode, "get_or_create_user_failed", err, n.log)
		return
	}

	token, err := api.CreateJWTToken(userEntry.UserID)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiGenToken: failed create token for user")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_create_token", err, n.log)
		return
	}

	type Out struct {
		Token string `json:"token"`
	}
	out := Out{
		Token: token,
	}

	c.JSON(http.StatusOK, out)
}

// @Summary Generate jwt guest token
// @Description Returns a new generated token for guest users
// @Tags auth
// @Success 200 {object} dto.User
// @Failure 400 {object} api.HTTPError
// @Router /api/v4/auth/guest-token [post]
func (n *Node) apiGuestToken(c *gin.Context) {
	visitorName, err := api.GenerateGuestName(c, n.db)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiGuestToken: failed to generate visitor name")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_generate_name", err, n.log)
		return
	}

	userEntry, err := n.apiCreateGuestUserByName(c, visitorName)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiGuestToken: failed create guest user")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_create_guest_user", err, n.log)
		return
	}

	token, err := api.CreateJWTToken(userEntry.UserID)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiGuestToken: failed create token for user")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_create_token", err, n.log)
		return
	}

	name := userEntry.UserID.String()
	if userEntry.Profile.Name != nil {
		name = *userEntry.Profile.Name
	}
	outUser := dto.User{
		ID:         userEntry.UserID,
		UserTypeID: userEntry.UserTypeID,
		Name:       name,
		JWTToken:   &token,
		CreatedAt:  userEntry.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  userEntry.UpdatedAt.Format(time.RFC3339),
		IsGuest:    true,
	}

	c.JSON(http.StatusOK, outUser)
}

// @Summary Remove wallet
// @Description Remove wallet
// @Tags auth
// @Security Bearer
// @Param body body node.apiDeleteWallet.InBody true "body params"
// @Success 200 {object} nil
// @Failure 400 {object} api.HTTPError
// @Router /api/v4/users/me/remove-wallet [delete]
func (n *Node) apiDeleteWallet(c *gin.Context) {
	type InBody struct {
		Wallet string `json:"wallet" binding:"required"`
	}
	var inBody InBody

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err := errors.WithMessage(err, "Node: apiDeleteWallet: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiDeleteWallet: failed to get user umid from context")
		api.AbortRequest(c, http.StatusInternalServerError, "get_user_id_failed", err, n.log)
		return
	}

	walletKey := universe.ReservedAttributes.Kusama.User.Wallet.Key
	walletAttributeID := entry.NewAttributeID(universe.GetKusamaPluginID(), walletKey)
	userAttributeID := entry.NewUserAttributeID(walletAttributeID, userID)

	modifyFn := func(current *entry.AttributePayload) (*entry.AttributePayload, error) {
		if current == nil {
			return current, errors.New("attribute is nil")
		}

		if current.Value == nil {
			return current, errors.New("attribute value is nil")
		}

		walletSlice := utils.GetFromAny((*current.Value)[walletKey], []any{})
		if walletSlice == nil {
			return current, errors.New("wallets slice is nil")
		}

		totalValidWallets := 0
		index := -1
		for i, w := range walletSlice {
			wallet, ok := w.(string)
			if !ok {
				return nil, errors.New("can not cast wallet item to string")
			}
			if strings.HasPrefix(wallet, "0x") {
				totalValidWallets++
			}
			if strings.ToLower(wallet) == strings.ToLower(inBody.Wallet) {
				index = i
			}
		}

		if totalValidWallets < 2 {
			return nil, errors.New("can not remove last wallet")
		}

		if index == -1 {
			return nil, errors.New("such wallet not attached to user")
		}

		newWallets := make([]any, 0)
		newWallets = append(newWallets, walletSlice[:index]...)
		newWallets = append(newWallets, walletSlice[index+1:]...)

		(*current.Value)[walletKey] = newWallets

		return current, nil
	}

	_, err = n.GetUserAttributes().Upsert(userAttributeID, modifyFn, true)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiDeleteWallet: failed to upsert user attribute")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_upsert", err, n.log)
		return
	}

	c.JSON(http.StatusOK, nil)
}
