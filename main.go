package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
	calicov3 "github.com/ustream/terraform-provider-calico-v3/calico"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return calicov3.Provider()
		},
	})
}
