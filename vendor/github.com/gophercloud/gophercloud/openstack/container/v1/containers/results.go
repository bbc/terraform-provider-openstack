package containers

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/container/v1/capsules"
	"github.com/gophercloud/gophercloud/pagination"
)

type commonResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts
// a capsule resource.
func (r commonResult) Extract() (*Container, error) {
	var s *Container
	err := r.ExtractInto(&s)
	return s, err
}

// GetResult represents the result of a get operation.
type GetResult struct {
	commonResult
}

// CreateResult is the response from a Create operation. Call its Extract
// method to interpret it as a Container.
type CreateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation.
type DeleteResult struct {
	gophercloud.ErrResult
}

// StartResult represents the result of a start operation.
type StartResult struct {
	gophercloud.ErrResult
}

type ContainerPage struct {
	pagination.LinkedPageBase
}

type Container struct {
	capsules.Container

	// The entrypoint for the container
	Entrypoint []string `json:"entrypoint"`

	// Whether the container has extra privileges
	Privileged bool `json:"privileged"`
}

// ExtractContainers accepts a Page struct, specifically a ContainerPage struct,
// and extracts the elements into an Container.
func ExtractContainers(r pagination.Page) ([]Container, error) {
	var s struct {
		Containers []Container `json:"containers"`
	}

	err := (r.(ContainerPage)).ExtractInto(&s)
	return s.Containers, err
}

// NextPageURL is invoked when a paginated collection of containers has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r ContainerPage) NextPageURL() (string, error) {
	var s struct {
		Next string `json:"next"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return s.Next, nil
}

// IsEmpty checks whether a ContainerPage struct is empty.
func (r ContainerPage) IsEmpty() (bool, error) {
	v, err := ExtractContainers(r)
	if err != nil {
		return false, err
	}

	return len(v) == 0, nil
}
