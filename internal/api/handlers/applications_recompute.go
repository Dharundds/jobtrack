package handlers

import (
	"context"

	"github.com/dharundds/jobtrack/internal/classify"
)

func (h *ApplicationsHandler) recomputeApplicationStatus(ctx context.Context, applicationID int64) error {
	app, err := h.Store.GetApplication(ctx, applicationID)
	if err != nil {
		return err
	}
	if app.StatusOverridden {
		return nil
	}
	events, err := h.Store.ListApplicationEvents(ctx, applicationID)
	if err != nil {
		return err
	}
	recs := make([]classify.EventRecord, len(events))
	for i, e := range events {
		recs[i] = classify.EventRecord{Type: e.Type, OccurredAt: e.OccurredAt}
	}
	status := classify.LatestStatusFromEvents(recs)
	return h.Store.UpdateApplicationStatusOnly(ctx, applicationID, status)
}
