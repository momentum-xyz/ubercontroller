package node

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/api"
	"github.com/momentum-xyz/ubercontroller/utils"
)

func (n *Node) apiUsersCheck(c *gin.Context) {
	inBody := struct {
		IDToken string `json:"idToken"`
	}{}

	if err := c.ShouldBindJSON(&inBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": errors.WithMessage(err, "failed to bind json").Error(),
		})
		return
	}

	accessToken, idToken, code, err := n.apiCheckTokens(c, api.GetTokenFromRequest(c), inBody.IDToken)
	if err != nil {
		c.AbortWithStatusJSON(code, gin.H{
			"message": errors.WithMessage(err, "failed to check tokens").Error(),
		})
		return
	}

	userEntry, httpCode, err := n.apiGetOrCreateUserFromTokens(c, accessToken, idToken)
	if err != nil {
		c.AbortWithStatusJSON(httpCode, gin.H{
			"message": errors.WithMessage(err, "failed get or create user from tokens").Error(),
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
		UserID: &userID,
	}

	// TODO: check issuer

	nodeSettings := n.GetNodeAttributes().GetValue(
		types.NewNodeAttributeIndex(universe.GetSystemPluginID(), types.NodeSettingsAttributeName),
	)
	if nodeSettings == nil {
		return nil, http.StatusInternalServerError, errors.Errorf("failed to get node settings: %+v", nodeSettings)
	}

	if idToken.Guest.IsGuest {
		guestUserTypeID := utils.GetFromAnyMap(*nodeSettings, "guest_user_type", uuid.Nil)
		if guestUserTypeID == uuid.Nil {
			return nil, http.StatusInternalServerError, errors.Errorf("failed to get guest user type id")
		}
		userEntry.UserTypeID = &guestUserTypeID

		if err := n.db.UsersUpsertUser(c, userEntry); err != nil {
			return nil, http.StatusInternalServerError, errors.WithMessagef(err, "failed to upsert guest: %s", userEntry.UserID)
		}
	} else {
		// TODO: check idToken web3 type

		if idToken.Web3Address == "" {
			return nil, http.StatusBadRequest, errors.Errorf("empty web3 address: %s", userEntry.UserID)
		}

		// TODO: validate idToken

		normUserTypeID := utils.GetFromAnyMap(*nodeSettings, "normal_user_type", uuid.Nil)
		if normUserTypeID == uuid.Nil {
			return nil, http.StatusInternalServerError, errors.Errorf("failed to get normal user type id")
		}
		userEntry.UserTypeID = &normUserTypeID

		if err := n.db.UsersUpsertUser(c, userEntry); err != nil {
			return nil, http.StatusInternalServerError, errors.WithMessagef(err, "failed to upsert user: %s", userEntry.UserID)
		}

		// TODO: add wallet
	}

	return userEntry, 0, nil
}
