package projects

import "github.com/gophercloud/gophercloud"

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("projects")
}

func getURL(client *gophercloud.ServiceClient, projectID string) string {
	return client.ServiceURL("projects", projectID)
}

func createURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("projects")
}

func deleteURL(client *gophercloud.ServiceClient, projectID string) string {
	return client.ServiceURL("projects", projectID)
}

func updateURL(client *gophercloud.ServiceClient, projectID string) string {
	return client.ServiceURL("projects", projectID)
}

//added by ccwings.cn
//list projects belongs user
func listByUserIdURL(client *gophercloud.ServiceClient,user_id string) string {
	return client.ServiceURL("users",user_id,"projects")
}