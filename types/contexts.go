package types

import (
	"context"
	"errors"

	"github.com/momentum-xyz/ubercontroller/config"
	"go.uber.org/zap"
)

// Typed contex.Context handling

const (
	loggerContextKey = "logger"
	configContextKey = "config"
)

type LoggerContext interface {
	context.Context
	Logger() *zap.SugaredLogger
}

type ConfigContext interface {
	context.Context
	Config() *config.Config
}

type NodeContext interface {
	LoggerContext
	ConfigContext
}

type nodeContext struct {
	context.Context
	log *zap.SugaredLogger
	cfg *config.Config
}

func (n nodeContext) Logger() *zap.SugaredLogger {
	return n.log
}

func (n nodeContext) Config() *config.Config {
	return n.cfg
}

func NewNodeContext(ctx context.Context, log *zap.SugaredLogger, cfg *config.Config) (NodeContext, error) {
	if log == nil {
		return nil, errors.New("NodeContext: log required")
	}
	if cfg == nil {
		return nil, errors.New("NodeContext: cfg required")
	}
	return nodeContext{
		Context: ctx,
		log:     log,
		cfg:     cfg,
	}, nil
}
