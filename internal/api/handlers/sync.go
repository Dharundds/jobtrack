package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	syncworker "github.com/dharundds/jobtrack/internal/sync"
)

type SyncHandler struct {
	Worker *syncworker.Worker
}

func (h *SyncHandler) Trigger(c *gin.Context) {
	started, status := h.Worker.TriggerManual(c.Request.Context())
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
