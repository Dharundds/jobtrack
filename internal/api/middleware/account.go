package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/dharundds/jobtrack/internal/api/apperrors"
	"github.com/dharundds/jobtrack/internal/store"
)

const ContextAccountKey = "account"

func RequireAccount(st *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		acc, err := st.GetAccount(c.Request.Context())
		if err != nil {
			if err == store.ErrNoAccount {
				AbortWithError(c, apperrors.BadRequest("not_connected", "connect Gmail before using this feature"))
				return
			}
			AbortWithError(c, apperrors.Internal("failed to load account"))
			return
		}
		if acc.EncRefreshToken == "" {
			AbortWithError(c, apperrors.BadRequest("not_connected", "connect Gmail before using this feature"))
			return
		}
		c.Set(ContextAccountKey, acc)
		c.Next()
	}
}

func AccountFromContext(c *gin.Context) (*store.Account, bool) {
	v, ok := c.Get(ContextAccountKey)
	if !ok {
		return nil, false
	}
	acc, ok := v.(*store.Account)
	return acc, ok
}
