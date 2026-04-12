// Package observability defines the SDK's pluggable Logger, Metrics, and
// Tracer interfaces along with no-op default implementations.
//
// Wire your own implementations into the client via the builder:
//
//	client, _ := cxone.NewClient().
//	    Tenant("acme").
//	    Region(cxone.RegionUS).
//	    AgentName("MyApp").
//	    APIKey(apiKey).
//	    Logger(observability.NewSlogLogger(slog.Default())).
//	    Build()
//
// Defaults are no-ops; the SDK is silent unless you opt in. Structured event
// names match CLAUDE.md §11.1 across every Checkmarx One SDK.
package observability

import (
	"context"
	"log/slog"
)

// Standard structured event names. Use these as the message argument to
// Logger.Info / Logger.Error so that downstream log processors can filter
// SDK events consistently across languages.
const (
	EventRequestStart = "cxone.request.start"
	EventRequestEnd   = "cxone.request.end"
	EventRequestRetry = "cxone.request.retry"
	EventAuthRefresh  = "cxone.auth.refresh"
	EventScanPoll     = "cxone.scan.poll"
	EventReportPoll   = "cxone.report.poll"
	EventError        = "cxone.error"
)

// Field is a structured key/value pair attached to a log event.
//
// We deliberately do NOT use slog.Attr here so that the Logger interface
// remains usable from logger backends that pre-date slog (zap, zerolog,
// logrus). Adapters convert as needed.
type Field struct {
	Key   string
	Value any
}

// F is a convenience constructor for Field.
func F(key string, value any) Field { return Field{Key: key, Value: value} }

// Logger is the SDK's structured-logging contract. Every public method on the
// SDK that does I/O emits one of the events declared above; implementations
// route those events wherever they wish.
//
// Implementations must be safe for concurrent use.
type Logger interface {
	// Debug logs verbose diagnostic detail (e.g. request bodies, retry math).
	Debug(ctx context.Context, msg string, fields ...Field)
	// Info logs routine SDK lifecycle events (request start/end, auth refresh).
	Info(ctx context.Context, msg string, fields ...Field)
	// Warn logs recoverable problems (a retried request, a degraded response).
	Warn(ctx context.Context, msg string, fields ...Field)
	// Error logs unrecoverable failures. Bearer tokens and form-body secrets
	// must already have been redacted from msg and field values before they
	// reach the logger — the SDK's funnel takes care of this for built-in
	// events; user code that calls Logger directly is responsible for its own
	// inputs.
	Error(ctx context.Context, msg string, fields ...Field)
}

// NopLogger is a Logger that discards every event. It is the default when no
// logger is supplied to the client builder.
type NopLogger struct{}

func (NopLogger) Debug(context.Context, string, ...Field) {}
func (NopLogger) Info(context.Context, string, ...Field)  {}
func (NopLogger) Warn(context.Context, string, ...Field)  {}
func (NopLogger) Error(context.Context, string, ...Field) {}

// slogLogger adapts a *slog.Logger to the Logger interface.
type slogLogger struct {
	l *slog.Logger
}

// NewSlogLogger wraps a *slog.Logger so it can be passed to the client builder.
// If l is nil, slog.Default() is used.
func NewSlogLogger(l *slog.Logger) Logger {
	if l == nil {
		l = slog.Default()
	}
	return &slogLogger{l: l}
}

func (s *slogLogger) log(ctx context.Context, level slog.Level, msg string, fields []Field) {
	if !s.l.Enabled(ctx, level) {
		return
	}
	attrs := make([]any, 0, len(fields)*2)
	for _, f := range fields {
		attrs = append(attrs, f.Key, f.Value)
	}
	s.l.Log(ctx, level, msg, attrs...)
}

func (s *slogLogger) Debug(ctx context.Context, msg string, fields ...Field) {
	s.log(ctx, slog.LevelDebug, msg, fields)
}
func (s *slogLogger) Info(ctx context.Context, msg string, fields ...Field) {
	s.log(ctx, slog.LevelInfo, msg, fields)
}
func (s *slogLogger) Warn(ctx context.Context, msg string, fields ...Field) {
	s.log(ctx, slog.LevelWarn, msg, fields)
}
func (s *slogLogger) Error(ctx context.Context, msg string, fields ...Field) {
	s.log(ctx, slog.LevelError, msg, fields)
}
