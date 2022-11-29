package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/utils"
)

var api = struct {
	ctx context.Context
	log *zap.SugaredLogger
	cfg *config.Config
}{}

type HTTPErrorPayload struct {
	Reason  string `json:"reason"`
	Message string `json:"message"`
}
type HTTPError struct {
	Error HTTPErrorPayload `json:"error"`
}

func Initialize(ctx context.Context, cfg *config.Config) error {
	log := utils.GetFromAny(ctx.Value(types.LoggerContextKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.LoggerContextKey))
	}

	api.ctx = ctx
	api.log = log
	api.cfg = cfg

	return nil
}

func AbortRequest(c *gin.Context, code int, reason string, err error, log *zap.SugaredLogger) {
	if code == http.StatusInternalServerError {
		log.Error(err)
	} else {
		log.Debug(err)
	}
	c.AbortWithStatusJSON(code, &HTTPError{Error: HTTPErrorPayload{
		Reason:  reason,
		Message: err.Error(),
	}})
}
