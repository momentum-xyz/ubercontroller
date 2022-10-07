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
	"github.com/zitadel/oidc/pkg/oidc"
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

func VerifyToken(ctx context.Context, tokenStr string) (oidc.IntrospectionResponse, error) {
	token, err := ParseToken(tokenStr)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to parse token")
	}

	oidcProvider, ok := api.oidcProviders.Load(token.Issuer)
	if !ok {
		provider, err := createProvider(token.Issuer)
		if err != nil {
			return nil, errors.WithMessagef(err, "failed to create provider: %s", token.Issuer)
		}
		api.oidcProviders.Store(token.Issuer, provider)
		oidcProvider = provider
	}

	resp, err := rs.Introspect(ctx, oidcProvider, tokenStr)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to introspect: %s", token.Issuer)
	}

	if !resp.IsActive() {
		return nil, errors.Errorf("token is not active: %s", tokenStr)
	}

	return resp, nil
}

func GetTokenFromRequest(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	return strings.Replace(authHeader, "Bearer ", "", -1)
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
	introspectURL := cfg.OIDCIntospectURLs[issuer]

	opts := make([]rs.Option, 0, 1)
	if introspectURL != "" {
		oidcConfig, err := client.Discover(issuer, http.DefaultClient)
		if err != nil {
			return nil, errors.WithMessagef(err, "failed to discover issuer: %s", issuer)
		}
		opts = append(opts, rs.WithStaticEndpoints(oidcConfig.TokenEndpoint, introspectURL))
	}

	provider, err := rs.NewResourceServerClientCredentials(
		issuer, cfg.OIDCClientIDs[issuer], cfg.OIDCSecrets[issuer], opts...,
	)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to create resource server client: %s", issuer)
	}

	return provider, nil
}
