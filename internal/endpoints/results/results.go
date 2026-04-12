// Package results implements the Checkmarx One Results API (low-level layer).
//
// The /api/results endpoint returns a paginated, cross-engine view of all
// findings for a single scan. Engine-specific result pages (sast-results,
// kics-results) live in their own subpackages and may be added later — most
// callers should use this aggregate endpoint and discriminate by
// ScanResult.Type.
package results

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/checkmarx-open-labs/cxone-sdk-golang/internal/transport"
	"github.com/checkmarx-open-labs/cxone-sdk-golang/models"
	"github.com/checkmarx-open-labs/cxone-sdk-golang/pagination"
)

// Path is the relative path to the results endpoint.
const Path = "api/results"

// List returns a single page of scan results. GET /api/results → 200.
//
// query MUST contain "scan-id" set to the scan whose results to fetch.
// Other supported query parameters: "limit" (default 10000),
// "offset" (default 0), "sort" (default "-severity"), and engine-specific
// filters such as "severity", "state", "status".
func List(ctx context.Context, e *transport.Executor, baseURL string, query url.Values) (*models.ScanResultsCollection, error) {
	var out models.ScanResultsCollection
	if err := transport.DoJSON(ctx, e, http.MethodGet,
		transport.JoinURL(baseURL, Path),
		query, nil, []int{http.StatusOK}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Iter returns an item-yielding iterator over the scan results matching
// query. Caller MUST set "scan-id" in query; the iterator preserves it
// across pages and overrides only "limit" and "offset".
//
// The /api/results endpoint can return very large pages (the platform
// default is 10,000) — pick cfg.PageSize to match the largest payload your
// caller can buffer comfortably.
func Iter(e *transport.Executor, baseURL string, query url.Values, cfg pagination.Config) *pagination.Iterator[*models.ScanResult] {
	cfg.OffsetIsByCount = true
	if cfg.PageSize == 0 {
		cfg.PageSize = 1000
	}
	return pagination.New(cfg, func(ctx context.Context, offset, limit int) ([]*models.ScanResult, error) {
		q := make(url.Values, len(query)+2)
		for k, vs := range query {
			cp := make([]string, len(vs))
			copy(cp, vs)
			q[k] = cp
		}
		q.Set("offset", strconv.Itoa(offset))
		q.Set("limit", strconv.Itoa(limit))
		coll, err := List(ctx, e, baseURL, q)
		if err != nil {
			return nil, err
		}
		return coll.Results, nil
	})
}
