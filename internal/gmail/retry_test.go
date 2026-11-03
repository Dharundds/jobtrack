package gmail

import (
	"fmt"
	"net/http"
	"testing"

	"google.golang.org/api/googleapi"
)

func TestIsRetryable(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil",
			err:  nil,
			want: false,
		},
		{
			name: "429",
			err:  &googleapi.Error{Code: http.StatusTooManyRequests},
			want: true,
		},
		{
			name: "503",
			err:  &googleapi.Error{Code: http.StatusServiceUnavailable},
			want: true,
		},
		{
			name: "403 rateLimitExceeded",
			err: &googleapi.Error{
				Code: http.StatusForbidden,
				Errors: []googleapi.ErrorItem{
					{Reason: "rateLimitExceeded", Message: "Quota exceeded"},
				},
			},
			want: true,
		},
		{
			name: "403 userRateLimitExceeded",
			err: &googleapi.Error{
				Code: http.StatusForbidden,
				Errors: []googleapi.ErrorItem{
					{Reason: "userRateLimitExceeded"},
				},
			},
			want: true,
		},
		{
			name: "403 insufficientPermissions",
			err: &googleapi.Error{
				Code: http.StatusForbidden,
				Errors: []googleapi.ErrorItem{
					{Reason: "insufficientPermissions"},
				},
			},
			want: false,
		},
		{
			name: "403 without reasons",
			err:  &googleapi.Error{Code: http.StatusForbidden},
			want: false,
		},
		{
			name: "404",
			err:  &googleapi.Error{Code: http.StatusNotFound},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsRetryable(tt.err); got != tt.want {
				t.Fatalf("IsRetryable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsRetryable_wrapped(t *testing.T) {
	inner := &googleapi.Error{
		Code: http.StatusForbidden,
		Errors: []googleapi.ErrorItem{{Reason: "rateLimitExceeded"}},
	}
	wrapped := fmt.Errorf("gmail: %w", inner)
	if !IsRetryable(wrapped) {
		t.Fatal("expected wrapped google error to be retryable")
	}
}
