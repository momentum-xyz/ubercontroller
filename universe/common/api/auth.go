package api

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ChainSafe/go-schnorrkel"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/zitadel/oidc/pkg/client"
	"github.com/zitadel/oidc/pkg/client/rs"

	"github.com/momentum-xyz/ubercontroller/utils"
)

func GetTokenFromRequest(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	return strings.TrimPrefix(authHeader, "Bearer ")
}

func GetTokenFromContext(c *gin.Context) (*jwt.Token, error) {
	value, ok := c.Get(TokenContextKey)
	if !ok {
		return nil, errors.Errorf("failed to get token value from context")
	}

	token := utils.GetFromAny(value, jwt.Token{})

	return &token, nil
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

func GetUserIDFromToken(token *jwt.Token) (uuid.UUID, error) {
	if token == nil {
		return uuid.Nil, errors.New("got nil token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok /*&& token.Valid*/ {
		userID, err := uuid.Parse(fmt.Sprint(claims["sub"])) // TODO! proper jwt parsing
		if err != nil {
			return uuid.Nil, errors.WithMessage(err, "failed to parse user id")
		}
		return userID, nil
	} else {
		return uuid.Nil, errors.New("Failed to get token claims")
	}

}

func GenerateChallenge(wallet string) (string, error) {
	return fmt.Sprintf(
		"Please sign this message with the private key for address %s to prove that you own it. %s",
		wallet, uuid.New().String(),
	), nil
}

func VerifyPolkadotSignature(wallet, challenge, signature string) (bool, error) {
	pub, err := schnorrkel.NewPublicKeyFromHex(wallet)
	if err != nil {
		return false, errors.WithMessage(err, "failed to get public key")
	}
	sig, err := schnorrkel.NewSignatureFromHex(signature)
	if err != nil {
		return false, errors.WithMessage(err, "failed to get signature")
	}

	trContext := []byte("substrate")
	trMessage := bytes.NewBufferString("<Bytes>")
	trMessage.WriteString(challenge)
	trMessage.WriteString("</Bytes>")

	transcript := schnorrkel.NewSigningContext(trContext, trMessage.Bytes())

	return pub.Verify(sig, transcript)
}

func createProvider(provider string) (rs.ResourceServer, error) {
	cfg := api.cfg.Auth
	oidcURL := cfg.OIDCURL
	clientID := cfg.GetIDByProvider(provider)
	secret := cfg.GetSecretByProvider(provider)
	introspectURL := cfg.GetIntrospectURLByProvider(provider)

	api.log.Infof("API: creating oidc provider: %s, url: %s, client id: %s, secret: %s, introspect url: %s",
		provider, oidcURL, clientID, secret, introspectURL)

	opts := make([]rs.Option, 0, 1)
	if introspectURL != "" {
		oidcConfig, err := client.Discover(oidcURL, http.DefaultClient)
		if err != nil {
			return nil, errors.WithMessagef(err, "failed to discover provider: %s", provider)
		}
		opts = append(opts, rs.WithStaticEndpoints(oidcConfig.TokenEndpoint, introspectURL))
	}

	oidcProvider, err := rs.NewResourceServerClientCredentials(oidcURL, clientID, secret, opts...)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to create resource server client: %s", provider)
	}

	return oidcProvider, nil
}

// SignJWTToken saves a jwt token with the given userID as subject
// and signed with the given secret
func SignJWTToken(userID string, secret []byte) (string, error) {
	claims := jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(4 * time.Hour).Unix(),
		Issuer:    "ubercontroller",
		Subject:   userID,
	}

	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedString, err := jwt.SignedString(secret)
	if err != nil {
		return "", err
	}

	return signedString, nil
}

func ValidateJWT(signedString string, secret []byte) (*jwt.Token, error) {
	parser := new(jwt.Parser)
	token, _, err := parser.ParseUnverified(signedString, jwt.MapClaims{})
	return token, err
	/* TODO:
	return jwt.Parse(signedString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid token %v", token.Header["alg"])
		}
		return secret, nil
	})
	*/
}
