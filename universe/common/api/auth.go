package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
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

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/utils"
)

type Token struct {
	Guest struct {
		IsGuest bool `json:"1"`
	}
	Issuer      string   `json:"iss"`
	Subject     string   `json:"sub"`
	Audience    []string `json:"aud"`
	Expiry      int      `json:"exp"`
	IssuedAt    int      `json:"iat"`
	Web3Address string   `json:"web3_address"`
	Web3Type    string   `json:"web3_type"`
	RawToken    string   `json:"-"`
}

func VerifyToken(ctx context.Context, token string) (Token, error) {
	parsedToken, err := ParseToken(token)
	if err != nil {
		return parsedToken, errors.WithMessage(err, "failed to parse token")
	}

	// TODO: change this!
	return parsedToken, nil
	//for _, provider := range api.cfg.Auth.OIDCProviders {
	//	if err := verifyTokenByProvider(ctx, provider, parsedToken); err == nil {
	//		return parsedToken, nil
	//	}
	//}
	//
	//return parsedToken, errors.Errorf("failed to verify token: %s", parsedToken.RawToken)
}

func verifyTokenByProvider(ctx context.Context, provider string, token Token) error {
	oidcProvider, ok := api.oidcProviders.Load(provider)
	if !ok {
		newProvider, err := createProvider(provider)
		if err != nil {
			return errors.WithMessagef(err, "failed to create provider: %s", provider)
		}
		api.oidcProviders.Store(provider, newProvider)
		oidcProvider = newProvider
	}

	resp, err := rs.Introspect(ctx, oidcProvider, token.RawToken)
	if err != nil {
		return errors.WithMessagef(err, "failed to introspect: %s", provider)
	}

	if !resp.IsActive() {
		return errors.Errorf("token is not active: %s", token.RawToken)
	}

	return nil
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

// TODO: Remove this when old tokens are removed
// Made as a second method here so we can use the jwt tokens
// and still not break node.apiCheckTokens - when the old auth
// flow is completely removed this method should become the
// only GetTokenFromContext
func GetJWTFromContext(c *gin.Context) (jwt.Token, error) {
	value, ok := c.Get(JWTContextKey)
	if !ok {
		return jwt.Token{}, errors.Errorf("failed to get token value from context")
	}

	token := utils.GetFromAny(value, jwt.Token{})

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
func SignJWTToken(userID string, secret []byte) (*entry.Token, error) {
	claims := jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(4 * time.Hour).Unix(),
		Issuer:    "ubercontroller",
		Subject:   userID,
	}

	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedString, err := jwt.SignedString(secret)
	if err != nil {
		return nil, err
	}

	exp := time.Unix(claims.ExpiresAt, 0)

	token := &entry.Token{
		SignedString: &signedString,
		Subject:      &claims.Subject,
		Issuer:       &claims.Issuer,
		IssuedAt:     time.Unix(claims.IssuedAt, 0),
		ExpiresAt:    &exp,
	}

	return token, nil
}

func ValidateJWT(signedString string, secret []byte) (*jwt.Token, error) {
	return jwt.Parse(signedString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid token %v", token.Header["alg"])

		}
		return secret, nil
	})
}
