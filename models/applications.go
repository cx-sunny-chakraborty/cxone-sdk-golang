package models

import "time"

// ApplicationsResponseModel is the paginated list-applications response.
type ApplicationsResponseModel struct {
	TotalCount         int           `json:"totalCount"`
	FilteredTotalCount int           `json:"filteredTotalCount"`
	Applications       []Application `json:"applications"`
}

// Application is the server-returned representation of one application.
type Application struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Criticality int               `json:"criticality"`
	Rules       []ApplicationRule `json:"rules"`
	ProjectIDs  []string          `json:"projectIds"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
	Tags        map[string]string `json:"tags"`
	Type        string            `json:"type"`
}

// ApplicationConfiguration is the body of PUT /api/applications/{id}.
type ApplicationConfiguration struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Type        string            `json:"type"`
	Criticality int               `json:"criticality"`
	Rules       []ApplicationRule `json:"rules"`
	Tags        map[string]string `json:"tags"`
}

// AssociateProjectModel is the body of POST /api/applications/{id}/projects.
type AssociateProjectModel struct {
	ProjectIDs []string `json:"projectIds"`
}

// ApplicationCreateRequest is the body of POST /api/applications.
//
// The Checkmarx One application create endpoint rejects omitted rule/tag
// fields on some tenants; defaults are initialized to non-nil empty slices
// and maps before send.
type ApplicationCreateRequest struct {
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Criticality int               `json:"criticality,omitempty"`
	Rules       []ApplicationRule `json:"rules"`
	Tags        map[string]string `json:"tags"`
}

// ApplicationPatch is the partial-update body for PATCH /api/applications/{id}
// (Checkmarx One 3.41+). Only non-nil fields are sent.
type ApplicationPatch struct {
	Name        *string            `json:"name,omitempty"`
	Description *string            `json:"description,omitempty"`
	Criticality *int               `json:"criticality,omitempty"`
	Rules       *[]ApplicationRule `json:"rules,omitempty"`
	Tags        *map[string]string `json:"tags,omitempty"`
}

// ApplicationProjectsModel is the body of POST/DELETE
// /api/applications/{id}/projects (direct-association variant) when the
// DIRECT_APP_ASSOCIATION_ENABLED feature flag is on.
type ApplicationProjectsModel struct {
	Projects []string `json:"projects"`
}

// ApplicationRule is one criticality / inclusion rule on an application.
//
// Named ApplicationRule (rather than the more generic Rule used in ast-cli)
// to avoid colliding with future rule types from other Checkmarx One areas.
type ApplicationRule struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Value string `json:"value"`
}
