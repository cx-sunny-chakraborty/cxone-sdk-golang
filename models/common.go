// Package models holds the strongly-typed request and response models for
// the Checkmarx One REST API.
//
// The wire format (JSON tags, field names, types) mirrors the API exactly so
// the SDK is forward compatible: unknown fields are tolerated by Go's
// encoding/json by default and new fields landing on the server side will
// not break decoding (CLAUDE.md §2 #9, §5).
//
// Naming uses Go-idiomatic exported PascalCase with explicit JSON tags. The
// struct names match the names in the Checkmarx One API documentation
// wherever possible so that cross-referencing with the swagger spec is
// straightforward.
package models

import "encoding/json"

// ErrorModel is the standard error envelope returned by Checkmarx One when an
// API request fails with a 4xx or 5xx status. The SDK's endpoint layer
// decodes this body and surfaces it via *cxerrors.ResponseError when an
// unexpected status is observed.
type ErrorModel struct {
	Message string `json:"message"`
	Type    string `json:"type,omitempty"`
	Code    int    `json:"code,omitempty"`
}

// WebError is an alternative error envelope used by some Checkmarx One
// services. The Data field carries opaque service-specific JSON.
type WebError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// PaginationEnvelope is embedded in collection-style response models. The
// concrete types declare the array element separately:
//
//	type ProjectsCollection struct {
//	    PaginationEnvelope
//	    Projects []ProjectResponseModel `json:"projects"`
//	}
type PaginationEnvelope struct {
	TotalCount         uint `json:"totalCount"`
	FilteredTotalCount uint `json:"filteredTotalCount"`
}
