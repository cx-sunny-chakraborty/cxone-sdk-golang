package cxone

import (
	"errors"
	"testing"
)

func TestAllRegionsCatalog(t *testing.T) {
	regions := AllRegions()
	if len(regions) != 9 {
		t.Fatalf("expected 9 regions, got %d", len(regions))
	}
	// Spot-check the canonical hosts from CLAUDE.md §13.
	expected := map[string][2]string{
		"US":        {"iam.checkmarx.net", "ast.checkmarx.net"},
		"US2":       {"us.iam.checkmarx.net", "us.ast.checkmarx.net"},
		"EU":        {"eu.iam.checkmarx.net", "eu.ast.checkmarx.net"},
		"EU2":       {"eu-2.iam.checkmarx.net", "eu-2.ast.checkmarx.net"},
		"DEU":       {"deu.iam.checkmarx.net", "deu.ast.checkmarx.net"},
		"ANZ":       {"anz.iam.checkmarx.net", "anz.ast.checkmarx.net"},
		"India":     {"ind.iam.checkmarx.net", "ind.ast.checkmarx.net"},
		"Singapore": {"sng.iam.checkmarx.net", "sng.ast.checkmarx.net"},
		"UAE":       {"mea.iam.checkmarx.net", "mea.ast.checkmarx.net"},
	}
	for _, r := range regions {
		want, ok := expected[r.Name]
		if !ok {
			t.Errorf("unexpected region in catalog: %q", r.Name)
			continue
		}
		if r.AuthHost != want[0] || r.APIHost != want[1] {
			t.Errorf("region %s: got auth=%q api=%q want auth=%q api=%q",
				r.Name, r.AuthHost, r.APIHost, want[0], want[1])
		}
		if r.Scheme != "https" {
			t.Errorf("region %s: scheme = %q, want https", r.Name, r.Scheme)
		}
	}
}

func TestRegionTokenURL(t *testing.T) {
	got, err := RegionUS.TokenURL("acme")
	if err != nil {
		t.Fatal(err)
	}
	want := "https://iam.checkmarx.net/auth/realms/acme/protocol/openid-connect/token"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRegionTokenURLEmptyTenant(t *testing.T) {
	_, err := RegionUS.TokenURL("")
	var endpointErr *EndpointError
	if !errors.As(err, &endpointErr) {
		t.Fatalf("expected *EndpointError, got %T: %v", err, err)
	}
}

func TestRegionAPIBaseURL(t *testing.T) {
	if got, want := RegionEU.APIBaseURL(), "https://eu.ast.checkmarx.net/api/"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRegionDisplayRootURL(t *testing.T) {
	if got, want := RegionEU.DisplayRootURL(), "https://eu.ast.checkmarx.net/"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNewCustomRegionValid(t *testing.T) {
	r, err := NewCustomRegion("iam.example.com", "ast.example.com")
	if err != nil {
		t.Fatal(err)
	}
	if r.Name != "custom" || r.Scheme != "https" {
		t.Errorf("unexpected region: %+v", r)
	}
	if got, _ := r.TokenURL("t"); got != "https://iam.example.com/auth/realms/t/protocol/openid-connect/token" {
		t.Errorf("token URL: %q", got)
	}
}

func TestNewCustomRegionRejectsBadHost(t *testing.T) {
	cases := []struct {
		name, auth, api string
	}{
		{"empty auth", "", "ast.example.com"},
		{"empty api", "iam.example.com", ""},
		{"scheme in auth", "https://iam.example.com", "ast.example.com"},
		{"path in api", "iam.example.com", "ast.example.com/api"},
		{"whitespace", "iam.example.com ", "ast.example.com"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewCustomRegion(tc.auth, tc.api)
			var ee *EndpointError
			if !errors.As(err, &ee) {
				t.Fatalf("expected *EndpointError, got %T: %v", err, err)
			}
		})
	}
}

func TestRegionWithSchemeRejectsBadScheme(t *testing.T) {
	_, err := RegionUS.WithScheme("ftp")
	var ee *EndpointError
	if !errors.As(err, &ee) {
		t.Fatalf("expected *EndpointError, got %T: %v", err, err)
	}
}

func TestRegionWithSchemeHTTPOverride(t *testing.T) {
	r, err := RegionUS.WithScheme("http")
	if err != nil {
		t.Fatal(err)
	}
	if got, _ := r.TokenURL("acme"); got != "http://iam.checkmarx.net/auth/realms/acme/protocol/openid-connect/token" {
		t.Errorf("expected http scheme in URL: %q", got)
	}
}
