package models

// AITelemetryEvent is the body of POST /api/telemetry/log — used to report
// AI-assisted scan events back to the platform for product analytics.
type AITelemetryEvent struct {
	AIProvider      string `json:"aiProvider"`
	ProblemSeverity string `json:"problemSeverity"`
	Type            string `json:"type"`
	SubType         string `json:"subType"`
	Agent           string `json:"agent"`
	Engine          string `json:"engine"`
	ScanType        string `json:"scanType"`
	Status          string `json:"status"`
	TotalCount      int    `json:"totalCount"`
	UniqueID        string `json:"uniqueId"`
}
