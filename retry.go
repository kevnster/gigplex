package gigplex

import "time"

// RetryPolicy decides whether to retry a failed job and how long to wait.
// Implement this interface to define custom retry behavior.
type RetryPolicy interface {
	ShouldRetry(attempt int, err error) bool
	Backoff(attempt int) time.Duration
}

// NoRetry never retries a failed job.
var NoRetry RetryPolicy = noRetry{}

type noRetry struct{}

func (noRetry) ShouldRetry(_ int, _ error) bool { return false }
func (noRetry) Backoff(_ int) time.Duration      { return 0 }

// DefaultRetry retries up to 3 times with exponential backoff starting at 2s.
var DefaultRetry RetryPolicy = exponentialBackoff{maxAttempts: 3, base: 2 * time.Second}

type exponentialBackoff struct {
	maxAttempts int
	base        time.Duration
}

func (e exponentialBackoff) ShouldRetry(attempt int, _ error) bool {
	return attempt < e.maxAttempts
}

func (e exponentialBackoff) Backoff(attempt int) time.Duration {
	d := e.base
	for i := 0; i < attempt; i++ {
		d *= 2
	}
	return d
}