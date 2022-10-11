package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/momentum-xyz/ubercontroller/universe/api"
)

func VerifyUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := api.VerifyToken(c, api.GetTokenFromRequest(c)); err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": err.Error(),
			})
		}
	}
}
