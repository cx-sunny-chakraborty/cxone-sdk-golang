package migration

import (
	"context"
	"fmt"
	"time"
)

// WaitOptions configures [WaitUntilComplete].
type WaitOptions struct {
	// PollInterval is the time between successive status polls.
	// Default 15 seconds.
	PollInterval time.Duration

	// MaxDuration is an upper bound on total wait time. Zero means no
	// limit (cancellation is still honored via ctx).
	MaxDuration time.Duration
}

// WaitUntilComplete polls the import until it reaches a terminal state
// (completed, partial, failed, or blank). Returns the final inspector.
//
// Cancellation is honored: if ctx is cancelled, the waiter returns the
// last-observed inspector and ctx.Err().
func WaitUntilComplete(ctx context.Context, insp *ImportInspector, opts WaitOptions) (*ImportInspector, error) {
	if insp == nil {
		return nil, fmt.Errorf("migration: inspector is nil")
	}
	if opts.PollInterval <= 0 {
		opts.PollInterval = 15 * time.Second
	}
	deadline, hasDeadline := computeDeadline(ctx, opts.MaxDuration)

	for {
		if insp.Completed() {
			return insp, nil
		}
		if insp.Failed() {
			return insp, fmt.Errorf("import %s %s", insp.ID(), insp.Status())
		}
		if err := sleepUntil(ctx, opts.PollInterval, deadline, hasDeadline); err != nil {
			return insp, err
		}
		if _, err := insp.Refresh(ctx); err != nil {
			return insp, err
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
