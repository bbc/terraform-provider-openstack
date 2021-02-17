package openstack

import (
	"fmt"
	"log"
	"time"

	"github.com/gophercloud/gophercloud/openstack/container/v1/extensions/attachinterfaces"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceContainerInterfaceAttachV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceContainerInterfaceAttachV1Create,
		Read:   resourceContainerInterfaceAttachV1Read,
		Delete: resourceContainerInterfaceAttachV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"port_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"network_id"},
			},

			"network_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"port_id"},
			},

			"container_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"fixed_ip": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"port_id"},
			},
		},
	}
}

func resourceContainerInterfaceAttachV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	containerClient, err := config.ContainerV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack container client: %s", err)
	}

	containerID := d.Get("container_id").(string)

	var portID string
	if v, ok := d.GetOk("port_id"); ok {
		portID = v.(string)
	}

	var networkID string
	if v, ok := d.GetOk("network_id"); ok {
		networkID = v.(string)
	}

	if networkID == "" && portID == "" {
		return fmt.Errorf("Must set one of network_id and port_id")
	}

	var fixedIP string
	if v, ok := d.GetOk("fixed_ip"); ok {
		fixedIP = v.(string)
	}

	attachOpts := attachinterfaces.CreateOpts{
		PortID:    portID,
		NetworkID: networkID,
		FixedIP:   fixedIP,
	}

	log.Printf("[DEBUG] openstack_container_interface_attach_v2 attach options: %#v", attachOpts)

	err = attachinterfaces.Create(containerClient, containerID, attachOpts).ExtractErr()
	if err != nil {
		return err
	}
	attachment, err := attachinterfaces.Get(containerClient, containerID, portID).Extract()
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ATTACHING"},
		Target:     []string{"ATTACHED"},
		Refresh:    containerInterfaceAttachV1AttachFunc(containerClient, containerID, portID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	if _, err = stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error creating openstack_container_interface_attach_v1 %s: %s", containerID, err)
	}

	// Use the container ID and port ID as the resource ID.
	id := fmt.Sprintf("%s/%s", containerID, portID)

	log.Printf("[DEBUG] Created openstack_container_interface_attach_v1 %s: %#v", id, attachment)

	d.SetId(id)

	return resourceContainerInterfaceAttachV1Read(d, meta)
}

func resourceContainerInterfaceAttachV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	containerClient, err := config.ContainerV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack container client: %s", err)
	}

	containerID, attachmentID, err := containerInterfaceAttachV1ParseID(d.Id())
	if err != nil {
		return err
	}

	attachment, err := attachinterfaces.Get(containerClient, containerID, attachmentID).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Error retrieving openstack_container_interface_attach_v1")
	}

	log.Printf("[DEBUG] Retrieved openstack_container_interface_attach_v1 %s: %#v", d.Id(), attachment)

	d.Set("container_id", containerID)
	d.Set("port_id", attachment.PortID)
	d.Set("network_id", attachment.NetID)
	d.Set("fixed_ip", attachment.FixedIP)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceContainerInterfaceAttachV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	containerClient, err := config.ContainerV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack container client: %s", err)
	}

	containerID, attachmentID, err := containerInterfaceAttachV1ParseID(d.Id())
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{""},
		Target:     []string{"DETACHED"},
		Refresh:    containerInterfaceAttachV1DetachFunc(containerClient, containerID, attachmentID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	if _, err = stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error detaching openstack_container_interface_attach_v1 %s: %s", d.Id(), err)
	}

	return nil
}
