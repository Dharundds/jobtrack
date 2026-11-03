package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/dharundds/jobtrack/internal/api/apperrors"
	"github.com/dharundds/jobtrack/internal/api/middleware"
	"github.com/dharundds/jobtrack/internal/store"
)

type StatsHandler struct {
	Store *store.Store
}

func (h *StatsHandler) Get(c *gin.Context) {
	st, err := h.Store.ComputeStats(c.Request.Context())
	if err != nil {
		middleware.AbortWithError(c, apperrors.Internal("failed to compute stats"))
		return
	}
	total := 0
	for _, n := range st.ByStatus {
		total += n
	}
	c.JSON(http.StatusOK, gin.H{
		"total_applications":       total,
		"by_status":                st.ByStatus,
		"response_rate":            st.ResponseRate,
		"avg_days_to_first_reply":  st.AvgDaysToFirstReply,
	})
}
