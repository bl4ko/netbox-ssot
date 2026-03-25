package client

import (
	"testing"
	"time"
)

func TestExponentialBackoff(t *testing.T) {
	tests := []struct {
		name    string
		attempt int
		want    time.Duration
	}{
		{
			name:    "attempt 0 returns initial backoff",
			attempt: 0,
			want:    initialBackoff,
		},
		{
			name:    "attempt 1 doubles",
			attempt: 1,
			want:    time.Duration(float64(initialBackoff) * backoffFactor),
		},
		{
			name:    "attempt 2 quadruples",
			attempt: 2,
			want:    time.Duration(float64(initialBackoff) * backoffFactor * backoffFactor),
		},
		{
			name:    "large attempt capped at max",
			attempt: 100,
			want:    maxBackoff,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := exponentialBackoff(tt.attempt)
			if got != tt.want {
				t.Errorf("exponentialBackoff(%d) = %v, want %v", tt.attempt, got, tt.want)
			}
		})
	}
}
