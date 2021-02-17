package attachinterfaces

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// List makes a request against the zun API to list the container's interfaces.
func List(client *gophercloud.ServiceClient, containerID string) pagination.Pager {
	return pagination.NewPager(client, listInterfaceURL(client, containerID), func(r pagination.PageResult) pagination.Page {
		return InterfacePage{pagination.SinglePageBase(r)}
	})
}

// Get requests details on a single interface attachment by the container and port IDs.
func Get(client *gophercloud.ServiceClient, containerID, portID string) (r GetResult) {
	allPages, err := List(client, containerID).AllPages()
	if err != nil {
		r.Err = err
		return
	}

	l, _ := ExtractInterfaces(allPages)
	for i := range l {
		if l[i].PortID == portID {
			// Found!
			r.Body = l[i]
			break
		}
	}
	return
}

// CreateOpts specifies parameters of a new interface attachment.
type CreateOpts struct {
	// PortID is the ID of the port for which you want to create an interface.
	// The NetworkID and PortID parameters are mutually exclusive.
	// If you do not specify the PortID parameter, the OpenStack Networking API
	// v2.0 allocates a port and creates an interface for it on the network.
	PortID string `q:"port,omitempty"`

	// NetworkID is the ID of the network for which you want to create an interface.
	// The NetworkID and PortID parameters are mutually exclusive.
	// If you do not specify the NetworkID parameter, the OpenStack Networking
	// API v2.0 uses the network information cache that is associated with the instance.
	NetworkID string `q:"network,omitempty"`

	// Fixed IP addresses. If you request a specific FixedIP address without a
	// NetworkID, the request returns a Bad Request (400) response code.
	FixedIP string `q:"fixed_ip,omitempty"`
}

// Create requests the creation of a new interface attachment on the server.
func Create(client *gophercloud.ServiceClient, containerID string, opts CreateOpts) (r CreateResult) {
	qs, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(attachInterfaceURL(client, containerID)+qs.String(), nil, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeleteOpts specifies parameters of an interface detattachment.
type DeleteOpts struct {
	// PortID is the ID of the port for which you want to remove an interface.
	// The NetworkID and PortID parameters are mutually exclusive.
	PortID string `q:"port,omitempty"`

	// NetworkID is the ID of the network for which you want to remove an interface.
	// The NetworkID and PortID parameters are mutually exclusive.
	NetworkID string `q:"network,omitempty"`
}

// Delete makes a request against the zun API to detach a single interface from the server.
// It needs server, and a network ID or port ID to make a such request.
func Delete(client *gophercloud.ServiceClient, containerID string, opts DeleteOpts) (r DeleteResult) {
	qs, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(detachInterfaceURL(client, containerID)+qs.String(), nil, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
