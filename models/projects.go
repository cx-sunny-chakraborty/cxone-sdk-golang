package models

import "time"

// Project is the Checkmarx One project create/update payload. Field naming
// and JSON tags match the Checkmarx One API contract — see
// ast-cli/internal/wrappers/projects.go for the source of truth.
type Project struct {
	Name           string            `json:"name,omitempty"`
	RepoURL        string            `json:"repoUrl,omitempty"`
	MainBranch     string            `json:"mainBranch,omitempty"`
	Origin         string            `json:"origin,omitempty"`
	ScmRepoID      string            `json:"scmRepoId,omitempty"`
	Tags           map[string]string `json:"tags,omitempty"`
	Groups         []string          `json:"groups,omitempty"`
	PrivatePackage bool              `json:"privatePackage,omitempty"`
	ApplicationIDs []string          `json:"applicationIds,omitempty"`
}

// ProjectResponseModel is the server-returned representation of a single
// project.
type ProjectResponseModel struct {
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	CreatedAt      time.Time         `json:"createdAt"`
	UpdatedAt      time.Time         `json:"updatedAt"`
	Groups         []string          `json:"groups"`
	Tags           map[string]string `json:"tags"`
	RepoURL        string            `json:"repoUrl"`
	MainBranch     string            `json:"mainBranch"`
	Origin         string            `json:"origin,omitempty"`
	ScmRepoID      string            `json:"scmRepoId,omitempty"`
	ApplicationIDs []string          `json:"applicationIds"`
}

// ProjectsCollectionResponseModel is the paginated list-projects response.
type ProjectsCollectionResponseModel struct {
	TotalCount         uint                   `json:"totalCount"`
	FilteredTotalCount uint                   `json:"filteredTotalCount"`
	Projects           []ProjectResponseModel `json:"projects"`
}

// ProjectConfiguration is one entry in the configuration patch payload sent
// to PATCH /api/configuration/project.
type ProjectConfiguration struct {
	Key           string `json:"key"`
	Name          string `json:"name"`
	Category      string `json:"category"`
	OriginLevel   string `json:"originLevel"`
	Value         string `json:"value"`
	ValueType     string `json:"valuetype"`
	AllowOverride bool   `json:"allowOverride"`
}

// BranchList is the response from GET /api/projects/branches.
//
// The endpoint returns a JSON array of strings, not an envelope, so it is
// modeled directly as a Go slice. The named type lets callers reference
// "models.BranchList" in signatures.
type BranchList []string

// ProjectPatch is the partial-update body for PATCH /api/projects/{id}
// (Checkmarx One 3.41+). Only non-nil fields are sent.
type ProjectPatch struct {
	Name       *string            `json:"name,omitempty"`
	Tags       *map[string]string `json:"tags,omitempty"`
	Groups     *[]string          `json:"groups,omitempty"`
	RepoURL    *string            `json:"repoUrl,omitempty"`
	MainBranch *string            `json:"mainBranch,omitempty"`
}

// ProjectApplicationsModel is the body of POST/DELETE
// /api/projects/{id}/applications (direct-association variant, gated by
// the DIRECT_APP_ASSOCIATION_ENABLED feature flag).
type ProjectApplicationsModel struct {
	ApplicationIDs []string `json:"applicationIds"`
}
