package roles

import "github.com/gophercloud/gophercloud"

const (
	ExtPath  = "OS-KSADM"
	RolePath = "roles"
	UserPath = "users"
)

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(ExtPath, RolePath, id)
}

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(ExtPath, RolePath)
}

func listURL(c *gophercloud.ServiceClient) string{
	return c.ServiceURL("roles")
}

func userRoleURL(c *gophercloud.ServiceClient, tenantID, userID, roleID string) string {
	return c.ServiceURL("projects", tenantID, UserPath, userID, RolePath,  roleID)
}
