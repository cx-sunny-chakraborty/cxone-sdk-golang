package models

// ExportRequestPayload is the body of POST /api/sca/export/requests.
//
// FileFormat is one of the platform-supported export formats: "ScanReportJson",
// "ScanReportXml", "CycloneDxJson", "CycloneDxXml", "SpdxJson", etc.
type ExportRequestPayload struct {
	ScanID           string           `json:"ScanID"`
	FileFormat       string           `json:"FileFormat"`
	ExportParameters ExportParameters `json:"ExportParameters,omitempty"`
}

// ExportParameters tweaks the contents of the exported document.
type ExportParameters struct {
	HideDevAndTestDependencies bool `json:"hideDevAndTestDependencies"`
	ShowOnlyEffectiveLicenses  bool `json:"showOnlyEffectiveLicenses"`
	ExcludePackages            bool `json:"excludePackages"`
	ExcludeLicenses            bool `json:"excludeLicenses"`
	ExcludeVulnerabilities     bool `json:"excludeVulnerabilities"`
	ExcludePolicies            bool `json:"excludePolicies"`
}

// ExportResponse is the immediate (202) response of POST /api/sca/export/requests.
type ExportResponse struct {
	ExportID string `json:"exportId"`
}

// ExportPollingResponse is the polling response of GET /api/sca/export/requests.
//
// ExportStatus is one of: "Requested", "InProgress", "Completed", "Failed".
type ExportPollingResponse struct {
	ExportID     string `json:"exportId"`
	ExportStatus string `json:"exportStatus"`
	FileURL      string `json:"fileUrl"`
	ErrorMessage string `json:"errorMessage"`
}
