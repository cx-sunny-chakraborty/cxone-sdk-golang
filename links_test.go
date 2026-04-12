package cxone_test

import (
	"strings"
	"testing"
)

func TestLinkBuilders(t *testing.T) {
	c := newWorkflowClient(t, nil)

	t.Run("ProjectLink", func(t *testing.T) {
		link := c.ProjectLink("proj-123")
		if !strings.Contains(link, "projects/proj-123/overview") {
			t.Errorf("link = %q", link)
		}
	})

	t.Run("ScanLink", func(t *testing.T) {
		link := c.ScanLink("scan-1", "proj-1")
		if !strings.Contains(link, "projects/proj-1/scans/scan-1") {
			t.Errorf("link = %q", link)
		}
	})

	t.Run("UserLink", func(t *testing.T) {
		link := c.UserLink("user-1")
		if !strings.Contains(link, "users/user-1") {
			t.Errorf("link = %q", link)
		}
		if !strings.Contains(link, "/auth/admin/") {
			t.Errorf("expected IAM host path, got %q", link)
		}
	})

	t.Run("GroupLink", func(t *testing.T) {
		link := c.GroupLink("group-1")
		if !strings.Contains(link, "groups/group-1") {
			t.Errorf("link = %q", link)
		}
	})

	t.Run("RoleLink", func(t *testing.T) {
		link := c.RoleLink("role-1")
		if !strings.Contains(link, "roles/role-1") {
			t.Errorf("link = %q", link)
		}
	})
}
