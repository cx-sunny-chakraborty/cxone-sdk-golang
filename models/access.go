package models

// AssignmentResponse is one row returned by the access-management entities-for
// endpoint.
type AssignmentResponse struct {
	EntityID     string `json:"entityID"`
	EntityType   string `json:"entityType"`
	EntityName   string `json:"entityName"`
	EntityRoles  []any  `json:"entityRoles"`
	ResourceID   string `json:"resourceID"`
	ResourceType string `json:"resourceType"`
	ResourceName string `json:"resourceName"`
}

// AssignmentPayload is the body of POST /api/access-management/.
type AssignmentPayload struct {
	EntityID     string `json:"entityID"`
	EntityType   string `json:"entityType"`
	EntityRoles  []any  `json:"entityRoles"`
	ResourceType string `json:"resourceType"`
	ResourceID   string `json:"resourceID"`
}

// AccessibleResource is one row in the [CheckAccessibleResources] response.
type AccessibleResource struct {
	ResourceID   string   `json:"resourceId"`
	ResourceType string   `json:"resourceType"`
	ResourceName string   `json:"resourceName"`
	Roles        []string `json:"roles"`
}

// AccessibleResourcesResponse is the envelope returned by
// GET /api/access-management/get-resources.
type AccessibleResourcesResponse struct {
	All       bool                 `json:"all"`
	Resources []AccessibleResource `json:"resources"`
}

// HasAccessResponse is the envelope returned by
// GET /api/access-management/has-access.
type HasAccessResponse struct {
	AccessGranted bool `json:"accessGranted"`
}
