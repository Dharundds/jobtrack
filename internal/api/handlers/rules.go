package handlers

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/dharundds/jobtrack/internal/api/apperrors"
	"github.com/dharundds/jobtrack/internal/api/middleware"
	"github.com/dharundds/jobtrack/internal/api/present"
	"github.com/dharundds/jobtrack/internal/store"
)

type RulesHandler struct {
	Store *store.Store
}

func (h *RulesHandler) List(c *gin.Context) {
	rules, err := h.Store.ListAllRules(c.Request.Context())
	if err != nil {
		middleware.AbortWithError(c, apperrors.Internal("failed to list rules"))
		return
	}
	out := make([]map[string]any, len(rules))
	for i, r := range rules {
		out[i] = present.RuleFromStore(r)
	}
	c.JSON(http.StatusOK, gin.H{"rules": out})
}

type createRuleRequest struct {
	Field     string `json:"field" binding:"required,oneof=subject body from"`
	Pattern   string `json:"pattern" binding:"required"`
	EventType string `json:"event_type" binding:"required"`
	Priority  int    `json:"priority"`
	Enabled   *bool  `json:"enabled"`
}

func (h *RulesHandler) Create(c *gin.Context) {
	var req createRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithError(c, apperrors.BadRequest("invalid_body", err.Error()))
		return
	}
	if _, err := regexp.Compile(req.Pattern); err != nil {
		middleware.AbortWithError(c, apperrors.BadRequest("invalid_pattern", "pattern must be valid regex"))
		return
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	id, err := h.Store.InsertRule(c.Request.Context(), req.Field, req.Pattern, req.EventType, req.Priority, enabled)
	if err != nil {
		middleware.AbortWithError(c, apperrors.Internal("failed to create rule"))
		return
	}
	r, err := h.Store.GetRule(c.Request.Context(), id)
	if err != nil {
		middleware.AbortWithError(c, apperrors.Internal("failed to load rule"))
		return
	}
	c.JSON(http.StatusCreated, present.RuleFromStore(*r))
}

type patchRuleRequest struct {
	Field     *string `json:"field" binding:"omitempty,oneof=subject body from"`
	Pattern   *string `json:"pattern"`
	EventType *string `json:"event_type"`
	Priority  *int    `json:"priority"`
	Enabled   *bool   `json:"enabled"`
}

func (h *RulesHandler) Patch(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		middleware.AbortWithError(c, apperrors.BadRequest("invalid_id", "invalid rule id"))
		return
	}
	var req patchRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithError(c, apperrors.BadRequest("invalid_body", err.Error()))
		return
	}
	if req.Pattern != nil {
		if _, err := regexp.Compile(*req.Pattern); err != nil {
			middleware.AbortWithError(c, apperrors.BadRequest("invalid_pattern", "pattern must be valid regex"))
			return
		}
	}
	if err := h.Store.UpdateRule(c.Request.Context(), id, req.Field, req.Pattern, req.EventType, req.Priority, req.Enabled); err != nil {
		if err.Error() == "rule not found" {
			middleware.AbortWithError(c, apperrors.NotFound("rule not found"))
			return
		}
		middleware.AbortWithError(c, apperrors.Internal("failed to update rule"))
		return
	}
	r, err := h.Store.GetRule(c.Request.Context(), id)
	if err != nil {
		middleware.AbortWithError(c, apperrors.NotFound("rule not found"))
		return
	}
	c.JSON(http.StatusOK, present.RuleFromStore(*r))
}

func (h *RulesHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		middleware.AbortWithError(c, apperrors.BadRequest("invalid_id", "invalid rule id"))
		return
	}
	if err := h.Store.DeleteRule(c.Request.Context(), id); err != nil {
		if err.Error() == "rule not found" {
			middleware.AbortWithError(c, apperrors.NotFound("rule not found"))
			return
		}
		middleware.AbortWithError(c, apperrors.Internal("failed to delete rule"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
