package flavors

import (
	"github.com/gophercloud/gophercloud"
)

func getURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("flavors", id)
}

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("flavors", "detail")
}

func createURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("flavors")
}

func deleteURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("flavors", id)
}

func addAccessURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("flavors",id,"action")
}

func deleteAccessURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("flavors",id,"action")
}

func listAccessURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("flavors",id,"os-flavor-access")
}

