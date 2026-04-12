// Package transport implements the SDK's request funnel — the single
// ExecuteRequest path that every HTTP call (high-level workflows and the
// internal endpoint layer alike) flows through. It applies authentication,
// retries, observability, the mandatory Checkmarx One headers (CLAUDE.md
// §4.2), and the user-supplied interceptor pipeline.
//
// This package is reachable from both the root cxone package and the
// internal/endpoints subtree thanks to Go's internal-package access rules,
// so endpoint functions can call Executor.Do directly without leaking the
// transport types onto the public surface. The root cxone package exposes
// the user-facing types ([Interceptor], [Request], [Response]) as type
// aliases.
package transport

import (
	"context"
	"io"
	"net/http"
	"time"
)

// Request describes one outgoing HTTP request as the SDK sees it. It is
// passed through the interceptor chain before [Executor] hands it to the
// underlying http.Client.
//
// The fields are intentionally minimal: callers build a Request, set the
// fields they care about, and let the executor merge in the mandatory
// headers, authentication, correlation id, and retry policy.
type Request struct {
	// Method is the HTTP method ("GET", "POST", …). Required.
	Method string

	// URL is the absolute, fully-qualified URL of the endpoint. The executor
	// does not perform base-URL resolution; callers (typically the endpoint
	// layer) compose the URL from the client's APIBaseURL plus the
	// endpoint-specific path.
	URL string

	// Header carries per-request headers. The executor will overlay the
	// mandatory cxone headers (Authorization, Accept, User-Agent,
	// CorrelationId) on top of these — duplicates are replaced.
	Header http.Header

	// Body is the request body, or nil for GET-style requests. The executor
	// reads it lazily and is responsible for closing it. If the body is a
	// *bytes.Buffer, *bytes.Reader, or *strings.Reader the executor will
	// re-seek it on retry; other readers will only be sent on the first
	// attempt and the request will not be retried after a transport failure.
	Body io.Reader

	// ContentType, if non-empty, is set as the Content-Type header before
	// the request is sent. Provided as a separate field so callers don't
	// need to allocate a Header map for the common case of a typed body.
	ContentType string

	// Timeout, if non-zero, overrides the client default timeout for this
	// one call (CLAUDE.md §4.6). Implemented by deriving a context.WithTimeout.
	Timeout time.Duration
}

// Response is the executor's view of a successful HTTP response. The body
// has NOT been read yet — callers are responsible for draining and closing
// it (typically via a small helper that decodes JSON and closes in one shot).
type Response struct {
	// StatusCode is the HTTP status code.
	StatusCode int

	// Header is the response header.
	Header http.Header

	// Body is the response body. Callers must Close it.
	Body io.ReadCloser

	// Attempt is the 1-indexed attempt number on which this response was
	// received (1 for first-try success, >1 if any retries occurred).
	Attempt int

	// CorrelationID echoes the per-client correlation id sent on the request.
	// Useful for surfacing in user-facing errors.
	CorrelationID string
}

// Interceptor is the user-extensible hook used to inspect and mutate
// requests/responses before they reach the network or the SDK caller.
//
// Implementations are registered on the client builder and invoked in
// registration order on the request side and reverse order on the response
// side (CLAUDE.md §11.4). Interceptors run inside the retry loop — Before
// runs once per attempt with the attempt number on the Request — so they
// can implement attempt-aware behavior such as logging or circuit breaking.
//
// Implementations must be safe for concurrent use.
type Interceptor interface {
	// Before runs prior to sending the HTTP request. Implementations may
	// mutate req.Header in place. Returning a non-nil error aborts the
	// request and surfaces the error to the caller without retrying.
	Before(ctx context.Context, req *Request, attempt int) error

	// After runs after a response is received (or, on transport error,
	// with resp == nil and err != nil). Implementations may inspect or
	// wrap err. Returning a non-nil error overrides the original error.
	// After may NOT consume resp.Body — that is the caller's job.
	After(ctx context.Context, req *Request, resp *Response, err error) error
}

// InterceptorFunc adapts a pair of plain functions to the Interceptor
// interface. Either function may be nil, in which case it is treated as a
// no-op for that side of the call.
type InterceptorFunc struct {
	BeforeFn func(ctx context.Context, req *Request, attempt int) error
	AfterFn  func(ctx context.Context, req *Request, resp *Response, err error) error
}

// Before implements [Interceptor].
func (f InterceptorFunc) Before(ctx context.Context, req *Request, attempt int) error {
	if f.BeforeFn == nil {
		return nil
	}
	return f.BeforeFn(ctx, req, attempt)
}

// After implements [Interceptor].
func (f InterceptorFunc) After(ctx context.Context, req *Request, resp *Response, err error) error {
	if f.AfterFn == nil {
		return err
	}
	return f.AfterFn(ctx, req, resp, err)
}
