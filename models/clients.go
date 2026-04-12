package models

// OIDCClient is one Keycloak/OIDC client entry as returned by the
// /auth/admin/realms/{tenant}/clients endpoint.
//
// The Checkmarx One platform uses OIDC clients for two purposes: the built-in
// "ast-app" client (the Checkmarx One application itself) and tenant-created
// service accounts used for API automation. Both flow through this type.
//
// The wire format is Keycloak's — additional fields (redirectUris, protocol,
// attributes bag, etc.) are tolerated and ignored by the decoder. Callers
// needing round-trippable access to the raw attribute map should use
// [OIDCClient.Attributes].
type OIDCClient struct {
	ID       string `json:"id,omitempty"`
	ClientID string `json:"clientId,omitempty"`
	Name     string `json:"name,omitempty"`
	Enabled  bool   `json:"enabled"`

	// Secret is the client secret. Only populated on create responses and
	// [ClientsEndpoint.GetSecret] / [ClientsEndpoint.RegenerateSecret]
	// calls.
	Secret string `json:"secret,omitempty"`

	Protocol               string   `json:"protocol,omitempty"`
	PublicClient           bool     `json:"publicClient,omitempty"`
	ServiceAccountsEnabled bool     `json:"serviceAccountsEnabled,omitempty"`
	StandardFlowEnabled    bool     `json:"standardFlowEnabled,omitempty"`
	FrontchannelLogout     bool     `json:"frontchannelLogout,omitempty"`
	RedirectURIs           []string `json:"redirectUris,omitempty"`

	// Attributes is Keycloak's free-form attribute bag. Values are always
	// strings on the wire; semantic decoding (expiry timestamps, email
	// lists, etc.) is the caller's responsibility.
	Attributes map[string]string `json:"attributes,omitempty"`
}

// OIDCClientScope is one client-scope entry.
type OIDCClientScope struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Protocol    string `json:"protocol,omitempty"`
}

// OIDCClientFilter carries the supported query parameters for GET /clients.
type OIDCClientFilter struct {
	First        int
	Max          int
	ClientID     string
	Q            string
	Search       *bool
	ViewableOnly *bool
}

// ClientSecretResponse is the minimal envelope Keycloak returns from
// GET/POST /clients/{id}/client-secret.
type ClientSecretResponse struct {
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}
