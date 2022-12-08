package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/database"
)

func AuthorizeAttributes(log *zap.SugaredLogger, db database.DB, attributeType any) gin.HandlerFunc {
	return func(c *gin.Context) {
		return
	}
}
