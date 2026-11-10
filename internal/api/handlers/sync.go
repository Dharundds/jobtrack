package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/dharundds/jobtrack/internal/api/apperrors"
	"github.com/dharundds/jobtrack/internal/api/middleware"
	syncworker "github.com/dharundds/jobtrack/internal/sync"
)

type SyncHandler struct {
	Worker *syncworker.Worker
}

func (h *SyncHandler) Trigger(c *gin.Context) {
	ctx := c.Request.Context()
	if strings.EqualFold(c.Query("full"), "true") || strings.EqualFold(c.Query("full"), "1") {
		if err := h.Worker.ResetSyncHistory(ctx); err != nil {
			middleware.AbortWithError(c, apperrors.Internal("failed to reset sync history"))
			return
		}
	}
	started, status := h.Worker.TriggerManual(ctx)
	if !started {
		c.JSON(http.StatusAccepted, gin.H{
			"ok":      true,
			"message": "sync already running",
			"status":  status,
		})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"ok": true, "message": "sync started"})
}

func (h *SyncHandler) Status(c *gin.Context) {
	st := h.Worker.StatusSnapshot(c.Request.Context())
	c.JSON(http.StatusOK, st)
}
