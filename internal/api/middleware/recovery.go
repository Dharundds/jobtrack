package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log := LoggerFromContext(c)
				log.Error().
					Interface("panic", r).
					Bytes("stack", debug.Stack()).
					Msg("panic recovered")

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": gin.H{
						"code":    "internal",
						"message": "internal server error",
					},
				})
			}
		}()
		c.Next()
	}
}
