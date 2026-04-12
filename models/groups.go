package models

// Group is one IAM group entry.
//
// This type backs both the lightweight PIP groups endpoint (which returns
// only ID + Name) and the Keycloak admin groups endpoint (which returns the
// full hierarchy, role bindings, attribute bag, etc.). The JSON decoder
// tolerates missing fields; legacy callers that only need ID and Name still
// work.
type Group struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`

	// Path is the tenant-absolute path of the group ("/Parent/Child").
	Path string `json:"path,omitempty"`

	// ParentID is the uuid of the immediate parent group, or empty for
	// top-level groups.
	ParentID string `json:"parentId,omitempty"`

	// SubGroups is the immediate children of this group, as returned by
	// /groups/{id}/children. Populated on demand.
	SubGroups []Group `json:"subGroups,omitempty"`

	// SubGroupCount is the number of immediate children, as reported by
	// Keycloak on /groups/{id}. May differ from len(SubGroups) when the
	// hierarchy has not been walked yet.
	SubGroupCount uint64 `json:"subGroupCount,omitempty"`

	// RealmRoles is the list of realm-role names directly bound to this
	// group (not inherited).
	RealmRoles []string `json:"realmRoles,omitempty"`

	// ClientRoles maps a Keycloak client id (e.g. "ast-app") to the list
	// of client-role names bound to this group.
	ClientRoles map[string][]string `json:"clientRoles,omitempty"`

	// Attributes is the free-form attribute bag Keycloak stores alongside
	// the group.
	Attributes map[string][]string `json:"attributes,omitempty"`
}

// GroupCountResponse is the envelope Keycloak returns from /groups/count.
type GroupCountResponse struct {
	Count uint64 `json:"count"`
}

// GroupFilter carries the supported query parameters for GET /groups and
// GET /groups/count.
type GroupFilter struct {
	First int
	Max   int

	BriefRepresentation *bool
	PopulateHierarchy   *bool
	Search              string
	Exact               *bool
	// TopLevel selects only top-level groups on /groups/count (the "top"
	// query parameter). Not meaningful for /groups itself.
	TopLevel bool
}
