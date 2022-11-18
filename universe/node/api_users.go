package node

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/common/api"
	"github.com/momentum-xyz/ubercontroller/universe/common/api/dto"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/merge"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

// @Summary Check user
// @Schemes
// @Description Checks if a logged in user exists in the database and is onboarded, otherwise create new one
// @Tags users
// @Accept json
// @Produce json
// @Param body body node.apiUsersCheck.Body true "body params"
// @Success 200 {object} dto.User
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/users/check [post]
func (n *Node) apiUsersCheck(c *gin.Context) {
	type Body struct {
		IDToken string `json:"idToken" binding:"required"`
	}
	inBody := Body{}

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiUsersCheck: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	accessToken, idToken, code, err := n.apiCheckTokens(c, api.GetTokenFromRequest(c), inBody.IDToken)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiUsersCheck: failed to check tokens")
		api.AbortRequest(c, code, "invalid_tokens", err, n.log)
		return
	}

	userEntry, httpCode, err := n.apiGetOrCreateUserFromTokens(c, accessToken, idToken)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiUsersCheck: failed get or create user from tokens")
		api.AbortRequest(c, httpCode, "failed_to_get_or_create_user", err, n.log)
		return
	}

	// TODO: add invitation

	userProfileEntry := userEntry.Profile

	outProfile := dto.Profile{
		Bio:         userProfileEntry.Bio,
		Location:    userProfileEntry.Location,
		AvatarHash:  userProfileEntry.AvatarHash,
		ProfileLink: userProfileEntry.ProfileLink,
		OnBoarded:   userProfileEntry.OnBoarded,
	}

	outBody := dto.User{
		ID:         userEntry.UserID.String(),
		UserTypeID: userEntry.UserTypeID.String(),
		Profile:    outProfile,
		CreatedAt:  userEntry.CreatedAt.String(),
	}
	if userEntry.UpdatedAt != nil {
		outBody.UpdatedAt = utils.GetPTR(userEntry.UpdatedAt.String())
	}
	if idToken.Web3Address != "" {
		outBody.Wallet = &idToken.Web3Address
	}
	if userProfileEntry.Name != nil {
		outBody.Name = *userProfileEntry.Name
	}

	c.JSON(http.StatusOK, outBody)
}

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
	token, err := api.GetTokenFromContext(c)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiUsersGetMe: failed to get token from context")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_token", err, n.log)
		return
	}

	userID, err := api.GetUserIDFromToken(token)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiUsersGetMe: failed to get user id from token")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_user_id", err, n.log)
		return
	}

	userEntry, err := n.db.UsersGetUserByID(c, userID)
	if err != nil {
		api.AbortRequest(c, http.StatusNotFound, "user_not_found", err, n.log)
		return
	}
	userProfileEntry := userEntry.Profile

	outUser := dto.User{
		ID: userEntry.UserID.String(),
		Profile: dto.Profile{
			Bio:         userProfileEntry.Bio,
			Location:    userProfileEntry.Location,
			AvatarHash:  userProfileEntry.AvatarHash,
			ProfileLink: userProfileEntry.ProfileLink,
			OnBoarded:   userProfileEntry.OnBoarded,
		},
		CreatedAt: userEntry.CreatedAt.String(),
	}
	if userEntry.UserTypeID != nil {
		outUser.UserTypeID = userEntry.UserTypeID.String()
	}
	if userEntry.UpdatedAt != nil {
		outUser.UpdatedAt = utils.GetPTR(userEntry.UpdatedAt.String())
	}
	if token.Web3Address != "" {
		outUser.Wallet = &token.Web3Address
	}
	if userProfileEntry != nil {
		if userProfileEntry.Name != nil {
			outUser.Name = *userProfileEntry.Name
		}
	}

	c.JSON(http.StatusOK, outUser)
}

func (n *Node) apiCheckTokens(c *gin.Context, accessToken, idToken string) (api.Token, api.Token, int, error) {
	parsedAccessToken, err := api.VerifyToken(c, accessToken)
	if err != nil {
		return api.Token{}, api.Token{}, http.StatusForbidden, errors.WithMessage(err,
			"failed to verify access token",
		)
	}

	parsedIDToken, err := api.ParseToken(idToken)
	if err != nil {
		return parsedAccessToken, parsedIDToken, http.StatusBadRequest, errors.WithMessage(err,
			"failed to parse id token",
		)
	}

	if parsedIDToken.Subject != parsedAccessToken.Subject {
		return parsedAccessToken, parsedIDToken, http.StatusBadRequest, errors.WithMessage(
			errors.Errorf("%s != %s", parsedIDToken.Subject, parsedAccessToken.Subject),
			"tokens subject mismatch",
		)
	}

	return parsedAccessToken, parsedIDToken, 0, nil
}

