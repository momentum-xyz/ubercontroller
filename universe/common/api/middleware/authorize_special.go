package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AuthorizeSpecial(log *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		return
	}
}
