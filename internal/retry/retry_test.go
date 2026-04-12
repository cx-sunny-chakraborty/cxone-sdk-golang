package retry

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"testing"
	"time"
)

func TestClassifyStatus(t *testing.T) {
	cases := []struct {
		status int
		want   Outcome
	}{
		{200, OutcomeSuccess},
		{201, OutcomeSuccess},
		{299, OutcomeSuccess},
		{400, OutcomeFatal},
		{401, OutcomeFatal}, // funnel handles 401 separately
		{403, OutcomeFatal},
		{404, OutcomeFatal},
		{500, OutcomeRetry},
		{502, OutcomeRetry},
		{503, OutcomeRetry},
		{504, OutcomeRetry},
		{505, OutcomeFatal},
	}
	for _, tc := range cases {
		if got := ClassifyStatus(tc.status); got != tc.want {
			t.Errorf("ClassifyStatus(%d)=%v, want %v", tc.status, got, tc.want)
		}
	}
}

func TestIsTransientNetworkError(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{"context canceled", context.Canceled, false},
		{"deadline exceeded", context.DeadlineExceeded, false},
		{"random error", errors.New("nope"), false},
		{"connection reset msg", errors.New("read tcp: connection reset by peer"), true},
		{"net.OpError", &net.OpError{Op: "dial", Err: errors.New("oops")}, true},
		{"DNS error", &net.DNSError{Err: "lookup failed", Name: "x"}, true},
		{"url.Error wrapping op", &url.Error{Op: "Get", URL: "http://x", Err: &net.OpError{Op: "dial", Err: errors.New("x")}}, true},
		{"timeout msg", errors.New("Client.Timeout exceeded while awaiting headers"), true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsTransientNetworkError(tc.err); got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestDoSuccessFirstAttempt(t *testing.T) {
	calls := 0
	attempt, err := Do(context.Background(), Policy{MaxAttempts: 3}, func(int) (Outcome, error) {
		calls++
		return OutcomeSuccess, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if attempt != 1 || calls != 1 {
		t.Errorf("got attempt=%d calls=%d, want 1/1", attempt, calls)
	}
}

func TestDoRetriesUntilSuccess(t *testing.T) {
	p := Policy{MaxAttempts: 4, MaxDelay: 0}
	calls := 0
	attempt, err := Do(context.Background(), p, func(a int) (Outcome, error) {
		calls++
		if a < 3 {
			return OutcomeRetry, errors.New("transient")
		}
		return OutcomeSuccess, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if attempt != 3 || calls != 3 {
		t.Errorf("got attempt=%d calls=%d, want 3/3", attempt, calls)
	}
}

func TestDoExhaustion(t *testing.T) {
	p := Policy{MaxAttempts: 3, MaxDelay: 0}
	want := errors.New("boom")
	attempt, err := Do(context.Background(), p, func(int) (Outcome, error) {
		return OutcomeRetry, want
	})
	if !errors.Is(err, want) {
		t.Errorf("got err=%v, want %v", err, want)
	}
	if attempt != 3 {
		t.Errorf("got attempt=%d, want 3", attempt)
	}
}

func TestDoFatalNoRetry(t *testing.T) {
	want := errors.New("nope")
	calls := 0
	_, err := Do(context.Background(), Policy{MaxAttempts: 5}, func(int) (Outcome, error) {
		calls++
		return OutcomeFatal, want
	})
	if !errors.Is(err, want) || calls != 1 {
		t.Errorf("got err=%v calls=%d, want %v / 1", err, calls, want)
	}
}

func TestDoCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	calls := 0
	_, err := Do(ctx, Policy{MaxAttempts: 3}, func(int) (Outcome, error) {
		calls++
		return OutcomeRetry, errors.New("x")
	})
	if !errors.Is(err, context.Canceled) {
		t.Errorf("got err=%v, want context.Canceled", err)
	}
	if calls != 0 {
		t.Errorf("expected no attempts on pre-cancelled ctx, got %d", calls)
	}
}

func TestPolicyNextDelayJittered(t *testing.T) {
	// Deterministic source.
	p := Policy{MaxAttempts: 5, MaxDelay: 5 * time.Second, Randomize: true}.
		WithRNG(rand.New(rand.NewSource(42)))
	for i := 1; i < 5; i++ {
		d := p.nextDelay(i)
		if d < time.Second || d > 5*time.Second {
			t.Errorf("attempt %d: delay %v outside [1s,5s]", i, d)
		}
	}
}

func TestPolicyNextDelayFixedWhenNotRandomized(t *testing.T) {
	p := Policy{MaxAttempts: 3, MaxDelay: 4 * time.Second, Randomize: false}
	if got := p.nextDelay(1); got != 4*time.Second {
		t.Errorf("got %v, want 4s", got)
	}
}

func TestSleepRespectsContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	start := time.Now()
	err := Sleep(ctx, time.Hour)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("got %v, want DeadlineExceeded", err)
	}
	if time.Since(start) > 100*time.Millisecond {
		t.Errorf("sleep did not return promptly: %v", time.Since(start))
	}
}

// example for godoc
func ExampleDo() {
	calls := 0
	_, _ = Do(context.Background(), Policy{MaxAttempts: 2, MaxDelay: 0}, func(a int) (Outcome, error) {
		calls++
		if a == 1 {
			return OutcomeRetry, fmt.Errorf("nope")
		}
		return OutcomeSuccess, nil
	})
	fmt.Println(calls)
	// Output: 2
}
