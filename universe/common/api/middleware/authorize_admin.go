package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/universe/common/api"
)

func AuthorizeAdmin(log *zap.SugaredLogger, db database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		spaceID, err := uuid.Parse(c.Param("spaceID"))
		if err != nil {
			err := errors.WithMessage(err, "Middleware: AuthorizeAdmin: failed to parse space id")
			api.AbortRequest(c, http.StatusBadRequest, "invalid_space_id", err, log)
			return
		}

		userID, err := api.GetUserIDFromContext(c)
		if err != nil {
			err := errors.WithMessage(err, "Middleware: AuthorizeAdmin: failed to get user id from context")
			api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_user_id", err, log)
			return
		}

		isAdmin, err := db.UserSpaceCheckIsUserIndirectSpaceAdmin(c, userID, spaceID)
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
