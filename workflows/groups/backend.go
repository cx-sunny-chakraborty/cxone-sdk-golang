package groups

import "github.com/checkmarx-open-labs/cxone-sdk-golang/internal/transport"

// Backend bundles the dependencies the group workflow needs. Group admin
// endpoints target the Keycloak admin realm.
type Backend struct {
	Executor    *transport.Executor
	IAMAdminURL string
}
