package gmail

import (
	"context"
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/api/googleapi"
)

const (
	maxAttempts    = 8
	initialBackoff = 2 * time.Second
	maxBackoff     = 60 * time.Second
)

// IsRetryable reports whether err is a transient Gmail API error worth retrying.
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}
	var gerr *googleapi.Error
	if !errors.As(err, &gerr) {
		s := strings.ToLower(err.Error())
		return strings.Contains(s, "429") ||
			strings.Contains(s, "503") ||
			strings.Contains(s, "502")
	}
	if gerr.Code == http.StatusTooManyRequests {
		return true
	}
	if gerr.Code >= 500 && gerr.Code < 600 {
		return true
	}
	if gerr.Code == http.StatusForbidden {
		for _, e := range gerr.Errors {
			switch strings.ToLower(e.Reason) {
			case "ratelimitexceeded", "userratelimitexceeded":
				return true
			}
		}
	}
	return false
}

// IsRateLimitError is true when retries are exhausted for a quota/rate error.
func IsRateLimitError(err error) bool {
	return IsRetryable(err)
}

func retryAfter(err error) time.Duration {
	var gerr *googleapi.Error
	if !errors.As(err, &gerr) || gerr.Header == nil {
		return 0
	}
	raw := strings.TrimSpace(gerr.Header.Get("Retry-After"))
	if raw == "" {
		return 0
	}
	if sec, convErr := strconv.Atoi(raw); convErr == nil && sec > 0 {
		return time.Duration(sec) * time.Second
	}
	if t, convErr := http.ParseTime(raw); convErr == nil {
		d := time.Until(t)
		if d > 0 {
			return d
		}
	}
	return 0
}

func backoffWithJitter(attempt int) time.Duration {
	delay := initialBackoff
	for i := 0; i < attempt; i++ {
		delay *= 2
		if delay > maxBackoff {
			delay = maxBackoff
			break
		}
	}
	if delay > maxBackoff {
		delay = maxBackoff
	}
	// jitter in [0.75, 1.25]
	jitter := 0.75 + rand.Float64()*0.5
	return time.Duration(float64(delay) * jitter)
}

func withRetry(ctx context.Context, log zerolog.Logger, fn func() error) error {
	var lastErr error
	for attempt := 0; attempt < maxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		lastErr = fn()
		if lastErr == nil {
			return nil
		}
		if !IsRetryable(lastErr) {
			return lastErr
		}
		if attempt == maxAttempts-1 {
			break
		}
		wait := retryAfter(lastErr)
		if wait == 0 {
			wait = backoffWithJitter(attempt)
		}
		log.Warn().
			Err(lastErr).
			Int("attempt", attempt+1).
			Int("max_attempts", maxAttempts).
			Dur("wait", wait).
			Msg("gmail retry")
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(wait):
		}
	}
	return lastErr
}
