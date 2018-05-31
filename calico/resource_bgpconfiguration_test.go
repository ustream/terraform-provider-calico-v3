package calico

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/projectcalico/libcalico-go/lib/options"
	"strings"
	"testing"
)

func TestAccBgpConfiguration(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBgpConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateBgpConfigurationConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBgpConfigurationExists("calico_bgpconfiguration.test"),
					resource.TestCheckResourceAttr("calico_bgpconfiguration.test", "metadata.0.name", "testbgpconfiguration"),
					resource.TestCheckResourceAttr("calico_bgpconfiguration.test", "spec.0.log_severity_screen", "Warning"),
					resource.TestCheckResourceAttr("calico_bgpconfiguration.test", "spec.0.node_to_node_mesh_enabled", "false"),
				),
			},
			{
				Config: testUpdateBgpConfigurationConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBgpConfigurationExists("calico_bgpconfiguration.test"),
					resource.TestCheckResourceAttr("calico_bgpconfiguration.test", "metadata.0.name", "testbgpconfiguration2"),
					resource.TestCheckResourceAttr("calico_bgpconfiguration.test", "spec.0.log_severity_screen", "Info"),
					resource.TestCheckResourceAttr("calico_bgpconfiguration.test", "spec.0.node_to_node_mesh_enabled", "true"),
				),
			},
		},
	})
}

func testAccCheckBgpConfigurationDestroy(state *terraform.State) error {

	client := testAccProvider.Meta().(config).Client

	apis := getResourcesByType("calico_BgpConfiguration", state)

	if len(apis) != 1 {
		return fmt.Errorf("expecting only 1 BgpConfiguration resource found %v", len(apis))
	}

	_, err := client.BGPConfigurations().Get(ctx, apis[0].Primary.ID, options.GetOptions{})

	switch {
	case err == nil:
		return errors.New("Expected error, got none")
	case err != nil && !strings.Contains(err.Error(), "resource does not exist"):
		return fmt.Errorf("Expected 404, got %s", err)
	}

	return nil
}

func testAccCheckBgpConfigurationExists(resourceKey string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceKey]

		if !ok {
			return fmt.Errorf("not found: %s", resourceKey)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		BgpConfiguration, err := testAccProvider.Meta().(config).Client.BGPConfigurations().Get(ctx, rs.Primary.ID, options.GetOptions{})

		if err != nil {
			return err
		}

		if BgpConfiguration == nil {
			return fmt.Errorf("BgpConfiguration with id %v not found", rs.Primary.ID)
		}

		return nil
	}
}

const testCreateBgpConfigurationConfig = `
resource "calico_BgpConfiguration" "test" {
  metadata{
    name = "testbgpconfiguration"
  }
  spec{
    log_severity_screen = "Warning"
    node_to_node_mesh_enabled = "false"
  }
}`

const testUpdateBgpConfigurationConfig = `
resource "calico_BgpConfiguration" "test" {
  metadata{
    name = "testbgpconfiguration2"
  }
  spec{
    log_severity_screen = "Info"
    node_to_node_mesh_enabled = "true"
  }
}`
