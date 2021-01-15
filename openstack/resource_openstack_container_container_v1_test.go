package openstack

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/container/v1/containers"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestContainerContainerV1_basic(t *testing.T) {
	var container containers.Container

	resourceName := "openstack_container_container_v1.container_1"
	imageName := acctest.RandomWithPrefix("tf-acc-image")

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerContainerV1Basic(imageName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContainerContainerV1Exists(resourceName, &container),
					resource.TestCheckResourceAttr(resourceName, "name", "container_1"),
				),
			},
		},
	})
}

func testAccContainerContainerV1Basic(imageName string) string {
	return fmt.Sprintf(`
resource "openstack_container_container_v1" "container_1" {
  name = "container_1"
  image = "%s"
}
`, imageName)
}

func testAccCheckContainerContainerV1Exists(n string, container *containers.Container) resource.TestCheckFunc {
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

		found, err := containers.Get(containerClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.UUID != rs.Primary.ID {
			return fmt.Errorf("Cluster not found")
		}

		*container = *found

		return nil
	}
}
