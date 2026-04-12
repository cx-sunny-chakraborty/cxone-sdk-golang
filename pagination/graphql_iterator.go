package pagination

import (
	"context"
	"encoding/json"
	"fmt"
)

// GraphQLExecutor is the contract this package needs to issue a GraphQL
// query. The SDK provides a concrete implementation in
// internal/transport/graphql.go ([transport.GraphQLClient]); workflow
// packages obtain one from cxone.Client and pass it here. Defining the
// contract as an interface in pagination keeps this file dependency-free
// from the transport layer.
type GraphQLExecutor interface {
	// Execute sends one GraphQL query and decodes the data field into out.
	// Implementations should surface non-nil GraphQL errors[] arrays as a
	// non-nil Go error.
	Execute(ctx context.Context, query string, variables map[string]any, out any) error
}

// SkipTakeNames lets callers override the default GraphQL variable names
// "skip" and "take" if their schema uses different names (e.g. "offset" /
// "first"). Zero values resolve to "skip" and "take".
type SkipTakeNames struct {
	Skip string
	Take string
}

func (n SkipTakeNames) skip() string {
	if n.Skip != "" {
		return n.Skip
	}
	return "skip"
}

func (n SkipTakeNames) take() string {
	if n.Take != "" {
		return n.Take
	}
	return "take"
}

// GraphQLDecode is the user-supplied callback that converts the raw GraphQL
// `data` payload into a slice of typed items.
//
// Why a callback rather than a JSON path string? GraphQL responses nest the
// items at a schema-specific location (e.g. data → tenantInventory →
// licenses → nodes), and modeling that path generically is more trouble
// than it is worth. The decode callback can simply unmarshal into a typed
// shape and return the inner slice. Examples:
//
//	decode := func(data json.RawMessage) ([]LicenseRecord, error) {
//	    var resp struct {
//	        TenantInventory struct {
//	            Licenses struct {
//	                Nodes []LicenseRecord `json:"nodes"`
//	            } `json:"licenses"`
//	        } `json:"tenantInventory"`
//	    }
//	    if err := json.Unmarshal(data, &resp); err != nil {
//	        return nil, err
//	    }
//	    return resp.TenantInventory.Licenses.Nodes, nil
//	}
type GraphQLDecode[T any] func(data json.RawMessage) ([]T, error)

// NewGraphQL constructs an [Iterator] over a paginated GraphQL query.
//
// The query MUST accept the skip/take variables (or whatever names are set
// in vars). On each page fetch the iterator injects {skip, take} into the
// variables map alongside extraVars and calls client.Execute. The
// graphqlExecutor adapter is responsible for any retries or error
// translation.
//
// extraVars contains the schema-specific variables (filters, project ids,
// etc.) that should be sent on every page. The iterator does not mutate
// extraVars; it makes a per-page copy.
//
// CLAUDE.md §10.2: this is the "iterator base" the SCA tenant-level
// analyses (TenantLicenses, TenantPackages, TenantRisks) build on. Those
// concrete iterators live in the workflows/sca/analysis package and supply
// the query string + decode callback for their respective schemas.
func NewGraphQL[T any](
	client GraphQLExecutor,
	query string,
	extraVars map[string]any,
	decode GraphQLDecode[T],
	cfg Config,
	vars SkipTakeNames,
) *Iterator[T] {
	cfg.OffsetIsByCount = true
	if cfg.PageSize == 0 {
		cfg.PageSize = 100
	}
	skipName := vars.skip()
	takeName := vars.take()

	return New(cfg, func(ctx context.Context, offset, limit int) ([]T, error) {
		// Per-page copy of variables; the iterator owns skip/take.
		merged := make(map[string]any, len(extraVars)+2)
		for k, v := range extraVars {
			merged[k] = v
		}
		merged[skipName] = offset
		merged[takeName] = limit

		// Execute returns the data payload via an out pointer; we ask for
		// a json.RawMessage so the decode callback can interpret the
		// schema-specific nesting.
		var raw json.RawMessage
		if err := client.Execute(ctx, query, merged, &raw); err != nil {
			return nil, err
		}
		if len(raw) == 0 {
			return nil, nil
		}
		items, err := decode(raw)
		if err != nil {
			return nil, fmt.Errorf("graphql decode: %w", err)
		}
		return items, nil
	})
}
