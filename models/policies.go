package models

// PolicyResponseModel is the response of GET /api/policies/evaluate.
type PolicyResponseModel struct {
	Status     string   `json:"status"`
	BreakBuild bool     `json:"breakBuild"`
	Policies   []Policy `json:"policies"`
}

// Policy is one entry inside [PolicyResponseModel].
type Policy struct {
	Name          string   `json:"policyName"`
	BreakBuild    bool     `json:"breakBuild"`
	Status        string   `json:"status"`
	Description   string   `json:"description"`
	RulesViolated []string `json:"rulesViolated"`
	Tags          []string `json:"tags"`
}

// PRPolicy is the compact policy shape used by the PR-decoration endpoints.
type PRPolicy struct {
	Name       string   `json:"policyName"`
	RulesNames []string `json:"rulesNames"`
	BreakBuild bool     `json:"breakBuild"`
}

// PolicyDetail is the full management-endpoint representation of a policy
// from GET /api/policy_management_service_uri/policies/v2. The evaluation
// shape ([Policy] above) carries a compact per-scan view; PolicyDetail is
// the CRUD/inventory view.
type PolicyDetail struct {
	ID            string   `json:"id,omitempty"`
	Name          string   `json:"name,omitempty"`
	Description   string   `json:"description,omitempty"`
	IsActivated   bool     `json:"isActivated,omitempty"`
	DefaultPolicy bool     `json:"defaultPolicy,omitempty"`
	Tags          []string `json:"tags,omitempty"`
	Projects      []string `json:"projects,omitempty"`
	Rules         []string `json:"rules,omitempty"`
}

// PolicyListResponse is the envelope for GET /policies/v2.
type PolicyListResponse struct {
	Policies           []PolicyDetail `json:"policies"`
	FilteredTotalCount uint64         `json:"filteredPoliciesCount"`
}

// PolicyViolation is one row returned by GET
// /policy_management_service_uri/incidents/filters.
type PolicyViolation struct {
	ViolationID uint64 `json:"violationId,omitempty"`
	PolicyID    string `json:"policyId,omitempty"`
	PolicyName  string `json:"policyName,omitempty"`
	ProjectID   string `json:"projectId,omitempty"`
	ProjectName string `json:"projectName,omitempty"`
	ScanID      string `json:"scanId,omitempty"`
	ScanDate    string `json:"scanDate,omitempty"`
	Status      string `json:"status,omitempty"`
	BreakBuild  bool   `json:"breakBuild,omitempty"`
}

// PolicyViolationListResponse is the envelope for the incidents/filters
// endpoint.
type PolicyViolationListResponse struct {
	Incidents          []PolicyViolation `json:"incidents"`
	FilteredTotalCount uint64            `json:"filteredIncidentsCount"`
}

// PolicyViolationDetails is the rich per-project/per-scan violation view
// returned by the evaluation endpoint when called with astProjectId + scanId.
type PolicyViolationDetails struct {
	ProjectID  string   `json:"astProjectId,omitempty"`
	ScanID     string   `json:"scanId,omitempty"`
	Policies   []Policy `json:"policies,omitempty"`
	Status     string   `json:"status,omitempty"`
	BreakBuild bool     `json:"breakBuild,omitempty"`
}

// PolicyFilter is the query shape for /policies/v2 and /policies/v2?count=1.
type PolicyFilter struct {
	Limit uint64
	Page  uint64
}

// PolicyViolationFilter is the query shape for the incidents/filters endpoint.
type PolicyViolationFilter struct {
	Limit     uint64
	Page      uint64
	ProjectID string
	ScanID    string
	PolicyID  string
}
