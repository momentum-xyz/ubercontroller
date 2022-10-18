package node

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

func (n *Node) apiUsersCheck(c *gin.Context) {
	inBody := struct {
		IDToken string `json:"idToken"`
	}{}

	if err := c.ShouldBindJSON(&inBody); err != nil {
		n.log.Debug(errors.WithMessage(err, "Node: apiUsersCheck: failed to bind json"))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
		})
		return
	}

	accessToken, idToken, code, err := n.apiCheckTokens(c, api.GetTokenFromRequest(c), inBody.IDToken)
	if err != nil {
		n.log.Debug(errors.WithMessage(err, "Node: apiUsersCheck: failed to check tokens"))
		c.AbortWithStatusJSON(code, gin.H{
			"message": "invalid tokens",
		})
		return
	}

	userEntry, httpCode, err := n.apiGetOrCreateUserFromTokens(c, accessToken, idToken)
	if err != nil {
		n.log.Error(errors.WithMessage(err, "Node: apiUsersCheck: failed get or create user from tokens"))
		c.AbortWithStatusJSON(httpCode, gin.H{
			"message": "failed to get or create user",
		})
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
	}
	if idToken.Web3Address != "" {
		outBody.Wallet = &idToken.Web3Address
	}
	if userProfileEntry.Name != nil {
		outBody.Name = *userProfileEntry.Name
	}
	if userEntry.CreatedAt != nil {
		outBody.CreatedAt = userEntry.CreatedAt.String()
	}
	if userEntry.UpdatedAt != nil {
		outBody.UpdatedAt = utils.GetPTR(userEntry.UpdatedAt.String())
	}

	c.JSON(http.StatusOK, outBody)
}

func (n *Node) apiUsersGetMe(c *gin.Context) {
	token, err := api.GetTokenFromContext(c)
	if err != nil {
		n.log.Error(errors.WithMessage(err, "Node: apiUsersGetMe: failed to get token from context"))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "failed to get token",
		})
		return
	}

	userID, err := api.GetUserIDFromToken(token)
	if err != nil {
		n.log.Error(errors.WithMessage(err, "Node: apiUsersGetMe: failed to get user id from token"))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "failed to get user id",
		})
		return
	}

	userEntry, err := n.db.UsersGetUserByID(c, userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "user not found",
		})
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
	}
	if userEntry.UserTypeID != nil {
		outUser.UserTypeID = userEntry.UserTypeID.String()
	}
	if token.Web3Address != "" {
		outUser.Wallet = &token.Web3Address
	}
	if userEntry.CreatedAt != nil {
		outUser.CreatedAt = userEntry.CreatedAt.String()
	}
	if userEntry.UpdatedAt != nil {
		outUser.UpdatedAt = utils.GetPTR(userEntry.UpdatedAt.String())
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
		entry.NewAttributeID(universe.GetSystemPluginID(), universe.NodeSettingsNodeAttributeName),
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
		userAttribute := &entry.UserAttribute{
			UserAttributeID: entry.NewUserAttributeID(
				entry.NewAttributeID(
					universe.GetKusamaPluginID(),
					universe.WalletKusamaUserAttributeName,
				),
				userEntry.UserID,
			),
		}

		walletAddressKey := universe.WalletWalletKusamaUserAttributeKey
		modifyFn := func(current *entry.AttributePayload) *entry.AttributePayload {
			newValue := func() *entry.AttributeValue {
				value := entry.NewAttributeValue()
				(*value)[walletAddressKey] = []string{idToken.Web3Address}
				return value
			}

			if current == nil {
				return entry.NewAttributePayload(newValue(), nil)
			}

			if current.Value == nil {
				current.Value = newValue()
				return current
			}

			address := utils.GetFromAnyMap(*current.Value, walletAddressKey, []any{idToken.Web3Address})
			for i := range address {
				if address[i] == idToken.Web3Address {
					// we don't know where address slice was coming from
					(*current.Value)[walletAddressKey] = address
					return current
				}
			}

			(*current.Value)[walletAddressKey] = append(address, idToken.Web3Address)

			return current
		}

		if err := n.db.UserAttributesUpsertUserAttribute(n.ctx, userAttribute, modifyFn); err != nil {
			// TODO: think about rollback
			return nil, http.StatusInternalServerError, errors.WithMessagef(
				err, "failed to upsert user attribute for user: %s", userEntry.UserID,
			)
		}

		n.log.Infof("Node: apiGetOrCreateUserFromTokens: wallet %q added to user: %s", idToken.Web3Address, userEntry.UserID)
	}

	return userEntry, 0, nil
}
