package audit

import (
	"context"
	"fmt"
	"time"
)

// WaitOptions configures [WaitForRequest].
type WaitOptions struct {
	// PollInterval is the time between successive request-status polls.
	// Default 5 seconds. A keep-alive heartbeat is sent every other poll.
	PollInterval time.Duration

	// MaxDuration is an upper bound on total wait time. Zero means no
	// limit (cancellation is still honored via ctx).
	MaxDuration time.Duration
}

// WaitForRequest polls the asynchronous request identified by requestID
// until it completes (or fails), sending periodic keep-alives so the
// session doesn't expire. Returns the final Value from the request status
// response.
//
// Cancellation is honored: if ctx is cancelled, the waiter returns ctx.Err().
func WaitForRequest(ctx context.Context, insp *SessionInspector, requestID string, opts WaitOptions) (any, error) {
	if insp == nil {
		return nil, fmt.Errorf("audit: inspector is nil")
	}
	if requestID == "" {
		return nil, fmt.Errorf("audit: requestID is required")
	}
	if opts.PollInterval <= 0 {
		opts.PollInterval = 5 * time.Second
	}
	deadline, hasDeadline := computeDeadline(ctx, opts.MaxDuration)

	keepAliveCounter := 0
	for {
		status, err := insp.PollRequest(ctx, requestID)
		if err != nil {
			return nil, err
		}
		if status.ErrorCode != 0 && status.ErrorMessage != "" {
			return status.Value, fmt.Errorf("audit request %s: error %d: %s", requestID, status.ErrorCode, status.ErrorMessage)
		}
		if status.Status == "Failed" {
			return status.Value, fmt.Errorf("audit request %s failed: %v", requestID, status.Value)
		}
		if status.Completed {
			return status.Value, nil
		}

		if err := sleepUntil(ctx, opts.PollInterval, deadline, hasDeadline); err != nil {
			return nil, err
		}

		keepAliveCounter++
		if keepAliveCounter%2 == 0 {
			_ = insp.KeepAlive(ctx)
		}
	}
}

func computeDeadline(ctx context.Context, maxDur time.Duration) (time.Time, bool) {
	if maxDur > 0 {
		return time.Now().Add(maxDur), true
	}
	if d, ok := ctx.Deadline(); ok {
		return d, true
	}
	return time.Time{}, false
}

func sleepUntil(ctx context.Context, interval time.Duration, deadline time.Time, hasDeadline bool) error {
	if hasDeadline {
		if remaining := time.Until(deadline); remaining < interval {
			interval = remaining
			if interval <= 0 {
				return context.DeadlineExceeded
			}
		}
	}
	t := time.NewTimer(interval)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
