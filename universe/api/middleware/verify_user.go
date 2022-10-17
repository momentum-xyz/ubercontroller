package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/universe/api"
)

func VerifyUser(log *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := api.VerifyToken(c, api.GetTokenFromRequest(c))
		if err != nil {
			log.Debug(errors.WithMessage(err, "Middleware: VerifyUser: failed to verify token"))
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "failed to verify access token",
			})
			return
		}
		c.Set(api.TokenContextKey, token)
	}
}
