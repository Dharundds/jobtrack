package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/dharundds/jobtrack/internal/api/apperrors"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}
		err := c.Errors.Last().Err

		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			c.JSON(appErr.HTTPStatus, gin.H{
				"error": gin.H{
					"code":    appErr.Code,
					"message": appErr.Message,
				},
			})
			return
		}

		log := LoggerFromContext(c)
		log.Error().Err(err).Msg("unhandled error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "internal",
				"message": "internal server error",
			},
		})
	}
}

func AbortWithError(c *gin.Context, err error) {
	_ = c.Error(err)
	c.Abort()
}
