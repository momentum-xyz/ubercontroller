package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/pkg/errors"
)

type Token struct {
	Issuer   string   `json:"iss"`
	Subject  string   `json:"sub"`
	Audience []string `json:"aud"`
	Expiry   int      `json:"exp"`
	IssuedAt int      `json:"iat"`
	RawToken string   `json:"-"`
}

func VerifyToken(ctx context.Context, token string) (Token, error) {
	parsedToken, err := ParseToken(token)
	if err != nil {
		return parsedToken, errors.WithMessage(err, "failed to parse token")
	}

	return parsedToken, nil
}

func GetTokenFromRequest(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	return strings.TrimPrefix(authHeader, "Bearer ")
}

func GetTokenFromContext(c *gin.Context) (Token, error) {
	value, ok := c.Get(TokenContextKey)
	if !ok {
		return Token{}, errors.Errorf("failed to get token value from context")
	}

	token := utils.GetFromAny(value, Token{})

	return token, nil
}

func GetUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	token, err := GetTokenFromContext(c)
	if err != nil {
		return uuid.Nil, errors.WithMessage(err, "failed to get token from context")
	}
	userID, err := GetUserIDFromToken(token)
	if err != nil {
		return uuid.Nil, errors.WithMessage(err, "failed to get user id from token")
	}
	return userID, nil
}

func GetUserIDFromToken(token Token) (uuid.UUID, error) {
	userID, err := uuid.Parse(token.Subject)
	if err != nil {
		return uuid.Nil, errors.WithMessage(err, "failed to parse user id")
	}
	return userID, nil
}

func ParseToken(token string) (Token, error) {
	var parsedToken Token

	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return parsedToken, errors.Errorf("invalid token, expected 3 parts got %d", len(parts))
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return parsedToken, errors.WithMessage(err, "invalid token payload")
	}
	if err := json.Unmarshal(payload, &parsedToken); err != nil {
		return parsedToken, errors.WithMessage(err, "failed to unmarshal payload")
	}

	parsedToken.RawToken = token

	return parsedToken, nil
}
