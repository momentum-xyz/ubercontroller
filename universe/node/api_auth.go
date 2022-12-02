package node

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/common/api"
	"github.com/momentum-xyz/ubercontroller/universe/common/api/dto"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/pkg/errors"
)

// @Summary Generate auth challenge
// @Schemes
// @Description Returns a new generated challenge based on params
// @Tags auth
// @Accept json
// @Produce json
// @Param query query node.apiGenChallenge.InQuery true "query params"
// @Success 200 {object} node.apiGenChallenge.Out
// @Failure 400 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /api/v4/auth/challenge [get]
func (n *Node) apiGenChallenge(c *gin.Context) {
	type InQuery struct {
		Wallet string `form:"wallet" binding:"required"`
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

	challengesKey := universe.Attributes.Kusama.Challenges.Key
	modifyFn := func(current *entry.AttributeValue) (*entry.AttributeValue, error) {
		if current == nil {
			current = entry.NewAttributeValue()
		}

		challenges := utils.GetFromAnyMap(*current, challengesKey, make(map[string]any))
		challenges[inQuery.Wallet] = challenge

		// store challenges because we don't know where we got it from
		(*current)[challengesKey] = challenges

		return current, nil
	}

	if _, err := n.UpdateNodeAttributeValue(
		entry.NewAttributeID(universe.GetKusamaPluginID(), universe.Attributes.Kusama.Challenges.Name),
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

// @Summary Generate auth token
// @Schemes
// @Description Returns a new generated token based on params
// @Tags auth
// @Accept json
// @Produce json
// @Param body body node.apiGenToken.InBody true "body params"
// @Success 200 {object} node.apiGenToken.Out
// @Failure 400 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /api/v4/auth/token [post]
func (n *Node) apiGenToken(c *gin.Context) {
	type InBody struct {
		Wallet          string `json:"wallet" binding:"required"`
		SignedChallenge string `json:"signedChallenge" binding:"required"`
	}
	var inBody InBody

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err := errors.WithMessage(err, "Node: apiGenToken: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	attributeID := entry.NewAttributeID(universe.GetKusamaPluginID(), universe.Attributes.Kusama.Challenges.Name)

	value, ok := n.GetNodeAttributeValue(attributeID)
	if !ok {
		err := errors.Errorf("Node: apiGenToken: node attribute not found")
		api.AbortRequest(c, http.StatusInternalServerError, "attribute_not_found", err, n.log)
		return
	}

	var challenge string
	if value != nil {
		store := utils.GetFromAnyMap(*value, universe.Attributes.Kusama.Challenges.Key, (map[string]any)(nil))
		if store != nil {
			challenge = utils.GetFromAnyMap(store, inBody.Wallet, "")
		}
	}
	if challenge == "" {
		err := errors.Errorf("Node: apiGenToken: challenge not found")
		api.AbortRequest(c, http.StatusNotFound, "challenge_not_found", err, n.log)
		return
	}

	modifyFn := func(current *entry.AttributeValue) (*entry.AttributeValue, error) {
		if current == nil {
			return current, nil
		}

		challenges := utils.GetFromAnyMap(*current, universe.Attributes.Kusama.Challenges.Key, (map[string]any)(nil))
		if challenges == nil {
			return current, nil
		}

		delete(challenges, inBody.Wallet)

		return current, nil
	}

	if _, err := n.UpdateNodeAttributeValue(attributeID, modifyFn, true); err != nil {
		err := errors.WithMessage(err, "Node: apiGenToken: failed to update node attribute value")
		api.AbortRequest(c, http.StatusInternalServerError, "attribute_update_failed", err, n.log)
		return
	}

	valid, err := api.VerifyPolkadotSignature(inBody.Wallet, challenge, inBody.SignedChallenge)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGenToken: failed to verify polkadot signature")
		api.AbortRequest(c, http.StatusBadRequest, "signature_validation_failed", err, n.log)
		return
	}

	if !valid {
		err := errors.Errorf("Node: apiGenToken: invalid signature")
		api.AbortRequest(c, http.StatusForbidden, "invalid_signature", err, n.log)
		return
	}

	type Out struct {
		Token string `json:"token"`
	}

	out := Out{
		Token: "my super secret token",
	}

	c.JSON(http.StatusOK, out)
}

// @Summary Generate jwt guest token
// @Schemes
// @Description Returns a new generated token based on params
// @Tags auth
// @Accept json
// @Produce json
// @Param body body node.apiGuestToken.InBody true "body params"
// @Success 200 {object} dto.User
// @Failure 400 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /api/v4/auth/guest-token [post]
func (n *Node) apiGuestToken(c *gin.Context) {
	type Body struct {
		Name string `json:"name" binding:"required"`
	}
	inBody := Body{}

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiGuestToken: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	userEntry, err := n.apiCreateUserByName(c, &inBody.Name)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiGuestToken: failed get or create user from tokens")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_or_create_user", err, n.log)
		return
	}

	// get jwt secret to sign token
	jwtKeyAttribute, ok := n.GetNodeAttributeValue(entry.NewAttributeID(universe.GetSystemPluginID(), universe.Attributes.Node.Settings.Name))
	if !ok {
		err := errors.New("Node: apiGuestToken: failed to get jwt_key_attribute")
		api.AbortRequest(c, http.StatusInternalServerError, "no_jwt_key", err, n.log)
		return
	}

	secret := utils.GetFromAnyMap(*jwtKeyAttribute, "secret", "")

	token, signedString, err := api.SignJWTToken(userEntry.UserID.String(), []byte(fmt.Sprintf("%v", secret)))
	if err != nil {
		err = errors.WithMessage(err, "Node: apiGuestToken: failed get or create user from tokens")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_or_create_user", err, n.log)
		return
	}

	userEntry.Token.SignedString = &signedString

	claims := token.Claims.(jwt.StandardClaims)

	userEntry.Token.Subject = &claims.Subject
	userEntry.Token.Issuer = &claims.Issuer
	unixExpires := time.Unix(claims.ExpiresAt, 0)
	userEntry.Token.ExpiresAt = &unixExpires
	userEntry.Token.IssuedAt = time.Unix(claims.IssuedAt, 0)

	// add tokens to new jsonb column in users
	if err := n.db.UsersUpsertUser(c, userEntry); err != nil {
		err = errors.WithMessage(err, "Node: apiGuestToken: failed to upsert user")
		api.AbortRequest(c, http.StatusInternalServerError, "invalid_request_query", err, n.log)
		return
	}

	expString := userEntry.Token.ExpiresAt.String()
	issuedString := userEntry.Token.IssuedAt.String()

	outUser := dto.User{
		ID:        userEntry.UserID.String(),
		Name:      *userEntry.Profile.Name,
		CreatedAt: userEntry.CreatedAt.String(),
		JWTToken: dto.JWTToken{
			SignedString: userEntry.Token.SignedString,
			Issuer:       userEntry.Token.Issuer,
			Subject:      userEntry.Token.Subject,
			ExpiresAt:    &expString,
			IssuedAt:     &issuedString,
		},
	}

	if userEntry.UserTypeID != nil {
		outUser.UserTypeID = userEntry.UserTypeID.String()
	}
	if userEntry.UpdatedAt != nil {
		outUser.UpdatedAt = utils.GetPTR(userEntry.UpdatedAt.String())
	}

	c.JSON(http.StatusOK, token)
}
