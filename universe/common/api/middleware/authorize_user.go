package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AuthorizeUser(log *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get userId
		// check user_space recursively for rights
		// check if user role is present
	}
}
