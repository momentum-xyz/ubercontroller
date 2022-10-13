package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"

	"github.com/momentum-xyz/ubercontroller/universe/api"
)

func VerifyUser(log *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := api.VerifyToken(c, api.GetTokenFromRequest(c)); err != nil {
			log.Debug(errors.WithMessage(err, "Middleware: VerifyUser: failed to verify user"))
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "failed to verify access token",
			})
		}
	}
}
