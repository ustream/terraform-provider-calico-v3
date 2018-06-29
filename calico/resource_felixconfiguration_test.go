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

func TestAccFelixConfiguration(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFelixConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateFelixConfigurationConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFelixConfigurationExists("calico_felixconfiguration.test"),
					resource.TestCheckResourceAttr("calico_felixconfiguration.test", "metadata.0.name", "testfelixconfiguration"),
					resource.TestCheckResourceAttr("calico_felixconfiguration.test", "spec.0.chain_insert_mode", "Insert"),
					resource.TestCheckResourceAttr("calico_felixconfiguration.test", "spec.0.default_endpoint_to_host_action", "Drop"),
					resource.TestCheckResourceAttr("calico_felixconfiguration.test", "spec.0.ignore_loose_rpf", "true"),
					resource.TestCheckResourceAttr("calico_felixconfiguration.test", "spec.0.interface_exclude", "kube-ipvs1"),
					resource.TestCheckResourceAttr("calico_felixconfiguration.test", "spec.0.interface_prefix", "calicotest"),
				),
			},
			{
				Config: testUpdateFelixConfigurationConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFelixConfigurationExists("calico_felixconfiguration.test"),
					resource.TestCheckResourceAttr("calico_felixconfiguration.test", "metadata.0.name", "testfelixconfiguration2"),
					resource.TestCheckResourceAttr("calico_felixconfiguration.test", "spec.0.chain_insert_mode", "Append"),
					resource.TestCheckResourceAttr("calico_felixconfiguration.test", "spec.0.default_endpoint_to_host_action", "Return"),
					resource.TestCheckResourceAttr("calico_felixconfiguration.test", "spec.0.ignore_loose_rpf", "false"),
					resource.TestCheckResourceAttr("calico_felixconfiguration.test", "spec.0.interface_exclude", "kube-ipvs2"),
					resource.TestCheckResourceAttr("calico_felixconfiguration.test", "spec.0.interface_prefix", "calico"),
				),
			},
		},
	})
}

func testAccCheckFelixConfigurationDestroy(state *terraform.State) error {

	client := testAccProvider.Meta().(config).Client

	apis := getResourcesByType("calico_felixconfiguration", state)

	if len(apis) != 1 {
		return fmt.Errorf("expecting only 1 FelixConfiguration resource found %v", len(apis))
	}

	_, err := client.FelixConfigurations().Get(ctx, apis[0].Primary.ID, options.GetOptions{})

	switch {
	case err == nil:
		return errors.New("Expected error, got none")
	case err != nil && !strings.Contains(err.Error(), "resource does not exist"):
		return fmt.Errorf("Expected 404, got %s", err)
	}

	return nil
}

func testAccCheckFelixConfigurationExists(resourceKey string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceKey]

		if !ok {
			return fmt.Errorf("not found: %s", resourceKey)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		FelixConfiguration, err := testAccProvider.Meta().(config).Client.FelixConfigurations().Get(ctx, rs.Primary.ID, options.GetOptions{})

		if err != nil {
			return err
		}

		if FelixConfiguration == nil {
			return fmt.Errorf("FelixConfiguration with id %v not found", rs.Primary.ID)
		}

		return nil
	}
}

const testCreateFelixConfigurationConfig = `
resource "calico_felixconfiguration" "test" {
  metadata{
    name = "testfelixconfiguration"
  }
  spec{
    chain_insert_mode = "Insert"
	default_endpoint_to_host_action = "Drop"
	ignore_loose_rpf = true
	interface_exclude = "kube-ipvs1"
	interface_prefix = "calicotest"
  }
}`

const testUpdateFelixConfigurationConfig = `
resource "calico_felixconfiguration" "test" {
  metadata{
    name = "testfelixconfiguration2"
  }
  spec{
    chain_insert_mode = "Append"
	default_endpoint_to_host_action = "Return"
	ignore_loose_rpf = false
	interface_exclude = "kube-ipvs2"
	interface_prefix = "calico"
  }
}`
