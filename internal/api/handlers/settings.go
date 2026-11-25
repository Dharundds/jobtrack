package handlers

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/dharundds/jobtrack/internal/api/apperrors"
	"github.com/dharundds/jobtrack/internal/api/middleware"
	"github.com/dharundds/jobtrack/internal/config"
	"github.com/dharundds/jobtrack/internal/llm"
	"github.com/dharundds/jobtrack/internal/store"
)

type SettingsHandler struct {
	Store *store.Store
	Cfg   *config.Config
}

type settingsResponse struct {
	SyncIntervalMinutes int    `json:"sync_interval_minutes"`
	GmailQuery          string `json:"gmail_query"`
	LLMBaseURL          string `json:"llm_base_url"`
	LLMModel            string `json:"llm_model"`
	LLMAPIKeySet        bool   `json:"llm_api_key_set"`
	LLMConcurrency      int    `json:"llm_concurrency"`
}

func (h *SettingsHandler) loadResponse(ctx context.Context) (settingsResponse, error) {
	gmailQ, _ := h.Store.GmailQuery(ctx)
	intervalRaw, _ := h.Store.GetSetting(ctx, "sync_interval_minutes")
	interval := store.ParseSyncIntervalMinutes(intervalRaw, h.Cfg.SyncIntervalMin)
	base, _ := h.Store.GetSetting(ctx, store.SettingLLMBaseURL)
	model, _ := h.Store.GetSetting(ctx, store.SettingLLMModel)
	keySet, _ := h.Store.LLMAPIKeyConfigured(ctx)
	conc, _ := h.Store.LLMConcurrency(ctx, 2)
	return settingsResponse{
		SyncIntervalMinutes: interval,
		GmailQuery:          gmailQ,
		LLMBaseURL:          base,
		LLMModel:            model,
		LLMAPIKeySet:        keySet,
		LLMConcurrency:      conc,
	}, nil
}

func (h *SettingsHandler) Get(c *gin.Context) {
	ctx := c.Request.Context()
	out, err := h.loadResponse(ctx)
	if err != nil {
		middleware.AbortWithError(c, apperrors.Internal("failed to load settings"))
		return
	}
	c.JSON(http.StatusOK, out)
}

type patchSettingsRequest struct {
	SyncIntervalMinutes *int    `json:"sync_interval_minutes" binding:"omitempty,min=1,max=1440"`
	GmailQuery          *string `json:"gmail_query" binding:"omitempty,min=1"`
	LLMBaseURL          *string `json:"llm_base_url"`
	LLMModel            *string `json:"llm_model"`
	LLMAPIKey           *string `json:"llm_api_key"`
	LLMConcurrency      *int    `json:"llm_concurrency" binding:"omitempty,min=1,max=8"`
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
	if req.LLMConcurrency != nil {
		if err := h.Store.SetSetting(ctx, store.SettingLLMConcurrency, strconv.Itoa(*req.LLMConcurrency)); err != nil {
			middleware.AbortWithError(c, apperrors.Internal("failed to save llm settings"))
			return
		}
	}
	if req.LLMBaseURL != nil || req.LLMModel != nil || req.LLMAPIKey != nil {
		cur, err := h.Store.LLMConfig(ctx, h.Cfg.TokenEncKey)
		if err != nil {
			middleware.AbortWithError(c, apperrors.Internal("failed to load llm settings"))
			return
		}
		base := cur.BaseURL
		model := cur.Model
		apiKey := cur.APIKey
		updateKey := false
		if req.LLMBaseURL != nil {
			base = strings.TrimSpace(*req.LLMBaseURL)
		}
		if req.LLMModel != nil {
			model = strings.TrimSpace(*req.LLMModel)
		}
		if req.LLMAPIKey != nil {
			updateKey = true
			apiKey = *req.LLMAPIKey
		}
		if err := h.Store.SetLLMConfig(ctx, h.Cfg.TokenEncKey, llm.NormalizeBaseURL(base), model, apiKey, updateKey); err != nil {
			middleware.AbortWithError(c, apperrors.Internal("failed to save llm settings"))
			return
		}
	}
	h.Get(c)
}

func (h *SettingsHandler) TestLLM(c *gin.Context) {
	ctx := c.Request.Context()
	cfg, err := h.Store.LLMConfig(ctx, h.Cfg.TokenEncKey)
	if err != nil {
		middleware.AbortWithError(c, apperrors.Internal("failed to load llm settings"))
		return
	}
	if !cfg.Configured() {
		middleware.AbortWithError(c, apperrors.BadRequest("llm_not_configured", "set base URL and model first"))
		return
	}
	client := llm.NewClient(cfg)
	ok := client.Ping(ctx)
	c.JSON(http.StatusOK, gin.H{"ok": ok, "reachable": ok})
}