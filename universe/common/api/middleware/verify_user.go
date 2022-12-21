package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/types/generic"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/common/api"
)

var secretStore = generic.NewSyncMap[string, *[]byte](0)

func VerifyUser(log *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		secret, err := GetJWTSecret()
		if err != nil {
			err = errors.WithMessage(err, "Middleware: VerifyUser: failed to get jwt secret")
			api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_secret", err, log)
			return
		}

		token, err := api.ValidateJWTWithSecret(api.GetTokenFromRequest(c), secret)
		if err != nil {
			err = errors.WithMessage(err, "Middleware: VerifyUser: failed to verify token")
			api.AbortRequest(c, http.StatusForbidden, "failed_to_verify_access_token", err, log)
			return
		}
		c.Set(api.TokenContextKey, *token)
	}
}

func GetJWTSecret() ([]byte, error) {
	var secret []byte

	store, ok := secretStore.Load(universe.Attributes.Node.JWTKey.Key)
	if !ok {
		jwtSecret, err := api.GetJWTSecret()
		if err != nil {
			return nil, errors.New("failed to get jwt secret")
		}

		secretStore.Store(universe.Attributes.Node.JWTKey.Key, &jwtSecret)
		secret = jwtSecret
	} else {
		secret = *store
	}

	return secret, nil
}
