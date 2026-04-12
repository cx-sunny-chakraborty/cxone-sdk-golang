package models

// SastMetadataModel is the response of GET /api/sast-metadata.
type SastMetadataModel struct {
	TotalCount int                `json:"totalCount"`
	Scans      []SastScanMetadata `json:"scans"`
	Missing    []string           `json:"missing"`
}

// SastScanMetadata holds metadata about one SAST scan: line counts,
// incremental status, query preset, etc.
//
// Named SastScanMetadata (not just "Scans") to avoid collision with the
// list-of-scans response shape in [ScansCollectionResponseModel].
type SastScanMetadata struct {
	ScanID                  string `json:"scanId,omitempty"`
	ProjectID               string `json:"projectId,omitempty"`
	Loc                     int    `json:"loc,omitempty"`
	FileCount               int    `json:"fileCount,omitempty"`
	IsIncremental           bool   `json:"isIncremental,omitempty"`
	IsIncrementalCanceled   bool   `json:"isIncrementalCanceled,omitempty"`
	IncrementalCancelReason string `json:"incrementalCancelReason,omitempty"`
	BaseID                  string `json:"baseId,omitempty"`
	AddedFilesCount         int    `json:"addedFilesCount,omitempty"`
	ChangedFilesCount       int    `json:"changedFilesCount,omitempty"`
	DeletedFilesCount       int    `json:"deletedFilesCount,omitempty"`
	QueryPreset             string `json:"queryPreset,omitempty"`
}
