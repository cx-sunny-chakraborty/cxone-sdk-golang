package cxone

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientBuilderRequiresRegion(t *testing.T) {
	_, err := NewClient().Tenant("acme").AgentName("ua").APIKey("k").Build()
	var ce *ConfigurationError
	if !errors.As(err, &ce) || ce.Field != "Region" {
		t.Fatalf("got %v", err)
	}
}

func TestClientBuilderRequiresTenant(t *testing.T) {
	_, err := NewClient().Region(RegionUS).AgentName("ua").APIKey("k").Build()
	var ce *ConfigurationError
	if !errors.As(err, &ce) || ce.Field != "Tenant" {
		t.Fatalf("got %v", err)
	}
}

func TestClientBuilderRequiresAgentName(t *testing.T) {
	_, err := NewClient().Region(RegionUS).Tenant("acme").APIKey("k").Build()
	var ce *ConfigurationError
	if !errors.As(err, &ce) || ce.Field != "AgentName" {
		t.Fatalf("got %v", err)
	}
}

func TestClientBuilderRequiresAuth(t *testing.T) {
	_, err := NewClient().Region(RegionUS).Tenant("acme").AgentName("ua").Build()
	var ce *ConfigurationError
	if !errors.As(err, &ce) || ce.Field != "Auth" {
		t.Fatalf("got %v", err)
	}
}

func TestClientBuilderBuildsValidClient(t *testing.T) {
	c, err := NewClient().
		Region(RegionUS).
		Tenant("acme").
		AgentName("ua").
		APIKey("k").
		Build()
	if err != nil {
		t.Fatal(err)
	}
	if c.Region().Name != "US" {
		t.Errorf("region = %q", c.Region().Name)
	}
	if c.Tenant() != "acme" {
		t.Errorf("tenant = %q", c.Tenant())
	}
	if len(c.CorrelationID()) == 0 {
		t.Errorf("correlationID empty")
	}
	if err := c.Close(); err != nil {
		t.Errorf("close: %v", err)
	}
	// Idempotent close.
	if err := c.Close(); err != nil {
		t.Errorf("second close: %v", err)
	}
}

func TestClientBuilderCustomRegionAndHTTPClient(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	region, err := NewCustomRegion("iam.example.com", "ast.example.com")
	if err != nil {
		t.Fatal(err)
	}
	c, err := NewClient().
		Region(region).
		Tenant("acme").
		AgentName("ua").
		APIKey("k").
		HTTPClient(srv.Client()).
		Build()
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	if c.Region().Name != "custom" {
		t.Errorf("region.Name = %q", c.Region().Name)
	}
}

func TestClientCorrelationIDIsStableAcrossLifetime(t *testing.T) {
	c, err := NewClient().
		Region(RegionUS).
		Tenant("t").
		AgentName("ua").
		APIKey("k").
		Build()
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	id1 := c.CorrelationID()
	id2 := c.CorrelationID()
	if id1 != id2 {
		t.Errorf("correlationID changed across calls: %q vs %q", id1, id2)
	}
}

func TestClientBuilderTwoClientsHaveDifferentCorrelationIDs(t *testing.T) {
	build := func() *Client {
		c, err := NewClient().Region(RegionUS).Tenant("t").AgentName("ua").APIKey("k").Build()
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { c.Close() })
		return c
	}
	a, b := build(), build()
	if a.CorrelationID() == b.CorrelationID() {
		t.Errorf("two clients share a correlationID: %q", a.CorrelationID())
	}
}

func TestClientBuilderRejectsBadProxy(t *testing.T) {
	_, err := NewClient().
		Region(RegionUS).
		Tenant("t").
		AgentName("ua").
		APIKey("k").
		Proxy("://bogus").
		Build()
	var ce *ConfigurationError
	if !errors.As(err, &ce) || ce.Field != "Proxy" {
		t.Fatalf("got %v", err)
	}
}
