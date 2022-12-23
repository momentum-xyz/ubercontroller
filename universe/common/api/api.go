package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HTTPErrorPayload struct {
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

type HTTPError struct {
	Error HTTPErrorPayload `json:"error"`
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
