package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/zitadel/oidc/pkg/client"
	"github.com/zitadel/oidc/pkg/client/rs"
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
	providerName := "web3"
	if parsedToken.Guest.IsGuest {
		providerName = "guest"
	}

	oidcProvider, ok := api.oidcProviders.Load(providerName)
	if !ok {
		provider, err := createProvider(providerName)
		if err != nil {
			return parsedToken, errors.WithMessagef(err, "failed to create provider: %s", providerName)
		}
		api.oidcProviders.Store(providerName, provider)
		oidcProvider = provider
	}

	resp, err := rs.Introspect(ctx, oidcProvider, token)
	if err != nil {
		return parsedToken, errors.WithMessagef(err, "failed to introspect: %s", providerName)
	}

	if !resp.IsActive() {
		return parsedToken, errors.Errorf("token is not active: %s", token)
	}

	return parsedToken, nil
}

func GetTokenFromRequest(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	return strings.TrimPrefix(authHeader, "Bearer ")
}

func ParseToken(token string) (Token, error) {
	var parsedToken Token

	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return parsedToken, errors.Errorf("invalid token, expected 2 parts got %d", len(parts))
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return parsedToken, errors.WithMessage(err, "invalid token payload")
	}
	if err := json.Unmarshal(payload, &token); err != nil {
		return parsedToken, errors.WithMessage(err, "failed to unmarshal payload")
	}

	parsedToken.RawToken = token

	return parsedToken, nil
}

func createProvider(provider string) (rs.ResourceServer, error) {
	cfg := api.cfg.Auth
	oidcURL := cfg.OIDCURL
	clientID := cfg.GetIDByProvider(provider)
	secret := cfg.GetSecretByProvider(provider)
	introspectURL := cfg.GetIntrospectURLByProvider(provider)

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
