package containers

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// CreateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Create operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type CreateOptsBuilder interface {
	ToCapsuleCreateMap() (map[string]interface{}, error)
}

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToContainerListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the container attributes you want to see returned. Marker and Limit are used
// for pagination.
type ListOpts struct {
	Name       string `q:"name"`
	Image      string `q:"image"`
	ProjectID  string `q:"project_id"`
	UserID     string `q:"user_id"`
	Memory     int    `q:"memory"`
	Host       string `q:"host"`
	TaskState  string `q:"task_state"`
	Status     string `q:"status"`
	AutoRemove string `q:"auto_remove"`
}

// ToContainerListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToContainerListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List makes a request against the API to list containers accessible to you.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToContainerListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ContainerPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get requests details on a single container, by ID.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(getURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 203},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateOpts is the common options struct used in this package's Create
// operation.
type CreateOpts struct {
	Name             string                       `json:"name,omitempty"`
	Image            string                       `json:"image" required:"true"`
	Command          []string                     `json:"command,omitempty"`
	CPU              float64                      `json:"cpu,omitempty"`
	Memory           string                       `json:"memory,omitempty"`
	Workdir          string                       `json:"workdir,omitempty"`
	Labels           map[string]string            `json:"labels,omitempty"`
	Environment      map[string]string            `json:"environment,omitempty"`
	RestartPolicy    map[string]string            `json:"restart_policy,omitempty"`
	Interactive      *bool                        `json:"interactive,omitempty"`
	TTY              *bool                        `json:"tty,omitempty"`
	ImageDriver      string                       `json:"image_driver,omitempty"`
	SecurityGroups   []string                     `json:"security_groups,omitempty"`
	Nets             []map[string]string          `json:"nets,omitempty"`
	Runtime          string                       `json:"runtime,omitempty"`
	Hostname         string                       `json:"hostname,omitempty"`
	AutoRemove       *bool                        `json:"auto_remove,omitempty"`
	AutoHeal         *bool                        `json:"auto_heal,omitempty"`
	AvailabilityZone string                       `json:"availability_zone,omitempty"`
	Hints            map[string]string            `json:"hints,omitempty"`
	Mounts           []map[string]string          `json:"mounts,omitempty"`
	Privileged       *bool                        `json:"privileged,omitempty"`
	Healthcheck      map[string]string            `json:"healthcheck,omitempty"`
	ExposedPorts     map[string]map[string]string `json:"exposed_ports,omitempty"`
	Host             string                       `json:"host,omitempty"`
	Entrypoint       []string                     `json:"entrypoint,omitempty"`
}

// ToCapsuleCreateMap assembles a request body based on the contents of
// a CreateOpts.
func (opts CreateOpts) ToCapsuleCreateMap() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	return b, err
}

// Create implements create container request.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToCapsuleCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(createURL(client), b, &r.Body, &gophercloud.RequestOpts{OkCodes: []int{202}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete implements Container delete request.
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(deleteURL(client, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// StopAndDelete implements Container stop and delete request.
func StopAndDelete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(deleteURL(client, id)+"?stop=True", nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ForceDelete implements Container force delete request.
func ForceDelete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(deleteURL(client, id)+"?force=True", nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Start implements start capsule request.
func Start(client *gophercloud.ServiceClient, id string) (r StartResult) {
	resp, err := client.Post(startURL(client, id), nil, nil, &gophercloud.RequestOpts{OkCodes: []int{202}})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
