package observability

import "context"

// Tracer is the SDK's tracing contract. It is intentionally minimal and free
// of any third-party dependency so the SDK does not pull in OpenTelemetry by
// default — bring your own adapter and the OTel import lives in your code,
// not in this module.
//
// Span attributes set by the SDK's funnel always include:
//
//   - http.method            (string)
//   - http.url               (string, sanitized)
//   - http.status_code       (int, when a response was received)
//   - cxone.correlation_id   (string, the per-client UUID from §4.2)
//   - cxone.attempt          (int, 1-indexed)
//
// CLAUDE.md §11.3 calls for OpenTelemetry-compatible spans; the [Span]
// interface's method shapes are deliberately compatible with otel/trace.Span
// so a thin adapter is trivial.
type Tracer interface {
	// Start opens a span with the given name. The returned context carries
	// the new span; calling [Span.End] on the returned span closes it.
	Start(ctx context.Context, name string) (context.Context, Span)
}

// Span represents an in-progress unit of work. Implementations must be safe
// for concurrent use.
type Span interface {
	// SetAttribute attaches a typed key/value pair to the span.
	SetAttribute(key string, value any)

	// SetStatus marks the span as failed (when err != nil) or successful.
	SetStatus(err error)

	// End closes the span and records its duration.
	End()
}

// NopTracer is the default Tracer implementation. It returns a no-op span
// that satisfies the interface without recording anything.
type NopTracer struct{}

// nopSpan is the [Span] returned by [NopTracer].
type nopSpan struct{}

// Start implements [Tracer].
func (NopTracer) Start(ctx context.Context, _ string) (context.Context, Span) {
	return ctx, nopSpan{}
}

func (nopSpan) SetAttribute(string, any) {}
func (nopSpan) SetStatus(error)          {}
func (nopSpan) End()                     {}
