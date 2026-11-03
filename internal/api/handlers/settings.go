package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/dharundds/jobtrack/internal/api/apperrors"
	"github.com/dharundds/jobtrack/internal/api/middleware"
	"github.com/dharundds/jobtrack/internal/config"
	"github.com/dharundds/jobtrack/internal/store"
)

type SettingsHandler struct {
	Store *store.Store
	Cfg   *config.Config
}

func (h *SettingsHandler) Get(c *gin.Context) {
	ctx := c.Request.Context()
	gmailQ, _ := h.Store.GmailQuery(ctx)
	intervalRaw, _ := h.Store.GetSetting(ctx, "sync_interval_minutes")
	interval := store.ParseSyncIntervalMinutes(intervalRaw, h.Cfg.SyncIntervalMin)
	c.JSON(http.StatusOK, gin.H{
		"sync_interval_minutes": interval,
		"gmail_query":           gmailQ,
	})
}

type patchSettingsRequest struct {
	SyncIntervalMinutes *int    `json:"sync_interval_minutes" binding:"omitempty,min=1,max=1440"`
	GmailQuery          *string `json:"gmail_query" binding:"omitempty,min=1"`
}

func (h *SettingsHandler) Patch(c *gin.Context) {
	var req patchSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithError(c, apperrors.BadRequest("invalid_body", err.Error()))
		return
	}
	ctx := c.Request.Context()
	if req.SyncIntervalMinutes != nil {
		if err := h.Store.SetSetting(ctx, "sync_interval_minutes", strconv.Itoa(*req.SyncIntervalMinutes)); err != nil {
			middleware.AbortWithError(c, apperrors.Internal("failed to save settings"))
			return
		}
	}
	if req.GmailQuery != nil {
		if err := h.Store.SetSetting(ctx, "gmail_query", *req.GmailQuery); err != nil {
			middleware.AbortWithError(c, apperrors.Internal("failed to save settings"))
			return
		}
	}
	h.Get(c)
}
