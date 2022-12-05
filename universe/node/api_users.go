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

	return *parsedAccessToken, http.StatusOK, nil
}

func (n *Node) apiCreateGuestUserByName(c *gin.Context, name string) (*entry.User, error) {
	ue := &entry.User{
		UserID: uuid.New(),
		Profile: &entry.UserProfile{
			Name: &name,
		},
	}

	userTypeAttributeValue, ok := n.GetNodeAttributeValue(
		entry.NewAttributeID(universe.GetSystemPluginID(), universe.Attributes.Node.GuestUserType.Name),
	)
	if !ok || userTypeAttributeValue == nil {
		return nil, errors.Errorf("failed to get user type attribute value")
	}

	// user is always guest
	guestUserType := utils.GetFromAnyMap(*userTypeAttributeValue, universe.Attributes.Node.GuestUserType.Key, "")
	guestUserTypeID, err := uuid.Parse(guestUserType)
	if err != nil {
		return nil, errors.Errorf("failed to parse guest user type id")
	}
	ue.UserTypeID = &guestUserTypeID

	err = n.db.UsersUpsertUser(c, ue)

	n.log.Infof("Node: apiCreateGuestUserByName: guest created: %s", ue.UserID)

	return ue, err
}
