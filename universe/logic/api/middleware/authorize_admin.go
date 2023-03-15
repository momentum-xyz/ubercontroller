package middleware

import (
	"github.com/momentum-xyz/ubercontroller/utils/mid"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
)

func AuthorizeAdmin(log *zap.SugaredLogger) gin.HandlerFunc {
	userObjects := universe.GetNode().GetUserObjects()

	return func(c *gin.Context) {
		objectID, err := mid.Parse(c.Param("objectID"))
		if err != nil {
			err := errors.WithMessage(err, "Middleware: AuthorizeAdmin: failed to parse object mid")
			api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, log)
			return
		}

		userID, err := api.GetUserIDFromContext(c)
		if err != nil {
			err := errors.WithMessage(err, "Middleware: AuthorizeAdmin: failed to get user mid from context")
			api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_user_id", err, log)
			return
		}

		isAdmin, err := userObjects.CheckIsIndirectAdmin(entry.NewUserObjectID(userID, objectID))
		if err != nil {
			err := errors.WithMessage(err, "Middleware: AuthorizeAdmin: failed to check is indirect admin")
			api.AbortRequest(c, http.StatusInternalServerError, "check_failed", err, log)
			return
		}

		if !isAdmin {
			err := errors.New("Middleware: AuthorizeAdmin: user is not admin")
			api.AbortRequest(c, http.StatusForbidden, "not_admin", err, log)
			return
		}
	}
}
