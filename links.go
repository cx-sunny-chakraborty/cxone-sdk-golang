package cxone

import "fmt"

// ProjectLink returns a browser-navigable URL to the project's overview page.
func (c *Client) ProjectLink(projectID string) string {
	return fmt.Sprintf("%sprojects/%s/overview", c.region.DisplayRootURL(), projectID)
}

// ScanLink returns a browser-navigable URL to the scan's results page.
func (c *Client) ScanLink(scanID, projectID string) string {
	return fmt.Sprintf("%sprojects/%s/scans/%s", c.region.DisplayRootURL(), projectID, scanID)
}

// UserLink returns a browser-navigable URL to the user's Keycloak admin page.
func (c *Client) UserLink(userID string) string {
	return fmt.Sprintf("%sauth/admin/%s/console/#/realms/%s/users/%s",
		c.iamRootURL(), c.tenant, c.tenant, userID)
}

// GroupLink returns a browser-navigable URL to the group's Keycloak admin page.
func (c *Client) GroupLink(groupID string) string {
	return fmt.Sprintf("%sauth/admin/%s/console/#/realms/%s/groups/%s",
		c.iamRootURL(), c.tenant, c.tenant, groupID)
}

// RoleLink returns a browser-navigable URL to the role's Keycloak admin page.
func (c *Client) RoleLink(roleID string) string {
	return fmt.Sprintf("%sauth/admin/%s/console/#/realms/%s/roles/%s",
		c.iamRootURL(), c.tenant, c.tenant, roleID)
}

// iamRootURL returns the IAM host root for building Keycloak console links.
func (c *Client) iamRootURL() string {
	return c.region.scheme() + "://" + c.region.AuthHost + "/"
}
