package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe/common/api"
	"github.com/momentum-xyz/ubercontroller/utils"
)

func AuthorizeAttributes(log *zap.SugaredLogger, db database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		type InBody struct {
			PluginID      string `json:"plugin_id" binding:"required"`
			AttributeName string `json:"attribute_name" binding:"required"`
		}

		inBody := InBody{}

		if err := c.ShouldBindJSON(&inBody); err != nil {
			err = errors.WithMessage(err, "Middleware: AuthorizeAttributes: failed to bind json")
			api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, log)
			return
		}

		userID, err := api.GetUserIDFromContext(c)
		if err != nil {
			err = errors.WithMessage(err, "Middleware: AuthorizeAttributes: failed to get user id from context")
			api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_user_id", err, log)
			return
		}

		spaceID, err := uuid.Parse(c.Param("spaceID"))
		if err != nil {
			err := errors.WithMessage(err, "Middleware: AuthorizeAttributes: failed to parse space id")
			api.AbortRequest(c, http.StatusBadRequest, "invalid_space_id", err, log)
			return
		}

		pluginID, err := uuid.Parse(inBody.PluginID)
		if err != nil {
			err := errors.WithMessage(err, "Middleware: AuthorizeAttributes: failed to parse plugin id")
			api.AbortRequest(c, http.StatusInternalServerError, "invalid_plugin_id", err, log)
			return
		}

		attributeID := entry.NewAttributeID(pluginID, inBody.AttributeName)
		attributeType, err := db.AttributeTypesGetAttributeTypeByID(c, attributeID)
		if err != nil {
			err := errors.WithMessage(err, "Middleware: AuthorizeAttributes: failed to get attribute type")
			api.AbortRequest(c, http.StatusBadRequest, "failed_to_get_attribute_type", err, log)
			return
		}

		options := attributeType.Options
		if options == nil {
			err := errors.WithMessage(err, "Middleware: AuthorizeAttributes: no options in attribute type")
			api.AbortRequest(c, http.StatusNotFound, "failed_to_get_options", err, log)
			return
		}

		permissions := utils.GetFromAnyMap(*options, "permissions", (map[string]any)(nil))
		if permissions == nil {
			err := errors.WithMessage(err, "Middleware: AuthorizeAttributes: no permissions in options")
			api.AbortRequest(c, http.StatusNotFound, "permissions_not_found", err, log)
			return
		}

		mutations := utils.GetFromAnyMap(*options, "permissions", (map[string]string)(nil))
		if mutations == nil {
			err := errors.WithMessage(err, "Middleware: AuthorizeAttributes: no mutations in permissions")
			api.AbortRequest(c, http.StatusNotFound, "mutations_not_found", err, log)
			return
		}

		userSpaceID := entry.NewUserSpaceID(userID, spaceID)
		userSpaceValue, err := db.UserSpaceGetUserSpaceValueByID(c, userSpaceID)
		if err != nil {
			err := errors.WithMessage(err, "Middleware: AuthorizeAttributes: failed to get userSpaceValue")
			api.AbortRequest(c, http.StatusBadRequest, "failed_to_get_user_space_value", err, log)
			return
		}

		if userSpaceValue == nil {
			err := errors.WithMessage(err, "Middleware: AuthorizeAttributes: no user space value found")
			api.AbortRequest(c, http.StatusNotFound, "mutations_not_found", err, log)
			return
		}

		role := utils.GetFromAnyMap(*userSpaceValue, "role", "")

		if mutations[role] == "" {
			return
		} else {
			c.Next()
		}
	}
}
