package transport

// HTTP header names and values mandated by CLAUDE.md §4.2 for every API
// request. Constants are exported so endpoint packages and tests can match on
// them without re-typing the magic strings.
const (
	// HeaderAccept is the verbatim Accept header required by the Checkmarx
	// One API. Do not change.
	HeaderAccept = "*/*; version=1.0"

	// HeaderCorrelationID is the per-client correlation id header name. The
	// VALUE is generated once per client and reused for the lifetime of the
	// client (CLAUDE.md §4.2).
	HeaderCorrelationID = "CorrelationId"
)
