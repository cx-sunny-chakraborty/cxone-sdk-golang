package models

// SASTQuery is one SAST query entry as returned by the preset-manager or
// query-editor endpoints.
type SASTQuery struct {
	QueryID            uint64 `json:"queryID,string,omitempty"`
	Level              string `json:"level,omitempty"`
	LevelID            string `json:"levelId,omitempty"`
	Path               string `json:"path,omitempty"`
	Name               string `json:"queryName,omitempty"`
	Group              string `json:"group,omitempty"`
	Language           string `json:"language,omitempty"`
	Severity           string `json:"severity,omitempty"`
	CweID              int64  `json:"cweID,omitempty"`
	IsExecutable       bool   `json:"isExecutable,omitempty"`
	QueryDescriptionID int64  `json:"queryDescriptionId,omitempty"`
	Custom             bool   `json:"custom,omitempty"`
	EditorKey          string `json:"key,omitempty"`
	SastID             uint64 `json:"sastId,omitempty"`
	Source             string `json:"source,omitempty"`
}

// SASTQueryGroup groups SAST queries under a category/group name and language.
type SASTQueryGroup struct {
	Name     string      `json:"name,omitempty"`
	Language string      `json:"language,omitempty"`
	Queries  []SASTQuery `json:"queries,omitempty"`
}

// SASTQueryLanguage groups query groups under a language.
type SASTQueryLanguage struct {
	Name        string           `json:"name,omitempty"`
	QueryGroups []SASTQueryGroup `json:"queryGroups,omitempty"`
}

// SASTQueryCollection is a hierarchical view of SAST queries organized by
// language → group → query.
type SASTQueryCollection struct {
	QueryLanguages []SASTQueryLanguage `json:"queryLanguages,omitempty"`
}

// IACQuery is one IAC/KICS query entry.
type IACQuery struct {
	QueryID        string `json:"queryId,omitempty"`
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	DescriptionID  string `json:"descriptionId,omitempty"`
	DescriptionURL string `json:"descriptionUrl,omitempty"`
	Platform       string `json:"platform,omitempty"`
	Group          string `json:"group,omitempty"`
	Category       string `json:"category,omitempty"`
	Severity       string `json:"severity,omitempty"`
	CWE            string `json:"cwe,omitempty"`
	Level          string `json:"level,omitempty"`
	Path           string `json:"path,omitempty"`
	Source         string `json:"source,omitempty"`
	Custom         bool   `json:"custom,omitempty"`
	Key            string `json:"key,omitempty"`
}

// IACQueryGroup groups IAC queries under a group name and platform.
type IACQueryGroup struct {
	Name     string     `json:"name,omitempty"`
	Platform string     `json:"platform,omitempty"`
	Queries  []IACQuery `json:"queries,omitempty"`
}

// IACQueryPlatform groups query groups under a platform.
type IACQueryPlatform struct {
	Name        string          `json:"name,omitempty"`
	QueryGroups []IACQueryGroup `json:"queryGroups,omitempty"`
}

// IACQueryCollection is a hierarchical view of IAC queries organized by
// platform → group → query.
type IACQueryCollection struct {
	Platforms []IACQueryPlatform `json:"platforms,omitempty"`
}

// QueryFamily is one family entry from the preset-manager families endpoint.
type QueryFamily struct {
	Name       string   `json:"familyName,omitempty"`
	TotalCount uint64   `json:"totalCount,omitempty"`
	QueryIDs   []string `json:"queryIds,omitempty"`
}

// QueryFailure is one validation/execution error returned by the query
// editor.
type QueryFailure struct {
	QueryID string       `json:"query_id,omitempty"`
	Errors  []QueryError `json:"error,omitempty"`
}

// QueryError is one error entry inside a [QueryFailure].
type QueryError struct {
	Line    int    `json:"line,omitempty"`
	Column  int    `json:"column,omitempty"`
	Message string `json:"message,omitempty"`
}

// AuditQueryTree is one node in the hierarchical query tree returned by
// the query editor.
type AuditQueryTree struct {
	IsLeaf   bool             `json:"isLeaf,omitempty"`
	Title    string           `json:"title,omitempty"`
	Key      string           `json:"key,omitempty"`
	Data     AuditQueryData   `json:"data,omitempty"`
	Children []AuditQueryTree `json:"children,omitempty"`
}

// AuditQueryData carries metadata on a query tree node.
type AuditQueryData struct {
	Level    string `json:"level,omitempty"`
	Severity string `json:"severity,omitempty"`
	CWE      int64  `json:"cwe,omitempty"`
	Custom   bool   `json:"custom,omitempty"`
}

// PresetManagerEntry is the full CRUD representation of a SAST or IAC preset
// from the preset-manager endpoint. Distinguished from [PresetDetail] in
// presets.go which wraps the /api/presets descriptor used by PresetReader.
type PresetManagerEntry struct {
	ID                 string        `json:"id,omitempty"`
	Name               string        `json:"name,omitempty"`
	Description        string        `json:"description,omitempty"`
	Custom             bool          `json:"custom,omitempty"`
	AssociatedProjects uint64        `json:"associatedProjects,omitempty"`
	IsTenantDefault    bool          `json:"isTenantDefault,omitempty"`
	IsMigrated         bool          `json:"isMigrated,omitempty"`
	QueryFamilies      []QueryFamily `json:"queries,omitempty"`
}

// PresetListResponse wraps the preset-manager list endpoint envelope.
type PresetListResponse struct {
	Presets    []PresetManagerEntry `json:"presets,omitempty"`
	TotalCount int                 `json:"totalCount,omitempty"`
}

// PresetCreateRequest is the body of POST /preset-manager/{engine}/presets.
type PresetCreateRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	QueryIDs    []string `json:"queryIds,omitempty"`
}

// SASTAggregateSummary is one row from the SAST scan-summary aggregate.
type SASTAggregateSummary struct {
	Status    string `json:"status,omitempty"`
	QueryID   uint64 `json:"queryID,string,omitempty"`
	QueryName string `json:"queryName,omitempty"`
	Severity  string `json:"severity,omitempty"`
	Language  string `json:"language,omitempty"`
	Count     uint64 `json:"count,omitempty"`
}

// ScanSummary is the status summary returned by GET /api/scans/summary.
type ScanSummary struct {
	ScanID     string `json:"scanId,omitempty"`
	Status     string `json:"status,omitempty"`
	Engine     string `json:"engine,omitempty"`
	TotalCount int    `json:"totalCount,omitempty"`
}

// AMRole is an Access Management Phase 2 role, distinct from Keycloak roles.
type AMRole struct {
	ID          string `json:"id,omitempty"`
	TenantID    string `json:"tenantId,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	SystemRole  bool   `json:"systemRole,omitempty"`
}

// AMPermission is an Access Management permission entry.
type AMPermission struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	ChildIDs    []string `json:"childIds,omitempty"`
	ParentIDs   []string `json:"parentIds,omitempty"`
	Category    string   `json:"category,omitempty"`
}
