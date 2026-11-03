package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/dharundds/jobtrack/internal/logging"
	"github.com/rs/zerolog"
)

const (
	ContextLoggerKey = "logger"
	ContextRequestID = "request_id"
	HeaderRequestID  = "X-Request-ID"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(HeaderRequestID)
		if id == "" {
			id = uuid.New().String()
		}
		c.Header(HeaderRequestID, id)
		c.Set(ContextRequestID, id)

		reqLog := logging.L().With().
			Str("request_id", id).
			Logger()
		c.Set(ContextLoggerKey, reqLog)
		c.Next()
	}
}

func LoggerFromContext(c *gin.Context) zerolog.Logger {
	if v, ok := c.Get(ContextLoggerKey); ok {
		if l, ok := v.(zerolog.Logger); ok {
			return l
		}
	}
	return logging.L()
}
