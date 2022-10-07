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
	Name        string   `json:"name"`
	Email       string   `json:"email"`
}

func VerifyToken(ctx context.Context, token string) (Token, error) {
	parsedToken, err := ParseToken(token)
	if err != nil {
		return parsedToken, errors.WithMessage(err, "failed to parse token")
	}

	oidcProvider, ok := api.oidcProviders.Load(parsedToken.Issuer)
	if !ok {
		provider, err := createProvider(parsedToken.Issuer)
		if err != nil {
			return parsedToken, errors.WithMessagef(err, "failed to create provider: %s", parsedToken.Issuer)
		}
		api.oidcProviders.Store(parsedToken.Issuer, provider)
		oidcProvider = provider
	}

	resp, err := rs.Introspect(ctx, oidcProvider, token)
	if err != nil {
		return parsedToken, errors.WithMessagef(err, "failed to introspect: %s", parsedToken.Issuer)
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

func ParseToken(tokenStr string) (Token, error) {
	var token Token

	parts := strings.Split(tokenStr, ".")
	if len(parts) < 2 {
		return token, errors.Errorf("invalid token, expected 2 parts got %d", len(parts))
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return token, errors.WithMessage(err, "invalid token payload")
	}
	if err := json.Unmarshal(payload, &token); err != nil {
		return token, errors.WithMessage(err, "failed to unmarshal payload")
	}

	return token, nil
}

func createProvider(issuer string) (rs.ResourceServer, error) {
	cfg := api.cfg.Auth

	provider, ok := utils.GetKeyByValueFromMap(cfg.OIDCURLs, issuer)
	if !ok {
		return nil, errors.Errorf("failed to get oidc provider by oidc url: %s", issuer)
	}

	clientID := cfg.OIDCClientIDs[provider]
	secret := cfg.OIDCSecrets[provider]
	introspectURL := cfg.OIDCIntospectURLs[provider]

	opts := make([]rs.Option, 0, 1)
	if introspectURL != "" {
		oidcConfig, err := client.Discover(issuer, http.DefaultClient)
		if err != nil {
			return nil, errors.WithMessagef(err, "failed to discover issuer: %s", issuer)
		}
		opts = append(opts, rs.WithStaticEndpoints(oidcConfig.TokenEndpoint, introspectURL))
	}

	oidcProvider, err := rs.NewResourceServerClientCredentials(issuer, clientID, secret, opts...)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to create resource server client: %s", issuer)
	}

	return oidcProvider, nil
}
