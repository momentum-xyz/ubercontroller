package generic

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/types"
)

// TODO: why-o-why? get rid of this global.
var generic struct {
	ctx context.Context
	log *zap.SugaredLogger
}

func Initialize(ctx types.LoggerContext) error {
	log := ctx.Logger()
	if log == nil {
		return errors.New("failed to get logger from context")
	}

	generic.ctx = ctx
	generic.log = log

	return nil
}
