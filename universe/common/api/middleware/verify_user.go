package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/common/api"
	"github.com/momentum-xyz/ubercontroller/utils"
)

func VerifyUser(log *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get jwt secret
		jwtSecret, ok := universe.GetNode().GetNodeAttributeValue(entry.NewAttributeID(universe.GetSystemPluginID(), universe.Attributes.Node.JWTKey.Name))
		if !ok || jwtSecret == nil {
			err := errors.New("Middleware: VerifyUser: failed to get jwt_key attribute")
			api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_jwt_key", err, log)
			return
		}
		secret := utils.GetFromAnyMap(*jwtSecret, "secret", "")

		// auth
		token, err := api.ValidateJWT(api.GetTokenFromRequest(c), []byte(secret))
		if err != nil {
			err = errors.WithMessage(err, "Middleware: VerifyUser: failed to verify token")
			api.AbortRequest(c, http.StatusForbidden, "failed_to_verify_access_token", err, log)
			return
		}
		c.Set(api.TokenContextKey, *token)
	}
}
