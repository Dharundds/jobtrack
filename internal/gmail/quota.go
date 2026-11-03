package gmail

import (
	"context"

	"golang.org/x/time/rate"
)

// Gmail API quota unit costs (per method).
const (
	CostMessagesList    = 5
	CostMessagesGet     = 5
	CostHistoryList     = 2
	CostUsersGetProfile = 1
)

// QuotaLimiter throttles API calls to stay under a units-per-minute budget.
type QuotaLimiter struct {
	lim *rate.Limiter
}

func NewQuotaLimiter(unitsPerMinute int) *QuotaLimiter {
	if unitsPerMinute < 1 {
		unitsPerMinute = 2500
	}
	// tokens/sec; burst allows short bursts without blocking the first calls.
	perSec := rate.Limit(float64(unitsPerMinute) / 60.0)
	burst := unitsPerMinute / 10
	if burst < 5 {
		burst = 5
	}
	return &QuotaLimiter{lim: rate.NewLimiter(perSec, burst)}
}

func (q *QuotaLimiter) Wait(ctx context.Context, units int) error {
	if units < 1 {
		units = 1
	}
	return q.lim.WaitN(ctx, units)
}
