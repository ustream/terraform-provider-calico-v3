package main

import (
	"github.com/cdiscount/terraform-provider-calico/calico"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return calico.Provider()
		},
	})
}
