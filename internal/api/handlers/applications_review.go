package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/dharundds/jobtrack/internal/api/apperrors"
	"github.com/dharundds/jobtrack/internal/api/middleware"
	"github.com/dharundds/jobtrack/internal/llm"
	"github.com/dharundds/jobtrack/internal/logging"
	"github.com/dharundds/jobtrack/internal/store"
)

// ReviewWithLLM re-checks a listing using the local LLM; may update fields or delete the row.
func (h *ApplicationsHandler) ReviewWithLLM(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		middleware.AbortWithError(c, apperrors.BadRequest("invalid_id", "invalid application id"))
		return
	}
	ctx := c.Request.Context()

	cfg, err := h.Store.LLMConfig(ctx, h.Cfg.TokenEncKey)
	if err != nil || !cfg.Configured() {
		middleware.AbortWithError(c, apperrors.BadRequest("llm_not_configured", "configure LLM in Settings first"))
		return
	}
	client := llm.NewClient(cfg)
	if !client.Ping(ctx) {
		middleware.AbortWithError(c, apperrors.BadRequest("llm_unavailable", "LLM server not reachable; start LM Studio and try again"))
		return
	}

	acc, err := h.Store.GetAccount(ctx)
	if err != nil {
		middleware.AbortWithError(c, apperrors.BadRequest("not_connected", "connect Gmail first"))
		return
	}

	if _, err := h.Store.GetApplicationFull(ctx, id); err != nil {
		if err == store.ErrApplicationNotFound {
			middleware.AbortWithError(c, apperrors.NotFound("application not found"))
			return
		}
		middleware.AbortWithError(c, apperrors.Internal("failed to load application"))
		return
	}

	conc, _ := h.Store.LLMConcurrency(ctx, 2)
	log := logging.L().With().Str("component", "review-llm").Logger()
	deleted, err := llm.ReviewApplicationWithLLM(ctx, h.Store, log, client, acc.Email, id, conc)
	if err != nil {
		middleware.AbortWithError(c, apperrors.BadRequest("llm_review_failed", err.Error()))
		return
	}

	if deleted {
		c.JSON(http.StatusOK, gin.H{
			"action":              "deleted",
			"application_deleted": true,
		})
		return
	}

	updated, err := h.Store.GetApplicationFull(ctx, id)
	if err != nil {
		middleware.AbortWithError(c, apperrors.Internal("failed to load application"))
		return
	}
	j, err := h.loadApplicationJSON(c, *updated)
	if err != nil {
		middleware.AbortWithError(c, apperrors.Internal("failed to load application"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"action":              "updated",
		"application":         j,
		"application_deleted": false,
	})
}