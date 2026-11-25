package classify

import (
	"context"

	"github.com/dharundds/jobtrack/internal/store"
)

// RecomputeApplicationStatusFromEvents updates pipeline status from timeline unless manually overridden.
func RecomputeApplicationStatusFromEvents(ctx context.Context, st *store.Store, applicationID int64) error {
	app, err := st.GetApplication(ctx, applicationID)
	if err != nil {
		return err
	}
	if app.StatusOverridden {
		return nil
	}
	events, err := st.ListApplicationEvents(ctx, applicationID)
	if err != nil {
		return err
	}
	recs := make([]EventRecord, len(events))
	for i, e := range events {
		recs[i] = EventRecord{Type: e.Type, OccurredAt: e.OccurredAt}
	}
	status := LatestStatusFromEvents(recs)
	return st.UpdateApplicationStatusOnly(ctx, applicationID, status)
}