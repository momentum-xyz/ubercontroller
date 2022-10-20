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
			err = errors.WithMessage(err, "Middleware: VerifyUser: failed to verify token")
			log.Debug(err)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": map[string]string{"reason": "failed_to_verify_access_token", "message": err.Error()},
			})
			return
		}
		c.Set(api.TokenContextKey, token)
	}
}
