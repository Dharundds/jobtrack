package handlers

import (
	"net/http"
	"strconv"

	"github.com/dharundds/jobtrack/internal/api/apperrors"
	"github.com/dharundds/jobtrack/internal/api/middleware"
	"github.com/dharundds/jobtrack/internal/api/present"
	"github.com/dharundds/jobtrack/internal/store"
	"github.com/gin-gonic/gin"
)

type ApplicationsHandler struct {
	Store *store.Store
}

func (h *ApplicationsHandler) loadApplicationJSON(c *gin.Context, app store.ApplicationFull) (present.ApplicationJSON, error) {
	msgs, err := h.Store.ListMessagesForApplication(c.Request.Context(), app.ID)
	if err != nil {
		return present.ApplicationJSON{}, err
	}
	details, err := h.Store.ListEventDetailsForApplication(c.Request.Context(), app.ID)
	if err != nil {
		return present.ApplicationJSON{}, err
	}
	tid, err := h.Store.FirstThreadIDForApplication(c.Request.Context(), app.ID)
	if err != nil {
		return present.ApplicationJSON{}, err
	}
	return present.ApplicationFromStore(app, msgs, details, tid), nil
}

func (h *ApplicationsHandler) List(c *gin.Context) {
	status := c.Query("status")
	q := c.Query("q")
	awaiting := c.Query("awaiting_reply") == "true" || c.Query("awaiting_reply") == "1"

	apps, err := h.Store.ListApplicationsFull(c.Request.Context(), status, q)
	if err != nil {
		middleware.AbortWithError(c, apperrors.Internal("failed to list applications: "+err.Error()))
		return
	}

	out := make([]present.ApplicationJSON, 0, len(apps))
	for _, app := range apps {
		j, err := h.loadApplicationJSON(c, app)
		if err != nil {
			middleware.AbortWithError(c, apperrors.Internal("failed to load application"))
			return
		}
		if awaiting && !j.AwaitingMyReply {
			continue
		}
		out = append(out, j)
	}
	c.JSON(http.StatusOK, gin.H{"applications": out})
}

func (h *ApplicationsHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		middleware.AbortWithError(c, apperrors.BadRequest("invalid_id", "invalid application id"))
		return
	}
	app, err := h.Store.GetApplicationFull(c.Request.Context(), id)
	if err == store.ErrApplicationNotFound {
		middleware.AbortWithError(c, apperrors.NotFound("application not found"))
		return
	}
	if err != nil {
		middleware.AbortWithError(c, apperrors.Internal("failed to load application"))
		return
	}
	j, err := h.loadApplicationJSON(c, *app)
	if err != nil {
		middleware.AbortWithError(c, apperrors.Internal("failed to load application"))
		return
	}
	msgs, err := h.Store.ListMessagesForApplication(c.Request.Context(), id)
	if err != nil {
		middleware.AbortWithError(c, apperrors.Internal("failed to load messages"))
		return
	}
	msgOut := make([]present.MessageJSON, len(msgs))
	for i, m := range msgs {
		msgOut[i] = present.MessageFromStore(m, true)
	}
	events, err := h.Store.ListTimelineEvents(c.Request.Context(), id)
	if err != nil {
		middleware.AbortWithError(c, apperrors.Internal("failed to load events"))
		return
	}
	evOut := make([]present.EventJSON, len(events))
	for i, e := range events {
		evOut[i] = present.EventFromStore(e)
	}
	c.JSON(http.StatusOK, gin.H{
		"application": j,
		"messages":    msgOut,
		"events":      evOut,
	})
}

type patchApplicationRequest struct {
	Company *string `json:"company"`
	Role    *string `json:"role"`
	Status  *string `json:"status" binding:"omitempty,oneof=applied replied interview offer rejected"`
	Notes   *string `json:"notes"`
}

func (h *ApplicationsHandler) Patch(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		middleware.AbortWithError(c, apperrors.BadRequest("invalid_id", "invalid application id"))
		return
	}
	var req patchApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithError(c, apperrors.BadRequest("invalid_body", err.Error()))
		return
	}
	statusManual := req.Status != nil
	if err := h.Store.PatchApplication(c.Request.Context(), id, req.Company, req.Role, req.Status, req.Notes, statusManual); err != nil {
		if err == store.ErrApplicationNotFound {
			middleware.AbortWithError(c, apperrors.NotFound("application not found"))
			return
		}
		middleware.AbortWithError(c, apperrors.Internal("failed to update application"))
		return
	}
	app, err := h.Store.GetApplicationFull(c.Request.Context(), id)
	if err != nil {
		middleware.AbortWithError(c, apperrors.Internal("failed to load application"))
		return
	}
	j, err := h.loadApplicationJSON(c, *app)
	if err != nil {
		middleware.AbortWithError(c, apperrors.Internal("failed to load application"))
		return
	}
	c.JSON(http.StatusOK, j)
}

type mergeApplicationsRequest struct {
	TargetID  int64   `json:"target_id" binding:"required"`
	SourceIDs []int64 `json:"source_ids" binding:"required,min=1"`
}

func (h *ApplicationsHandler) Merge(c *gin.Context) {
	var req mergeApplicationsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AbortWithError(c, apperrors.BadRequest("invalid_body", err.Error()))
		return
	}
	if err := h.Store.MergeApplications(c.Request.Context(), req.TargetID, req.SourceIDs); err != nil {
		middleware.AbortWithError(c, apperrors.BadRequest("merge_failed", err.Error()))
		return
	}
	app, err := h.Store.GetApplicationFull(c.Request.Context(), req.TargetID)
	if err != nil {
		middleware.AbortWithError(c, apperrors.Internal("failed to load application"))
		return
	}
	j, err := h.loadApplicationJSON(c, *app)
	if err != nil {
		middleware.AbortWithError(c, apperrors.Internal("failed to load application"))
		return
	}
	c.JSON(http.StatusOK, j)
}
