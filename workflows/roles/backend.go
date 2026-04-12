package roles

import "github.com/checkmarx-open-labs/cxone-sdk-golang/internal/transport"

// Backend bundles the dependencies the role workflow needs. IAM role
// endpoints target the Keycloak admin realm.
type Backend struct {
	Executor    *transport.Executor
	IAMAdminURL string
}
