package containers

import "github.com/gophercloud/gophercloud"

func getURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("containers", id)
}

func createURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("containers")
}

// `listURL` is a pure function. `listURL(c)` is a URL for which a GET
// request will respond with a list of containers in the service `c`.
func listURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("containers")
}

func deleteURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("containers", id)
}

func startURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("containers", id, "start")
}
