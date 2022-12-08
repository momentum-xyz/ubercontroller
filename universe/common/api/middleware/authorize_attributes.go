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

		attributeID := entry.NewAttributeTypeID(pluginID, inBody.AttributeName)
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

		role := (*options)["permissions"]
		if role == nil {
			err := errors.WithMessage(err, "Middleware: AuthorizeAttributes: no role in permissions")
			api.AbortRequest(c, http.StatusNotFound, "failed_to_get_permissions", err, log)
			return
		}

		userSpaceID := entry.NewUserSpaceID(userID, spaceID)
		_, err = db.UserSpaceGetUserSpaceValueByID(c, userSpaceID)
		if err != nil {
			err := errors.WithMessage(err, "Middleware: AuthorizeAttributes: failed to get userSpaceValue")
			api.AbortRequest(c, http.StatusBadRequest, "failed_to_get_user_space_value", err, log)
			return
		}

		//if userRole != role {
		//	return
		//} else {
		//	c.Next()
		//}
	}
}
