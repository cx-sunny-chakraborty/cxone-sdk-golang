package observability

import "context"

// Metrics is the SDK's pluggable metrics contract. The funnel emits the four
// instruments declared in CLAUDE.md §11.2:
//
//   - Request count, dimensioned by HTTP method + status.
//   - Request latency, dimensioned by HTTP method.
//   - Retry count, dimensioned by reason ("status_5xx", "status_502",
//     "transport", "auth_401", …).
//   - Auth refresh count.
//
// Implementations must be safe for concurrent use. The default ([NopMetrics])
// discards every measurement.
type Metrics interface {
	// IncRequest increments the request counter for a completed request.
	// Status is the HTTP status code (0 if no response was received).
	IncRequest(ctx context.Context, method string, status int)

	// ObserveRequestLatency records request latency in milliseconds.
	ObserveRequestLatency(ctx context.Context, method string, milliseconds float64)

	// IncRetry increments the retry counter. Reason is a short, low-cardinality
	// label suitable for use as a metric dimension.
	IncRetry(ctx context.Context, reason string)

	// IncAuthRefresh increments the auth-refresh counter. Reason is "initial"
	// for the lazy first auth and "401" for refreshes triggered by an
	// unauthorized response.
	IncAuthRefresh(ctx context.Context, reason string)
}

// NopMetrics is the default Metrics implementation. It discards every event.
type NopMetrics struct{}

func (NopMetrics) IncRequest(context.Context, string, int)              {}
func (NopMetrics) ObserveRequestLatency(context.Context, string, float64) {}
func (NopMetrics) IncRetry(context.Context, string)                     {}
func (NopMetrics) IncAuthRefresh(context.Context, string)               {}
