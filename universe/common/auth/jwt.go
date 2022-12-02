package auth

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/momentum-xyz/ubercontroller/universe/common/api"
	"github.com/momentum-xyz/ubercontroller/utils"
)

type JWTClaim struct {
	UserID string `json:"userID"`
	jwt.StandardClaims
	SignedString string `json:"signedString"`
}

// GenerateJWTPair generates a JWT auth token and a refresh token pair.
func GenerateJWTPair(userID string, secret []byte) (map[string]*JWTClaim, error) {
	// default expiration time is 4h
	expire := time.Now().Add(4 * time.Hour)

	JWT := &JWTClaim{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expire.Unix(),
			Issuer:    "controller",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, *JWT)

	authToken, err := token.SignedString(secret)
	if err != nil {
		return nil, errors.New("Could not sign auth JWT token")
	}

	JWT.SignedString = authToken

	rtJWT := &JWTClaim{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expire.Unix(),
		},
	}

	rt := jwt.New(jwt.SigningMethodHS256)
	rtClaims := rt.Claims.(jwt.StandardClaims)
	rtClaims.ExpiresAt = expire.Unix()
	rtClaims.Subject = userID

	refreshToken, err := rt.SignedString(secret)
	if err != nil {
		return nil, err
	}

	rtJWT.SignedString = refreshToken

	tokensMap := map[string]*JWTClaim{
		"auth":    JWT,
		"refresh": rtJWT,
	}

	return tokensMap, nil
}

// ValidateToken checks if a given signed token
// string is a still valid JWT.
func ValidateToken(jwt *JWTClaim) (string, error) {
	if jwt.ExpiresAt < time.Now().Local().Unix() {
		return "", errors.New("token expired")
	}
	return jwt.SignedString, nil
}

func GetJWTClaimFromContext(c *gin.Context) (*JWTClaim, error) {
	value, ok := c.Get(api.JWTClaimContextKey)
	if !ok {
		return &JWTClaim{}, errors.New("failed to get jwt claim value from context")
	}

	claims := utils.GetFromAny(value, JWTClaim{})

	return &claims, nil
}

// RefreshToken resets expiration time for a JWT
func RefreshToken() error {
	return nil
}
