package openstack

import (
	"fmt"
	"log"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/container/v1/extensions/attachinterfaces"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func containerInterfaceAttachV1AttachFunc(
	containerClient *gophercloud.ServiceClient, containerID, portID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		va, err := attachinterfaces.Get(containerClient, containerID, portID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return va, "ATTACHING", nil
			}
			return va, "", err
		}

		return va, "ATTACHED", nil
	}
}

func containerInterfaceAttachV1DetachFunc(
	containerClient *gophercloud.ServiceClient, containerID, portID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to detach openstack_container_interface_attach_v1 %s from container %s",
			portID, containerID)

		va, err := attachinterfaces.Get(containerClient, containerID, portID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return va, "DETACHED", nil
			}
			return va, "", err
		}

		err = attachinterfaces.Delete(containerClient, containerID, attachinterfaces.DeleteOpts{PortID: portID}).ExtractErr()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return va, "DETACHED", nil
			}

			if _, ok := err.(gophercloud.ErrDefault400); ok {
				return nil, "", nil
			}

			return nil, "", err
		}

		log.Printf("[DEBUG] openstack_compute_interface_attach_v2 %s is still active.", portID)
		return nil, "", nil
	}
}

func containerInterfaceAttachV1ParseID(id string) (string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) < 2 {
		return "", "", fmt.Errorf("Unable to determine openstack_container_interface_attach_v1 %s ID", id)
	}

	containerID := idParts[0]
	attachmentID := idParts[1]

	return containerID, attachmentID, nil
}
