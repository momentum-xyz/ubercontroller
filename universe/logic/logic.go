package logic

import (
	"context"
	"errors"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/types"
	"go.uber.org/zap"
)

var logic struct {
	ctx context.Context
	log *zap.SugaredLogger
	cfg *config.Config
}

func Initialize(ctx interface {
	context.Context
	types.LoggerContext
	types.ConfigContext
}) error {
	log := ctx.Logger()
	if log == nil {
		return errors.New("failed to get logger from context")
	}
	cfg := ctx.Config()
	if cfg == nil {
		return errors.New("failed to get config from context")
	}

	logic.ctx = ctx
	logic.log = log
	logic.cfg = cfg

	return nil
}

func GetContext() context.Context {
	return logic.ctx
}

func GetLogger() *zap.SugaredLogger {
	return logic.log
}

func GetConfig() *config.Config {
	return logic.cfg
}
