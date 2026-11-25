package handlers

import (
	"context"

	"github.com/dharundds/jobtrack/internal/classify"
)

func (h *ApplicationsHandler) recomputeApplicationStatus(ctx context.Context, applicationID int64) error {
	return classify.RecomputeApplicationStatusFromEvents(ctx, h.Store, applicationID)
}
