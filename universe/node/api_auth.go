package node

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe/common/api"
	"github.com/momentum-xyz/ubercontroller/universe/common/api/dto"
	"github.com/momentum-xyz/ubercontroller/universe/common/auth"
	"github.com/pkg/errors"
)

func (n *Node) apiGetChallenge(c *gin.Context) {
	type Out struct {
		Challenge string `json:"challenge"`
	}
	out := Out{
		Challenge: "my super secret challenge",
	}

	c.JSON(http.StatusOK, out)
}

func (n *Node) apiGenToken(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

func (n *Node) apiGenerateGuestToken(c *gin.Context) {
	type InQuery struct {
		UserID string `form:"userID" json:"userID"`
	}
	var inQuery InQuery

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err = errors.WithMessage(err, "Node: apiGenerateGuestToken: failed to bind query parameters")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, n.log)
		return
	}

	predicate := func(k entry.AttributeID, v *entry.AttributePayload) bool {
		if v == nil {
			return false
		}
		key := *v.Value
		if _, ok := key["jwt_key"]; !ok {
			return false
		}
		return true
	}

	jwtAttribute := n.nodeAttributes.Filter(predicate)
	if len(jwtAttribute) != 1 {
		err := errors.New("Node: apiGenerateGuestToken: got more than 1 jwt_key")
		api.AbortRequest(c, http.StatusInternalServerError, "multiple_jwt_keys", err, n.log)
		return
	}

	var jwtAttributeID entry.AttributeID
	for k, _ := range jwtAttribute {
		jwtAttributeID = k
	}

	secretMap, ok := n.nodeAttributes.Load(jwtAttributeID)
	if !ok {
		err := errors.New("Node: apiGenerateGuestToken: no jwt secret")
		api.AbortRequest(c, http.StatusInternalServerError, "no_jwt_secret", err, n.log)
		return
	}

	secret := *secretMap.Value
	secretBytes, err := json.Marshal(secret)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiGenerateGuestToken: failed marshal JWT secret")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_marshal_jwt_secret", err, n.log)
		return
	}

	tokensMap, err := auth.GenerateJWTPair(inQuery.UserID, secretBytes)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiGenerateGuestToken: failed to generate JWT tokens pair")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_generate_jwt", err, n.log)
		return
	}

	uid := uuid.Nil
	if inQuery.UserID != "" {
		uid = uuid.MustParse(inQuery.UserID)
	}
	userEntry, err := n.db.UsersGetUserByID(c, uid)
	if err == nil {
		return
	}

	// userEntry.Auth = tokensMap

	// add tokens to new jsonb column in users
	if err := n.db.UsersUpsertUser(c, userEntry); err != nil {
		err = errors.WithMessage(err, "Node: apiGenerateGuestToken: failed to upsert user")
		api.AbortRequest(c, http.StatusInternalServerError, "invalid_request_query", err, n.log)
		return
	}

	authToken := tokensMap["auth"]
	response := dto.JWTToken{
		UserID:    authToken.UserID,
		Issuer:    authToken.StandardClaims.Issuer,
		Subject:   authToken.StandardClaims.Subject,
		IssuedAt:  strconv.FormatInt(authToken.StandardClaims.IssuedAt, 10),
		ExpiresAt: strconv.FormatInt(authToken.StandardClaims.ExpiresAt, 10),
	}

	c.JSON(http.StatusOK, response)
}

func (n *Node) apiRefreshJWT(c *gin.Context) {
	// TODO
	c.JSON(http.StatusOK, nil)
}
