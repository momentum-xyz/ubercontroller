package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AuthorizeAdmin(log *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get userId
		// check user_space recursively for rights
		// check if admin role is present
	}
}
