package calico

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/projectcalico/libcalico-go/lib/errors"
	"github.com/projectcalico/libcalico-go/lib/options"
)

func dataSourceNodes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCalicoNodeList,
		Schema: map[string]*schema.Schema{
			"calico_nodes": {
				Description: "The list of calico nodes",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceCalicoNodeList(d *schema.ResourceData, m interface{}) error {
	calicoClient := m.(config).Client
	nodesInterface := calicoClient.Nodes()
	nodes, err := nodesInterface.List(ctx, options.ListOptions{})

	if err != nil {
		if _, ok := err.(errors.ErrorResourceDoesNotExist); ok {
			d.SetId("")
			return nil
		} else {
			return err
		}
	}
	var nodeList []string
	for _, node := range nodes.Items {
		nodeList = append(nodeList, node.Name)
	}
	d.SetId(nodes.ListMeta.ResourceVersion)
	d.Set("calico_nodes", nodeList)

	return nil
}
