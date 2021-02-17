package attachinterfaces

import "github.com/gophercloud/gophercloud"

func listInterfaceURL(client *gophercloud.ServiceClient, containerID string) string {
	return client.ServiceURL("containers", containerID, "network_list")
}

func attachInterfaceURL(client *gophercloud.ServiceClient, containerID string) string {
	return client.ServiceURL("containers", containerID, "network_attach")
}

func detachInterfaceURL(client *gophercloud.ServiceClient, containerID string) string {
	return client.ServiceURL("containers", containerID, "network_detach")
}
