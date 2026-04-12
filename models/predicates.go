package models

import "time"

// Predicate state values for SAST/KICS triage. Use these constants instead
// of bare strings.
const (
	PredicateStateToVerify               = "ToVerify"
	PredicateStateConfirmed              = "Confirmed"
	PredicateStateNotExploitable         = "NotExploitable"
	PredicateStateProposedNotExploitable = "ProposedNotExploitable"
	PredicateStateUrgent                 = "Urgent"
)

// PredicateRequest is the body of POST /api/sast-results-predicates and the
// equivalent KICS endpoint. Either State (built-in) or CustomStateID (one of
// the tenant's custom states) should be set, not both.
type PredicateRequest struct {
	SimilarityID  string  `json:"similarityId"`
	ProjectID     string  `json:"projectId"`
	State         *string `json:"state,omitempty"`
	CustomStateID *int    `json:"customStateId,omitempty"`
	Comment       string  `json:"comment"`
	Severity      string  `json:"severity"`
}

// BasePredicate is the shared shape used by historical predicate entries.
type BasePredicate struct {
	SimilarityID string `json:"similarityId"`
	ProjectID    string `json:"projectId"`
	State        string `json:"state"`
	Severity     string `json:"severity"`
	Comment      string `json:"comment"`
	StateID      int    `json:"stateId"`
}

// Predicate is one historical state-change entry.
type Predicate struct {
	BasePredicate
	ID        string    `json:"ID"`
	CreatedBy string    `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
}

// PredicateHistory groups predicates per project for a given similarity id.
type PredicateHistory struct {
	ProjectID    string      `json:"projectId"`
	SimilarityID string      `json:"similarityId"`
	Predicates   []Predicate `json:"predicates"`
	TotalCount   int         `json:"totalCount"`
}

// PredicatesCollectionResponseModel is the GET /api/sast-results-predicates/{id}
// response envelope.
type PredicatesCollectionResponseModel struct {
	ScaResponse                any                `json:"scaPredicate,omitempty"`
	PredicateHistoryPerProject []PredicateHistory `json:"predicateHistoryPerProject"`
	TotalCount                 int                `json:"totalCount"`
}

// CustomState is one tenant-defined custom triage state.
type CustomState struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	IsAllowed bool   `json:"-"`
}

// CustomStateCreateRequest is the minimal body for POST /api/custom-states.
type CustomStateCreateRequest struct {
	Name string `json:"name"`
}

// ResultsChangeHistoryEntry is one row returned by
// GET /api/sast-results-predicates/changelog.
type ResultsChangeHistoryEntry struct {
	SimilarityID string    `json:"similarityId"`
	ProjectID    string    `json:"projectId"`
	ScanID       string    `json:"scanId,omitempty"`
	State        string    `json:"state,omitempty"`
	Severity     string    `json:"severity,omitempty"`
	Comment      string    `json:"comment,omitempty"`
	CreatedBy    string    `json:"createdBy,omitempty"`
	CreatedAt    time.Time `json:"createdAt,omitempty"`
}

// ResultsChangeHistoryResponse is the envelope returned by the changelog
// endpoint. TotalSimilarityIDs is a free-form string like
// "presenting 10 of 29"; higher-level helpers parse it.
type ResultsChangeHistoryResponse struct {
	Results            []ResultsChangeHistoryEntry `json:"results"`
	TotalSimilarityIDs string                      `json:"totalSimilarityIds"`
}

// ScaPredicateRequest is the body of the SCA package-vulnerabilities
// predicate endpoint.
type ScaPredicateRequest struct {
	PackageName     string      `json:"packageName"`
	PackageVersion  string      `json:"packageVersion"`
	PackageManager  string      `json:"packageManager"`
	VulnerabilityID string      `json:"vulnerabilityId"`
	ProjectIDs      []string    `json:"projectIds"`
	Actions         []ScaAction `json:"actions"`
}

// ScaAction is one action inside a [ScaPredicateRequest].
type ScaAction struct {
	ActionType string `json:"actionType"`
	Value      string `json:"value"`
	Comment    string `json:"comment"`
}
