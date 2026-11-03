package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/dharundds/jobtrack/internal/api/apperrors"
	"github.com/dharundds/jobtrack/internal/api/middleware"
	"github.com/dharundds/jobtrack/internal/classify"
	"github.com/dharundds/jobtrack/internal/store"
	"github.com/gin-gonic/gin"
)

func (h *ApplicationsHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		middleware.AbortWithError(c, apperrors.BadRequest("invalid_id", "invalid application id"))
		return
	}
	if err := h.Store.DeleteApplication(c.Request.Context(), id); err != nil {
		if err == store.ErrApplicationNotFound {
			middleware.AbortWithError(c, apperrors.NotFound("application not found"))
			return
		}
		middleware.AbortWithError(c, apperrors.Internal("failed to delete application"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *ApplicationsHandler) IgnoreMessage(c *gin.Context) {
	appID, msgID, ok := parseAppMessageIDs(c)
	if !ok {
		return
	}
	ctx := c.Request.Context()
	if err := h.Store.IgnoreMessage(ctx, appID, msgID); err != nil {
		if err == store.ErrMessageNotFound {
			middleware.AbortWithError(c, apperrors.NotFound("message not found"))
			return
		}
		middleware.AbortWithError(c, apperrors.Internal("failed to remove message"))
		return
	}
	remaining, _ := h.Store.CountMessagesForApplication(ctx, appID)
	if remaining > 0 {
		_ = h.recomputeApplicationStatus(ctx, appID)
	}
	c.JSON(http.StatusOK, gin.H{"ok": true, "application_deleted": remaining == 0})
}

type splitMessageRequest struct {
	Company *string `json:"company"`
	Role    *string `json:"role"`
	Status  *string `json:"status" binding:"omitempty,oneof=applied replied interview offer rejected"`
}

func (h *ApplicationsHandler) SplitMessage(c *gin.Context) {
	appID, msgID, ok := parseAppMessageIDs(c)
	if !ok {
		return
	}
	var req splitMessageRequest
	_ = c.ShouldBindJSON(&req)

	ctx := c.Request.Context()
	msg, err := h.Store.GetMessage(ctx, msgID)
	if err != nil {
		if err == store.ErrMessageNotFound {
			middleware.AbortWithError(c, apperrors.NotFound("message not found"))
			return
		}
		middleware.AbortWithError(c, apperrors.Internal("failed to load message"))
		return
	}
	if !msg.ApplicationID.Valid || msg.ApplicationID.Int64 != appID {
		middleware.AbortWithError(c, apperrors.NotFound("message not found"))
		return
	}

	company := derefString(req.Company)
	role := derefString(req.Role)
	if company == "" {
		company = classify.ExtractCompany(msg.FromAddr, msg.Subject, msg.BodyExcerpt)
	}
	if role == "" {
		role = classify.ExtractRole(msg.Subject, msg.BodyExcerpt)
	}
	if company == "" {
		middleware.AbortWithError(c, apperrors.BadRequest("company_required", "company is required to split"))
		return
	}

	status := classify.StatusFromEvent(msg.EventType)
	if req.Status != nil {
		status = *req.Status
	}

	newID, err := h.Store.SplitMessageToApplication(ctx, appID, msgID, company, role, status, appliedAtForSplit(msg))
	if err != nil {
		if err == store.ErrMessageNotFound {
			middleware.AbortWithError(c, apperrors.NotFound("message not found"))
			return
		}
		middleware.AbortWithError(c, apperrors.BadRequest("split_failed", err.Error()))
		return
	}

	if n, _ := h.Store.CountMessagesForApplication(ctx, appID); n > 0 {
		_ = h.recomputeApplicationStatus(ctx, appID)
	}
	_ = h.recomputeApplicationStatus(ctx, newID)

	app, err := h.Store.GetApplicationFull(ctx, newID)
	if err != nil {
		middleware.AbortWithError(c, apperrors.Internal("failed to load new application"))
		return
	}
	j, err := h.loadApplicationJSON(c, *app)
	if err != nil {
		middleware.AbortWithError(c, apperrors.Internal("failed to load new application"))
		return
	}

	sourceDeleted := false
	if n, _ := h.Store.CountMessagesForApplication(ctx, appID); n == 0 {
		sourceDeleted = true
	}

	c.JSON(http.StatusOK, gin.H{
		"application":         j,
		"application_id":      newID,
		"source_deleted":      sourceDeleted,
	})
}

func parseAppMessageIDs(c *gin.Context) (appID, msgID int64, ok bool) {
	appID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || appID <= 0 {
		middleware.AbortWithError(c, apperrors.BadRequest("invalid_id", "invalid application id"))
		return 0, 0, false
	}
	msgID, err = strconv.ParseInt(c.Param("mid"), 10, 64)
	if err != nil || msgID <= 0 {
		middleware.AbortWithError(c, apperrors.BadRequest("invalid_id", "invalid message id"))
		return 0, 0, false
	}
	return appID, msgID, true
}

func derefString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func appliedAtForSplit(msg *store.MessageRow) *time.Time {
	if msg.EventType != classify.EventApplied {
		return nil
	}
	t := msg.ReceivedAt
	return &t
}
