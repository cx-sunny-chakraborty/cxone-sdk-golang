// Package risks implements the Checkmarx One Risk Management and API
// Security risks-overview endpoints.
//
// These two endpoints come from different services on the platform but are
// closely related — both surface a "risk view" of an existing scan — so
// they share a single SDK package.
package risks

import (
	"context"
	"net/http"
	"net/url"

	"github.com/checkmarx-open-labs/cxone-sdk-golang/internal/transport"
	"github.com/checkmarx-open-labs/cxone-sdk-golang/models"
)

// Path constants for the risk endpoints.
const (
	// RiskManagementPathTemplate carries the {projectId} placeholder. The
	// scan id is supplied as a query parameter.
	RiskManagementPathTemplate = "api/risk-management/projects/%s/results"

	// APISecOverviewPathTemplate carries the {scanId} placeholder.
	APISecOverviewPathTemplate = "api/apisec/static/api/scan/%s/risks-overview"

	// APISecRisksPath is the paginated risks listing endpoint.
	APISecRisksPath = "api/apisec/static/api/risks"
)

// GetRiskManagement returns the ASPM risk-management view of a project's
// results for a given scan.
// GET /api/risk-management/projects/{projectId}/results?scanID={id} → 200.
func GetRiskManagement(ctx context.Context, e *transport.Executor, baseURL, projectID, scanID string) (*models.ASPMResult, error) {
	path := "api/risk-management/projects/" + projectID + "/results"
	q := url.Values{"scanID": {scanID}}
	var out models.ASPMResult
	if err := transport.DoJSON(ctx, e, http.MethodGet,
		transport.JoinURL(baseURL, path),
		q, nil, []int{http.StatusOK}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetAPISecOverview returns the API security risks overview for a scan.
// GET /api/apisec/static/api/scan/{scanId}/risks-overview → 200.
func GetAPISecOverview(ctx context.Context, e *transport.Executor, baseURL, scanID string) (*models.APISecResult, error) {
	path := "api/apisec/static/api/scan/" + scanID + "/risks-overview"
	var out models.APISecResult
	if err := transport.DoJSON(ctx, e, http.MethodGet,
		transport.JoinURL(baseURL, path),
		nil, nil, []int{http.StatusOK}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListAPISecRisks returns paginated API security risk entries.
// GET /api/apisec/static/api/risks?page&per_page&filtering&state&status&severity → 200.
func ListAPISecRisks(ctx context.Context, e *transport.Executor, baseURL string, query url.Values) (*models.APISecPaginatedResult, error) {
	var out models.APISecPaginatedResult
	if err := transport.DoJSON(ctx, e, http.MethodGet,
		transport.JoinURL(baseURL, APISecRisksPath),
		query, nil, []int{http.StatusOK}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
