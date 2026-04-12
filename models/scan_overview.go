package models

// SCSOverview is the response of GET /api/micro-engines/read/scans/{scanId}/scan-overview.
//
// Source-Code Security (SCS) is the umbrella for Checkmarx One's
// commit-history, secrets, scorecard, and 2ms scanners. The overview returns
// a per-engine status + risk summary.
type SCSOverview struct {
	Status               ScanStatus             `json:"status"`
	TotalRisksCount      int                    `json:"totalRisks,omitempty"`
	RiskSummary          map[string]int         `json:"riskSummary,omitempty"`
	MicroEngineOverviews []*MicroEngineOverview `json:"engineOverviews,omitempty"`
}

// MicroEngineOverview is the per-micro-engine overview entry inside [SCSOverview].
type MicroEngineOverview struct {
	Name        string         `json:"name"`
	FullName    string         `json:"fullName"`
	Status      ScanStatus     `json:"status"`
	TotalRisks  int            `json:"totalRisks"`
	RiskSummary map[string]any `json:"riskSummary"`
}
