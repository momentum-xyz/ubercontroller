package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/universe/common/api"
)

func VerifyUser(log *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {

		var secret []byte
		token, err := api.ValidateJWT(api.GetTokenFromRequest(c), secret)
		if err != nil {
			err = errors.WithMessage(err, "Middleware: VerifyUser: failed to verify token")
			api.AbortRequest(c, http.StatusForbidden, "failed_to_verify_access_token", err, log)
			return
		}
		c.Set(api.JWTContextKey, token)
	}
}
