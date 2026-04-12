package clients

import "github.com/checkmarx-open-labs/cxone-sdk-golang/internal/transport"

// Backend bundles the dependencies the client workflow needs.
type Backend struct {
	Executor    *transport.Executor
	IAMAdminURL string
}
