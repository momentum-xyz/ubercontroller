package node

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
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

/* TODO: remove

// @Summary Check user
// @Schemes
// @Description Checks if a logged in user exists in the database and is onboarded, otherwise create new one
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} dto.User
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/users/check [post]
func (n *Node) apiUsersCheck(c *gin.Context) {
	userEntry, httpCode, err := n.apiGetOrCreateUserFromTokens(c, api.GetTokenFromRequest(c))
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
		CreatedAt:  userEntry.CreatedAt.Format(time.RFC3339),
	}
	if userEntry.UpdatedAt != nil {
		outBody.UpdatedAt = utils.GetPTR(userEntry.UpdatedAt.Format(time.RFC3339))
	}
	if userProfileEntry.Name != nil {
		outBody.Name = *userProfileEntry.Name
	}

	c.JSON(http.StatusOK, outBody)
}
*/

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
		err = errors.WithMessage(err, "Node: apiUsersGetMe: failed to get user by id")
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
		CreatedAt: userEntry.CreatedAt.Format(time.RFC3339),
	}
	if userEntry.UserTypeID != nil {
		outUser.UserTypeID = userEntry.UserTypeID.String()
	}
	if userEntry.UpdatedAt != nil {
		outUser.UpdatedAt = utils.GetPTR(userEntry.UpdatedAt.Format(time.RFC3339))
	}
	if userProfileEntry != nil {
		if userProfileEntry.Name != nil {
			outUser.Name = *userProfileEntry.Name
		}
	}

	c.JSON(http.StatusOK, outUser)
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
func (n *Node) apiUsersGetById(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userID"))
	if err != nil {
		err = errors.WithMessage(err, "Node: apiUsersGetById: failed to parse user id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	userEntry, err := n.db.UsersGetUserByID(c, userID)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiUsersGetById: failed to get user by id")
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
		CreatedAt: userEntry.CreatedAt.Format(time.RFC3339),
	}
	if userEntry.UserTypeID != nil {
		outUser.UserTypeID = userEntry.UserTypeID.String()
	}
	if userEntry.UpdatedAt != nil {
		outUser.UpdatedAt = utils.GetPTR(userEntry.UpdatedAt.Format(time.RFC3339))
	}
	if userProfileEntry != nil {
		if userProfileEntry.Name != nil {
			outUser.Name = *userProfileEntry.Name
		}
	}

	c.JSON(http.StatusOK, outUser)
}

func (n *Node) apiParseJWT(c *gin.Context, token string) (jwt.Token, int, error) {
	// get jwt secret to sign token
	jwtKeyAttribute, ok := n.GetNodeAttributeValue(
		entry.NewAttributeID(universe.GetSystemPluginID(), universe.Attributes.Node.JWTKey.Name),
	)
	if !ok {
		return jwt.Token{}, http.StatusInternalServerError, errors.New("failed to get jwt_key")
	}
	secret := utils.GetFromAnyMap(*jwtKeyAttribute, universe.Attributes.Node.JWTKey.Key, "")

	parsedAccessToken, err := api.ValidateJWT(token, []byte(secret))
	if err != nil {
		return jwt.Token{}, http.StatusForbidden, errors.WithMessage(err, "failed to verify access token")
	}

	return *parsedAccessToken, 200, nil
}

func (n *Node) apiGetOrCreateUserFromTokens(c *gin.Context, accessToken string) (*entry.User, int, error) {
	jwt, httpCode, err := n.apiParseJWT(c, accessToken)
	if err != nil {
		err := errors.New("Node: apiGetOrCreateUserFromTokens: failed to get jwt_key_attribute")
		return nil, httpCode, err
	}
	userID, err := api.GetUserIDFromToken(jwt)
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
		&entry.AttributeValue{},
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

	return userEntry, 0, nil
}

func (n *Node) apiCreateGuestUserByName(c *gin.Context, name string) (*entry.User, error) {
	ue := &entry.User{
		UserID: uuid.New(),
		Profile: &entry.UserProfile{
			Name: &name,
		},
	}

	nodeSettings, ok := n.GetNodeAttributeValue(
		entry.NewAttributeID(universe.GetSystemPluginID(), universe.Attributes.Node.Settings.Name),
	)
	if !ok || nodeSettings == nil {
		return nil, errors.Errorf("failed to get node settings")
	}

	// user is always guest
	guestUserType := utils.GetFromAnyMap(*nodeSettings, "guest_user_type", "")
	guestUserTypeID, err := uuid.Parse(guestUserType)
	if err != nil {
		return nil, errors.Errorf("failed to parse guest user type id")
	}
	ue.UserTypeID = &guestUserTypeID

	err = n.db.UsersUpsertUser(c, ue)

	n.log.Infof("Node: apiCreateGuestUserByName: guest created: %s", ue.UserID)

	return ue, err
}
