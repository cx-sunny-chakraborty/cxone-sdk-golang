package models

import (
	"encoding/json"
	"fmt"
	"time"
)

// AnalyticsTimeLayout is the format Checkmarx One analytics endpoints use
// for start/end dates.
const AnalyticsTimeLayout = "2006-01-02 15:04:05"

// AnalyticsTime wraps [time.Time] with the analytics-API-specific
// "2006-01-02 15:04:05" serialization.
type AnalyticsTime struct {
	time.Time
}

// UnmarshalJSON parses an analytics timestamp.
func (t *AnalyticsTime) UnmarshalJSON(b []byte) error {
	if len(b) < 2 {
		return fmt.Errorf("analytics time: empty value")
	}
	s := string(b[1 : len(b)-1])
	parsed, err := time.Parse(AnalyticsTimeLayout, s)
	if err != nil {
		return fmt.Errorf("analytics time: %w", err)
	}
	t.Time = parsed
	return nil
}

// MarshalJSON renders an analytics timestamp.
func (t AnalyticsTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Time.Format(AnalyticsTimeLayout))
}

// AnalyticsFilter is the common filter body accepted by every analytics KPI.
// Zero-value fields are omitted from the request.
type AnalyticsFilter struct {
	Projects        []string       `json:"projects,omitempty"`
	Applications    []string       `json:"applications,omitempty"`
	Scanners        []string       `json:"scanners,omitempty"`
	ApplicationTags []string       `json:"applicationTags,omitempty"`
	ProjectTags     []string       `json:"projectTags,omitempty"`
	ScanTags        []string       `json:"scanTags,omitempty"`
	States          []string       `json:"states,omitempty"`
	Status          []string       `json:"status,omitempty"`
	Severities      []string       `json:"severities,omitempty"`
	BranchNames     []string       `json:"branchNames,omitempty"`
	Timezone        string         `json:"timezone,omitempty"`
	Groups          []string       `json:"groupIds,omitempty"`
	StartDate       *AnalyticsTime `json:"startDate,omitempty"`
	EndDate         *AnalyticsTime `json:"endDate,omitempty"`
}

// AnalyticsRequest is the full POST body for /api/data_analytics/analyticsAPI/v1.
type AnalyticsRequest struct {
	AnalyticsFilter
	KPI    string  `json:"kpi"`
	Limit  uint64  `json:"limit,omitempty"`
	Offset *uint64 `json:"offset,omitempty"`
}

// AnalyticsDistributionEntry is one slice of a "by X total" KPI.
type AnalyticsDistributionEntry struct {
	Label      string  `json:"label"`
	Density    float32 `json:"density"`
	Percentage float32 `json:"percentage"`
	Results    uint64  `json:"results"`
}

// AnalyticsDistributionBlock groups distribution entries by a primary label.
type AnalyticsDistributionBlock struct {
	Label  string                       `json:"label"`
	Values []AnalyticsDistributionEntry `json:"values"`
}

// AnalyticsDistributionStats is the envelope returned by every
// "vulnerabilitiesByXTotal" KPI.
type AnalyticsDistributionStats struct {
	Distribution []AnalyticsDistributionBlock `json:"distribution"`
	LOC          uint64                       `json:"loc"`
	Total        uint64                       `json:"total"`
}

// AnalyticsLabeledResultsCount is one label + count pair used by
// severity/state/aging stats.
type AnalyticsLabeledResultsCount struct {
	Label   string `json:"label"`
	Results int64  `json:"results"`
}

// AnalyticsSeverityAndStateEntry is one row in the severity-and-state KPI.
type AnalyticsSeverityAndStateEntry struct {
	Label      string                         `json:"label"`
	Results    int64                          `json:"results"`
	Severities []AnalyticsLabeledResultsCount `json:"severities"`
}

// AnalyticsAgingEntry is one row in the aging KPI (same shape as severity/state).
type AnalyticsAgingEntry = AnalyticsSeverityAndStateEntry

// AnalyticsAgingResponse wraps the agingTotal KPI envelope.
type AnalyticsAgingResponse struct {
	AgingAndSeverities []AnalyticsAgingEntry `json:"agingAndSeverities"`
}

// AnalyticsOverTimeEntry is one point in an over-time series.
type AnalyticsOverTimeEntry struct {
	Time  uint64        `json:"time"`
	Value float32       `json:"value"`
	Date  AnalyticsTime `json:"date"`
}

// AnalyticsOverTimeStats is one series in an over-time KPI.
type AnalyticsOverTimeStats struct {
	Label  string                   `json:"label"`
	Values []AnalyticsOverTimeEntry `json:"values"`
}

// AnalyticsOverTimeResponse wraps any over-time KPI envelope.
type AnalyticsOverTimeResponse struct {
	Distribution []AnalyticsOverTimeStats `json:"distribution"`
}

// AnalyticsMeanTimeEntry is one row in the MTTR response.
type AnalyticsMeanTimeEntry struct {
	Label    string `json:"label"`
	Results  int64  `json:"results"`
	MeanTime int64  `json:"meanTime"`
}

// AnalyticsMeanTimeStats is the MTTR envelope.
type AnalyticsMeanTimeStats struct {
	MeanTimeData      []AnalyticsMeanTimeEntry `json:"meanTimeData"`
	MeanTimeStateData []AnalyticsMeanTimeEntry `json:"meanTimeStateData"`
	TotalResults      int64                    `json:"totalResults"`
}

// AnalyticsVulnerabilityStats is one row in most-common/most-aging/all
// vulnerabilities KPIs.
type AnalyticsVulnerabilityStats struct {
	VulnerabilityName string                         `json:"vulnerabilityName"`
	Total             int64                          `json:"total"`
	Severities        []AnalyticsLabeledResultsCount `json:"severities"`
}

// AnalyticsAllVulnerabilitiesResponse is the paged envelope for the
// allVulnerabilities KPI.
type AnalyticsAllVulnerabilitiesResponse struct {
	AllVulnerabilities []AnalyticsVulnerabilityStats `json:"allVulnerabilities"`
	Page               int                           `json:"page"`
}

// AnalyticsIDEStatEntry is one row in the ideTotal KPI.
type AnalyticsIDEStatEntry struct {
	Label      string  `json:"label"`
	Scans      float32 `json:"scans"`
	Developers float32 `json:"developers"`
}

// AnalyticsIDETotalResponse wraps the ideTotal envelope.
type AnalyticsIDETotalResponse struct {
	IDEData []AnalyticsIDEStatEntry `json:"ideData"`
}

// AnalyticsIDEOverTimeEntry is one point in the ideOvertime KPI.
type AnalyticsIDEOverTimeEntry struct {
	Label      string        `json:"label"`
	Scans      float32       `json:"scans"`
	Developers float32       `json:"developers"`
	Date       AnalyticsTime `json:"date"`
}

// AnalyticsIDEOverTimeDistribution is one series in the ideOvertime KPI.
type AnalyticsIDEOverTimeDistribution struct {
	Label  string                      `json:"label"`
	Values []AnalyticsIDEOverTimeEntry `json:"values"`
}

// AnalyticsIDEOverTimeResponse wraps the ideOvertime envelope.
type AnalyticsIDEOverTimeResponse struct {
	Distribution []AnalyticsIDEOverTimeDistribution `json:"distribution"`
}
