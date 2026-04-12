package models

// AuditSession is the session handle returned by
// POST /api/query-editor/sessions.
//
// The session expires unless it is kept alive with periodic PATCH calls
// (see workflows/audit for the polling helper, not yet implemented).
type AuditSession struct {
	ID   string           `json:"id,omitempty"`
	Data AuditSessionData `json:"data,omitempty"`

	// Engine is the requested query engine ("sast" or "iac"). Filled by the
	// SDK from the create request; not part of the wire response.
	Engine string `json:"-"`

	// ProjectID / ApplicationID are the scopes the session was opened
	// against, if any. Filled by the SDK.
	ProjectID     string `json:"-"`
	ApplicationID string `json:"-"`

	// Languages is populated from the RequestID polling for SAST sessions.
	Languages []string `json:"-"`

	// Platforms is populated inline for IAC sessions (Data.QueryFilters).
	Platforms []string `json:"-"`
}

// AuditSessionData is the nested status envelope on an audit session response.
type AuditSessionData struct {
	Status       string   `json:"status,omitempty"`
	RequestID    string   `json:"requestId,omitempty"`
	QueryFilters []string `json:"queryFilters,omitempty"`
}

// AuditCreateRequest is the POST body for /api/query-editor/sessions.
//
// Either Filter is set (for a tenant-level session — "go", "java", "all", etc.)
// or ProjectID+ScanID are set (for a project/scan session). Scanner is
// required in both cases.
type AuditCreateRequest struct {
	Scanner   string `json:"scanner"`
	Filter    string `json:"filter,omitempty"`
	ProjectID string `json:"projectId,omitempty"`
	ScanID    string `json:"scanId,omitempty"`
}

// AuditRequestStatus is the response from
// GET /api/query-editor/sessions/{id}/requests/{requestId}. Status is one of
// "Running", "Completed", "Failed". Completed tracks Status=="Completed" but
// is exposed as a convenience flag by the polling endpoint.
type AuditRequestStatus struct {
	Completed    bool   `json:"completed,omitempty"`
	Value        any    `json:"value,omitempty"`
	ErrorCode    int    `json:"code,omitempty"`
	ErrorMessage string `json:"message,omitempty"`
	Status       string `json:"status,omitempty"`
}

// AuditScanSourceFile is one source-file entry returned by
// GET /api/query-editor/sessions/{id}/sources.
type AuditScanSourceFile struct {
	Filename string `json:"filename,omitempty"`
	Language string `json:"language,omitempty"`
	Size     int64  `json:"size,omitempty"`
}
