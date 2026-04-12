package models

// User is one IAM user as returned by the Keycloak admin endpoint
// (/auth/admin/realms/{tenant}/users).
//
// Field naming mirrors the wire format exactly. Unknown fields are tolerated
// by the JSON decoder — the Checkmarx One user payload carries many Keycloak
// attributes we do not surface.
type User struct {
	ID        string `json:"id,omitempty"`
	Username  string `json:"username,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email,omitempty"`
	Enabled   bool   `json:"enabled"`

	EmailVerified *bool `json:"emailVerified,omitempty"`

	// Attributes holds the raw Keycloak attribute bag (lastLogin, locale, etc).
	// Values are string slices in Keycloak's wire format.
	Attributes map[string][]string `json:"attributes,omitempty"`

	// RequiredActions is the list of Keycloak "required actions" (e.g.
	// "VERIFY_EMAIL"). Used when creating users that must complete setup
	// on first login.
	RequiredActions []string `json:"requiredActions,omitempty"`

	// FederatedIdentities associates the user with external identity
	// providers. Required for SAML users (see [UserFederatedIdentity]).
	FederatedIdentities []UserFederatedIdentity `json:"federatedIdentities,omitempty"`

	// Totp is the "time-based one-time password" flag Keycloak sets when
	// 2FA is enrolled; CreateSAMLUser sets this to false.
	Totp *bool `json:"totp,omitempty"`
}

// UserFederatedIdentity links an SDK user to an external identity provider
// entry. Populated by [CreateSAMLUser].
type UserFederatedIdentity struct {
	IdentityProvider string `json:"identityProvider"`
	UserID           string `json:"userId"`
	UserName         string `json:"userName"`
}

// UserFilter carries the supported query parameters for GET /users and
// GET /users/count on the Keycloak admin API.
//
// All fields are optional. Pointer fields represent tri-state booleans so
// callers can distinguish "not set" from "explicitly false".
type UserFilter struct {
	// First is the offset (Keycloak calls this "first").
	First int
	// Max is the page size (Keycloak calls this "max"). 0 means unset.
	Max int

	BriefRepresentation *bool
	Email               string
	EmailVerified       *bool
	Enabled             *bool
	Exact               *bool
	FirstName           string
	IDPAlias            string
	IDPUserID           string
	Q                   string
	Search              string
	Username            string
}
