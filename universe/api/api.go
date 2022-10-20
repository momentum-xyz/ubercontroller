package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/zitadel/oidc/pkg/client/rs"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/generic"
	"github.com/momentum-xyz/ubercontroller/utils"
)

var api = struct {
	ctx           context.Context
	log           *zap.SugaredLogger
	cfg           *config.Config
	oidcProviders *generic.SyncMap[string, rs.ResourceServer]
}{}

func Initialize(ctx context.Context, cfg *config.Config) error {
	log := utils.GetFromAny(ctx.Value(types.LoggerContextKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.LoggerContextKey))
	}

	api.ctx = ctx
	api.log = log
	api.cfg = cfg
	api.oidcProviders = generic.NewSyncMap[string, rs.ResourceServer]()

	return nil
}

func AbortRequest(c *gin.Context, code int, reason string, err error) {
	c.AbortWithStatusJSON(code, gin.H{
		"error": map[string]string{"reason": reason, "message": err.Error()},
	})
}
