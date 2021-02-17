package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/gophercloud/gophercloud/openstack/container/v1/extensions/attachinterfaces"
)

func TestAccContainerV1InterfaceAttach_basic(t *testing.T) {
	var ai attachinterfaces.Interface

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerV1InterfaceAttachDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerV1InterfaceAttachBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerV1InterfaceAttachExists("openstack_container_interface_attach_v1.ai_1", &ai),
				),
			},
		},
	})
}

func TestAccContainerV1InterfaceAttach_IP(t *testing.T) {
	var ai attachinterfaces.Interface

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerV1InterfaceAttachDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerV1InterfaceAttachIP(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerV1InterfaceAttachExists("openstack_container_interface_attach_v1.ai_1", &ai),
					testAccCheckContainerV1InterfaceAttachIP(&ai, "192.168.1.100"),
				),
			},
		},
	})
}

func testAccCheckContainerV1InterfaceAttachDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	containerClient, err := config.ContainerV1Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack container client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_container_interface_attach_v1" {
			continue
		}

		instanceID, portID, err := containerInterfaceAttachV1ParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = attachinterfaces.Get(containerClient, instanceID, portID).Extract()
		if err == nil {
			return fmt.Errorf("Volume attachment still exists")
		}
	}

	return nil
}

func testAccCheckContainerV1InterfaceAttachExists(n string, ai *attachinterfaces.Interface) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		containerClient, err := config.ContainerV1Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack container client: %s", err)
		}

		instanceID, portID, err := containerInterfaceAttachV1ParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		found, err := attachinterfaces.Get(containerClient, instanceID, portID).Extract()
		if err != nil {
			return err
		}

		//if found.instanceID != instanceID || found.PortID != portID {
		if found.PortID != portID {
			return fmt.Errorf("InterfaceAttach not found")
		}

		*ai = *found

		return nil
	}
}

func testAccCheckContainerV1InterfaceAttachIP(
	ai *attachinterfaces.Interface, ip string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, i := range ai.FixedIPs {
			if i.IPAddress == ip {
				return nil
			}
		}
		return fmt.Errorf("Requested ip (%s) does not exist on port", ip)
	}
}

func testAccContainerV1InterfaceAttachBasic() string {
	return fmt.Sprintf(`
resource "openstack_networking_port_v2" "port_1" {
  name = "port_1"
  network_id = "%s"
  admin_state_up = "true"
}

resource "openstack_container_container_v1" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
}

resource "openstack_container_interface_attach_v1" "ai_1" {
  instance_id = "${openstack_container_container_v1.instance_1.id}"
  port_id = "${openstack_networking_port_v2.port_1.id}"
}
`, osNetworkID, osNetworkID)
}

func testAccContainerV1InterfaceAttachIP() string {
	return fmt.Sprintf(`
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  network_id = "${openstack_networking_network_v2.network_1.id}"
  cidr = "192.168.1.0/24"
  ip_version = 4
  enable_dhcp = true
  no_gateway = true
}

resource "openstack_container_container_v1" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
}

resource "openstack_container_interface_attach_v1" "ai_1" {
  instance_id = "${openstack_container_container_v1.instance_1.id}"
  network_id = "${openstack_networking_network_v2.network_1.id}"
  fixed_ip = "192.168.1.100"
}
`, osNetworkID)
}
