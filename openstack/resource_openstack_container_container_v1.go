package openstack

import (
	"fmt"
	"log"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/container/v1/containers"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceContainerContainerV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceContainerContainerV1Create,
		Read:   resourceContainerContainerV1Read,
		Delete: resourceContainerContainerV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"addresses": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeMap,
					},
				},
				Computed: true,
			},

			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"links": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				Computed: true,
			},

			"disk": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status_detail": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"task_state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ports": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Computed: true,
			},

			"cpu_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"image": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"command": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: true,
			},

			"cpu": {
				Type:     schema.TypeFloat,
				Optional: true,
				ForceNew: true,
			},

			"memory": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"workdir": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"labels": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: true,
			},

			"environment": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: true,
			},

			"restart_policy": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: true,
			},

			"interactive": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"tty": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"image_driver": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"security_groups": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: true,
			},

			"nets": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				Optional: true,
				ForceNew: true,
			},

			"runtime": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"hostname": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"auto_remove": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"auto_heal": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"hints": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: true,
			},

			"mounts": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				Optional: true,
				ForceNew: true,
			},

			"privileged": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"healthcheck": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: true,
			},

			"exposed_ports": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				Optional: true,
				ForceNew: true,
			},

			"host": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"entrypoint": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: true,
			},

			"registry": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"registry_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func (c *Config) ContainerV1Client(region string) (*gophercloud.ServiceClient, error) {
	return c.CommonServiceClientInit(openstack.NewContainerV1, region, "container")
}

func resourceContainerContainerV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	containerInfraClient, err := config.ContainerV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack container client: %s", err)
	}

	containerInfraClient.Microversion = "1.31"

	// Get boolean parameters that will be passed by reference.
	interactive := d.Get("interactive").(bool)
	tty := d.Get("tty").(bool)
	autoRemove := d.Get("auto_remove").(bool)
	autoHeal := d.Get("auto_heal").(bool)
	//privileged := d.Get("privileged").(bool)

	createOpts := containers.CreateOptsV131{
		CreateOpts: containers.CreateOpts{
			Name:             d.Get("name").(string),
			Image:            d.Get("image").(string),
			Command:          unfoldListOfStrings(d.Get("command").([]interface{})),
			CPU:              d.Get("cpu").(float64),
			Memory:           d.Get("memory").(string),
			Workdir:          d.Get("workdir").(string),
			Labels:           unfoldMapToString(d.Get("labels").(map[string]interface{})),
			Environment:      unfoldMapToString(d.Get("environment").(map[string]interface{})),
			RestartPolicy:    unfoldMapToString(d.Get("restart_policy").(map[string]interface{})),
			Interactive:      &interactive,
			TTY:              &tty,
			ImageDriver:      d.Get("image_driver").(string),
			SecurityGroups:   unfoldListOfStrings(d.Get("security_groups").([]interface{})),
			Nets:             unfoldListOfMapToString(d.Get("nets").([]interface{})),
			Runtime:          d.Get("runtime").(string),
			Hostname:         d.Get("hostname").(string),
			AutoRemove:       &autoRemove,
			AutoHeal:         &autoHeal,
			AvailabilityZone: d.Get("availability_zone").(string),
			Hints:            unfoldMapToString(d.Get("hints").(map[string]interface{})),
			Mounts:           unfoldListOfMapToString(d.Get("mounts").([]interface{})),
			//Privileged:       &privileged,
			Healthcheck:  unfoldMapToString(d.Get("healthcheck").(map[string]interface{})),
			ExposedPorts: unfoldMapToMapToString(d.Get("exposed_ports").(map[string]interface{})),
			Host:         d.Get("host").(string),
			Entrypoint:   unfoldListOfStrings(d.Get("entrypoint").([]interface{})),
		},
		Registry: d.Get("registry").(string),
	}

	log.Printf("[DEBUG] openstack_container_container_v1 create options: %#v", createOpts)

	c, err := containers.Create(containerInfraClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating openstack_container_container_v1: %s", err)
	}

	d.SetId(c.UUID)

	log.Printf("[DEBUG] Created openstack_container_container_v1 %s: %#v", c.UUID, c)
	return resourceContainerContainerV1Read(d, meta)
}

func resourceContainerContainerV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	containerInfraClient, err := config.ContainerV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack container client: %s", err)
	}

	c, err := containers.Get(containerInfraClient, d.Id()).Extract()
	if err != nil {
		return fmt.Errorf("Error retrieving openstack_container_container_v1: %s", err)
	}

	log.Printf("[DEBUG] Retrieved openstack_container_container_v1 %s: %#v", d.Id(), c)

	// Fill out the addresses from the returned data
	addresses := map[string][]map[string]interface{}{}
	for key, net := range c.Addresses {
		addrs := []map[string]interface{}{}
		for _, addr := range net {
			addrs = append(addrs, map[string]interface{}{
				"preserve_on_delete": addr.PreserveOnDelete,
				"addr":               addr.Addr,
				"port":               addr.Port,
				"version":            addr.Version,
				"subnet_id":          addr.SubnetID,
			})
		}
		addresses[key] = addrs
	}

	d.Set("links", c.Links)
	d.Set("addresses", addresses)
	d.Set("name", c.Name)
	d.Set("image", c.Image)
	d.Set("labels", c.Labels)
	d.Set("image_driver", c.ImageDriver)
	d.Set("security_groups", c.SecurityGroups)
	d.Set("command", c.Command)
	d.Set("cpu", c.CPU)
	d.Set("memory", c.Memory)
	d.Set("workdir", c.WorkDir)
	d.Set("environment", c.Environment)
	d.Set("restart_policy", c.RestartPolicy)
	d.Set("interactive", c.Interactive)
	d.Set("tty", c.TTY)
	d.Set("hostname", c.HostName)
	d.Set("status", c.Status)
	d.Set("status_detail", c.StatusDetail)
	d.Set("host", c.Host)
	d.Set("task_state", c.TaskState)
	d.Set("status_reason", c.StatusReason)
	d.Set("ports", c.Ports)
	d.Set("privileged", c.Privileged)
	d.Set("healthcheck", c.Healthcheck)
	d.Set("user_id", c.UserID)
	d.Set("project_id", c.ProjectID)
	d.Set("disk", c.Disk)
	d.Set("registry_id", c.RegistryID)
	d.Set("cpu_policy", c.CPUPolicy)
	d.Set("entrypoint", c.Entrypoint)

	return nil
}

func resourceContainerContainerV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	containerClient, err := config.ContainerV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack container client: %s", err)
	}

	if err := containers.Delete(containerClient, d.Id()).ExtractErr(); err != nil {
		return CheckDeleted(d, err, "Error deleting openstack_container_container_v1")
	}

	return nil
}

func unfoldListOfStrings(l []interface{}) []string {
	r := make([]string, len(l))
	for i, v := range l {
		r[i] = v.(string)
	}
	return r
}

func unfoldListOfMapToString(l []interface{}) []map[string]string {
	r := make([]map[string]string, len(l))
	for i, v := range l {
		r[i] = unfoldMapToString(v.(map[string]interface{}))
	}
	return r
}

func unfoldMapToMapToString(m map[string]interface{}) map[string]map[string]string {
	r := make(map[string]map[string]string, len(m))
	for k, v := range m {
		r[k] = unfoldMapToString(v.(map[string]interface{}))
	}
	return r
}

func unfoldMapToString(m map[string]interface{}) map[string]string {
	r := make(map[string]string, len(m))
	for k, v := range m {
		r[k] = v.(string)
	}
	return r
}
