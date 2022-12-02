package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/universe/common/api"
	"github.com/momentum-xyz/ubercontroller/universe/common/auth"
)

func VerifyUser(log *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		claim, err := auth.GetJWTClaimFromContext(c)
		if err != nil {
			err = errors.WithMessage(err, "Middleware: VerifyUser: failed to get claims from context")
			api.AbortRequest(c, http.StatusForbidden, "failed_to_get_jwt_claim", err, log)
			return
		}
		token, err := auth.ValidateToken(claim)
		if err != nil {
			err = errors.WithMessage(err, "Middleware: VerifyUser: failed to verify token")
			api.AbortRequest(c, http.StatusForbidden, "failed_to_verify_access_token", err, log)
			return
		}
		c.Set(api.TokenContextKey, token)
	}
}
