package common

import (
	"context"
	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var common struct {
	ctx context.Context
	log *zap.SugaredLogger
	cfg *config.Config
}

func Initialize(ctx context.Context) error {
	log := utils.GetFromAny(ctx.Value(types.LoggerContextKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.LoggerContextKey))
	}
	cfg := utils.GetFromAny(ctx.Value(types.ConfigContextKey), (*config.Config)(nil))
	if cfg == nil {
		return errors.Errorf("failed to get config from context: %T", ctx.Value(types.ConfigContextKey))
	}

	common.ctx = ctx
	common.log = log
	common.cfg = cfg

	return nil
}

func GetContext() context.Context {
	return common.ctx
}

func GetLogger() *zap.SugaredLogger {
	return common.log
}

func GetConfig() *config.Config {
	return common.cfg
}
