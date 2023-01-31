package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
)

func VerifyUser(log *zap.SugaredLogger) gin.HandlerFunc {
	var secret []byte

	return func(c *gin.Context) {
		if secret == nil {
			jwtSecret, err := api.GetJWTSecret()
			if err != nil {
				err = errors.WithMessage(err, "Middleware: VerifyUser: failed to fetch jwt secret")
				api.AbortRequest(c, http.StatusInternalServerError, "failed_to_verify_access_token", err, log)
				return
			}

			secret = jwtSecret
		}

		token, err := api.ValidateJWTWithSecret(api.GetTokenFromRequest(c), secret)
		if err != nil {
			err = errors.WithMessage(err, "Middleware: VerifyUser: failed to verify token")
			api.AbortRequest(c, http.StatusUnauthorized, "failed_to_verify_access_token", err, log)
			return
		}
		c.Set(api.TokenContextKey, *token)
	}
}
