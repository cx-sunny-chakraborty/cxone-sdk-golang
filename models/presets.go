package models

// Preset shapes — placeholder definitions.
//
// ast-cli does not wrap the presets endpoint (it passes preset names through
// to the scan engine as opaque strings), so the field shapes here are
// derived from common Checkmarx One REST conventions and the cross-language
// SDK contract in CLAUDE.md §9.5. The exact wire format MUST be verified
// against the live https://ast.checkmarx.net/spec/v1 swagger before
// shipping a stable release; treat these structs as best-effort until then.
//
// Marked TODO(spec) so a future engineer can grep for them.

// PresetDescriptor is the lightweight summary of one SAST preset.
//
// TODO(spec): verify field names against /api/presets schema.
type PresetDescriptor struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description,omitempty"`
	Custom         bool   `json:"custom,omitempty"`
	QueriesCount   int    `json:"queriesCount,omitempty"`
}

// PresetCollection is the response of GET /api/presets — a paginated list of
// preset descriptors.
//
// TODO(spec): verify envelope shape.
type PresetCollection struct {
	TotalCount         int                `json:"totalCount"`
	FilteredTotalCount int                `json:"filteredTotalCount"`
	Presets            []PresetDescriptor `json:"presets"`
}

// PresetDetail is the full preset including its query family + query
// references. Returned by GET /api/presets/{id}.
//
// TODO(spec): verify field names.
type PresetDetail struct {
	PresetDescriptor
	QueryIDs   []string `json:"queryIds,omitempty"`
	Tenant     string   `json:"tenant,omitempty"`
}

// QueryDescriptor is the lightweight summary of one SAST query.
type QueryDescriptor struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Group    string `json:"group,omitempty"`
	Language string `json:"language,omitempty"`
	Severity string `json:"severity,omitempty"`
	Custom   bool   `json:"custom,omitempty"`
}

// QueryCollection is the response of GET /api/queries.
//
// TODO(spec): verify envelope.
type QueryCollection struct {
	TotalCount         int               `json:"totalCount"`
	FilteredTotalCount int               `json:"filteredTotalCount"`
	Queries            []QueryDescriptor `json:"queries"`
}
