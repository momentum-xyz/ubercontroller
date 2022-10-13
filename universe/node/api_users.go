package node

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/api"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

func (n *Node) apiUsersCheck(c *gin.Context) {
	inBody := struct {
		IDToken string `json:"idToken"`
	}{}

	if err := c.ShouldBindJSON(&inBody); err != nil {
		n.log.Warn(errors.WithMessage(err, "Node: apiUsersCheck: failed to bind json"))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
		})
		return
	}

	accessToken, idToken, code, err := n.apiCheckTokens(c, api.GetTokenFromRequest(c), inBody.IDToken)
	if err != nil {
		n.log.Warn(errors.WithMessage(err, "Node: apiUsersCheck: failed to check tokens"))
		c.AbortWithStatusJSON(code, gin.H{
			"message": "invalid tokens",
		})
		return
	}

	userEntry, httpCode, err := n.apiGetOrCreateUserFromTokens(c, accessToken, idToken)
	if err != nil {
		n.log.Warn(errors.WithMessage(err, "Node: apiUsersCheck: failed get or create user from tokens"))
		c.AbortWithStatusJSON(httpCode, gin.H{
			"message": "failed to get or create user",
		})
		return
	}

	// TODO: add invitation

	var onBoarder bool
	if userEntry.Profile != nil && userEntry.Profile.OnBoarded != nil && *userEntry.Profile.OnBoarded {
		onBoarder = true
	}

	c.JSON(http.StatusOK, gin.H{
		"userOnboarded": onBoarder,
	})
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
	userID, err := uuid.Parse(idToken.Subject)
	if err != nil {
		return nil, http.StatusBadRequest, errors.WithMessage(err, "failed to parse user id")
	}

	userEntry, err := n.db.UsersGetUserByID(c, userID)
	if err == nil {
		return userEntry, 0, nil
	}

	userEntry = &entry.User{
		UserID:  &userID,
		Profile: &entry.UserProfile{},
	}

	// TODO: check issuer

	nodeSettings, ok := n.GetNodeAttributeValue(
		entry.NewAttributeID(universe.GetSystemPluginID(), "node_settings"),
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
					"wallet",
				),
				*userEntry.UserID,
			),
		}

		walletAddressKey := "address"
		valueModifyFn := func(current *entry.AttributeValue) *entry.AttributeValue {
			if current == nil {
				value := entry.NewAttributeValue()
				(*value)[walletAddressKey] = []string{idToken.Web3Address}
				return value
			}

			address := utils.GetFromAnyMap(*current, walletAddressKey, []any{idToken.Web3Address})
			for i := range address {
				if address[i] == idToken.Web3Address {
					// we don't know where address slice was coming from
					(*current)[walletAddressKey] = address
					return current
				}
			}

			(*current)[walletAddressKey] = append(address, idToken.Web3Address)

			return current
		}

		if err := n.db.UserAttributesUpsertUserAttribute(
			n.ctx, userAttribute, valueModifyFn, modify.Nop[entry.AttributeOptions](),
		); err != nil {
			// TODO: think about rollback
			return nil, http.StatusInternalServerError, errors.WithMessagef(
				err, "failed to upsert user attribute for user: %s", userEntry.UserID,
			)
		}

		n.log.Infof("Node: apiGetOrCreateUserFromTokens: wallet %q added to user: %s", idToken.Web3Address, userEntry.UserID)
	}

	return userEntry, 0, nil
}
