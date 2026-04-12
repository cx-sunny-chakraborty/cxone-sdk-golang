package models

// ReportData is the inner payload shared by JSON and PDF report requests.
type ReportData struct {
	ScanID     string   `json:"scanId"`
	ProjectID  string   `json:"projectId"`
	BranchName string   `json:"branchName"`
	Host       string   `json:"host,omitempty"`
	Sections   []string `json:"sections,omitempty"`
	Scanners   []string `json:"scanners,omitempty"`
	Email      []string `json:"email,omitempty"`
}

// JSONReportRequest is the body of POST /api/reports/json.
type JSONReportRequest struct {
	Data       ReportData `json:"data"`
	ReportType string     `json:"reportType"`
	ReportName string     `json:"reportName"`
	FileFormat string     `json:"fileFormat"`
}

// JSONReportResponse is the immediate (202) response of POST /api/reports/json.
type JSONReportResponse struct {
	ReportID string `json:"reportId"`
}

// JSONReportPollingResponse is the response of GET /api/reports/json/{id}.
//
// Status is one of: "requested", "started", "completed", "failed".
type JSONReportPollingResponse struct {
	URL      string `json:"url,omitempty"`
	ReportID string `json:"reportId"`
	Status   string `json:"status"`
}

// PDFReportRequest is the body of POST /api/reports/pdf.
type PDFReportRequest struct {
	ReportName string     `json:"reportName"`
	ReportType string     `json:"reportType"`
	FileFormat string     `json:"fileFormat"`
	Data       ReportData `json:"data"`
}

// PDFReportResponse is the immediate (202) response of POST /api/reports/pdf.
type PDFReportResponse struct {
	ReportID string `json:"reportId"`
}

// PDFReportPollingResponse is the response of GET /api/reports/pdf/{id}.
type PDFReportPollingResponse struct {
	ReportID string `json:"reportId"`
	Status   string `json:"status"`
	URL      string `json:"url,omitempty"`
}

// Report status values returned in the polling responses.
const (
	ReportStatusRequested = "requested"
	ReportStatusStarted   = "started"
	ReportStatusCompleted = "completed"
	ReportStatusFailed    = "failed"
)
