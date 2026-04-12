//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	cxone "github.com/checkmarx-open-labs/cxone-sdk-golang"
)

func TestListSCMIntegrations(t *testing.T) {
	dualAuth(t, func(t *testing.T, c *cxone.Client) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		integrations, err := c.Advanced().SCM().ListIntegrations(ctx)
		if err != nil {
			t.Fatalf("list SCM integrations: %v", err)
		}
		t.Logf("SCM integrations = %d", len(integrations))
		for _, integ := range integrations {
			t.Logf("  %s (type=%s repos=%d)", integ.DisplayName, integ.Type, integ.RepoCount)
		}
	})
}
