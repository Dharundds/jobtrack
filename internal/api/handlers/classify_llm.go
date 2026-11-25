package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/dharundds/jobtrack/internal/api/apperrors"
	"github.com/dharundds/jobtrack/internal/api/middleware"
	"github.com/dharundds/jobtrack/internal/store"
	syncworker "github.com/dharundds/jobtrack/internal/sync"
)

type ClassifyLLMHandler struct {
	Worker *syncworker.Worker
	Store  *store.Store
}

func (h *ClassifyLLMHandler) Pending(c *gin.Context) {
	ctx := c.Request.Context()
	n, err := h.Store.CountApplicationsNeedingLLMReview(ctx)
	if err != nil {
		middleware.AbortWithError(c, apperrors.Internal("count applications needing llm review"))
		return
	}
	unlinked, _ := h.Store.CountUnlinkedMessages(ctx)
	c.JSON(http.StatusOK, gin.H{"pending": n, "unlinked_messages": unlinked})
}

type classifyLLMBody struct {
	// Limit is how many pipeline applications to LLM-review (2 through pending count), newest activity first.
	Limit int `json:"limit" binding:"required,min=2"`
}

func (h *ClassifyLLMHandler) Trigger(c *gin.Context) {
	var body classifyLLMBody
	if err := c.ShouldBindJSON(&body); err != nil {
		middleware.AbortWithError(c, apperrors.BadRequest("invalid_limit", "limit must be an integer from 2 upward"))
		return
	}

	ctx := c.Request.Context()
	pending, err := h.Store.CountApplicationsNeedingLLMReview(ctx)
	if err != nil {
		middleware.AbortWithError(c, apperrors.Internal("count applications needing llm review"))
		return
	}
	if pending == 0 {
		middleware.AbortWithError(c, apperrors.BadRequest("no_pending", "no pipeline applications need LLM review"))
		return
	}
	limit := body.Limit
	if limit > pending {
		middleware.AbortWithError(c, apperrors.BadRequest("limit_too_large", "limit cannot exceed pending applications"))
		return
	}

	started, status := h.Worker.TriggerLLMClassify(ctx, limit)
	if !started {
		c.JSON(http.StatusAccepted, gin.H{
			"ok":      true,
			"message": "llm classify already running or sync in progress",
			"status":  status,
		})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"ok": true, "message": "llm classify started"})
}