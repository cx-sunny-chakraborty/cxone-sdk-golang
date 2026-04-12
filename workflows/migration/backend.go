package migration

import "github.com/checkmarx-open-labs/cxone-sdk-golang/internal/transport"

// Backend bundles the dependencies the migration workflow needs.
type Backend struct {
	Executor *transport.Executor
	BaseURL  string
}