func (n *Node) apiGetOrCreateUserFromTokens(c *gin.Context, accessToken, idToken api.Token) (*entry.User, int, error) {
	userID, err := api.GetUserIDFromToken(accessToken)
	if err != nil {
		return nil, http.StatusBadRequest, errors.WithMessage(err, "failed to get user id from token")
	}

	userEntry, err := n.db.UsersGetUserByID(c, userID)
	if err == nil {
		return userEntry, 0, nil
	}

	userEntry = &entry.User{
		UserID:  userID,
		Profile: &entry.UserProfile{},
	}

	// TODO: check issuer

	nodeSettings, ok := n.GetNodeAttributeValue(
		entry.NewAttributeID(universe.GetSystemPluginID(), universe.Attributes.Node.Settings.Name),
	)
	if !ok {
		return nil, http.StatusInternalServerError, errors.Errorf("failed to get node settings")
	}

	if idToken.Guest.IsGuest {
		guestUserType := utils.GetFromAnyMap(*nodeSettings, "guest_user_type", "")
		guestUserTypeID, err := uuid.Parse(guestUserType)
		if err != nil {
			return nil, http.StatusInternalServerError, errors.Errorf("failed to parse guest user type id")
		}
		userEntry.UserTypeID = &guestUserTypeID

		if err := n.db.UsersUpsertUser(c, userEntry); err != nil {
			return nil, http.StatusInternalServerError, errors.WithMessagef(
				err, "failed to upsert guest: %s", userEntry.UserID,
			)
		}

		n.log.Infof("Node: apiGetOrCreateUserFromTokens: guest created: %s", userEntry.UserID)
	} else {
		// TODO: check idToken web3 type

		if idToken.Web3Address == "" {
			return nil, http.StatusBadRequest, errors.Errorf("empty web3 address: %s", userEntry.UserID)
		}

		// TODO: validate idToken

		normUserType := utils.GetFromAnyMap(*nodeSettings, "normal_user_type", "")
		normUserTypeID, err := uuid.Parse(normUserType)
		if err != nil {
			return nil, http.StatusInternalServerError, errors.Errorf("failed to parse normal user type id")
		}
		userEntry.UserTypeID = &normUserTypeID

		if err := n.db.UsersUpsertUser(c, userEntry); err != nil {
			return nil, http.StatusInternalServerError, errors.WithMessagef(err, "failed to upsert user: %s", userEntry.UserID)
		}

		n.log.Infof("Node: apiGetOrCreateUserFromTokens: user created: %s", userEntry.UserID)

		// adding wallet to user attributes
		userAttributeID := entry.NewUserAttributeID(
			entry.NewAttributeID(
				universe.GetKusamaPluginID(), universe.Attributes.Kusama.User.Wallet.Name,
			),
			userEntry.UserID,
		)

		walletAddressKey := universe.Attributes.Kusama.User.Wallet.Key
		newPayload := entry.NewAttributePayload(
			&entry.AttributeValue{
				walletAddressKey: []string{idToken.Web3Address},
			},
			nil,
		)

		walletAddressKeyPath := ".Value." + walletAddressKey
		if _, err := n.db.UserAttributesUpsertUserAttribute(
			n.ctx, userAttributeID,
			modify.MergeWith(
				newPayload,
				merge.NewTrigger(walletAddressKeyPath, merge.AppendTriggerFn),
				merge.NewTrigger(walletAddressKeyPath, merge.UniqueTriggerFn),
			)); err != nil {
			// TODO: think about rollback
			return nil, http.StatusInternalServerError, errors.WithMessagef(
				err, "failed to upsert user attribute for user: %s", userEntry.UserID,
			)
		}

		n.log.Infof("Node: apiGetOrCreateUserFromTokens: wallet %q added to user: %s", idToken.Web3Address, userEntry.UserID)
	}

	return userEntry, 0, nil
}
