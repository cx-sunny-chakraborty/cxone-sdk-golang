package models

// Role is a Keycloak role entry — either a realm (IAM) role or a client
// (application) role. Checkmarx One application roles live under the
// "ast-app" client id; IAM roles live at the realm level.
//
// The wire format field names are Keycloak's, preserved verbatim.
type Role struct {
	// ID is the role UUID.
	ID string `json:"id,omitempty"`

	// Name is the role name (e.g. "ast-admin").
	Name string `json:"name,omitempty"`

	// Description is the human-readable role description.
	Description string `json:"description,omitempty"`

	// ContainerID is the owning Keycloak container: a client uuid (for
	// application roles) or the tenant realm name (for IAM roles).
	ContainerID string `json:"containerId,omitempty"`

	// Composite is true if this role bundles other roles.
	Composite bool `json:"composite,omitempty"`

	// ClientRole is true for client-scoped roles, false for realm roles.
	ClientRole bool `json:"clientRole,omitempty"`

	// Attributes mirrors Keycloak's raw attribute bag.
	Attributes map[string][]string `json:"attributes,omitempty"`
}

// RoleComposites is the set of roles that compose a composite role, split
// into realm-level and client-level entries as Keycloak returns them.
type RoleComposites struct {
	Realm  []Role            `json:"realm,omitempty"`
	Client map[string][]Role `json:"client,omitempty"`
}
