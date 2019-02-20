package calico

import (
	"github.com/hashicorp/terraform/helper/schema"
	api "github.com/projectcalico/libcalico-go/lib/apis/v3"
	"github.com/projectcalico/libcalico-go/lib/errors"
	"github.com/projectcalico/libcalico-go/lib/options"
)

func resourceHostEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceHostEndpointCreate,
		Read:   resourceHostEndpointRead,
		Update: resourceHostEndpointUpdate,
		Delete: resourceHostEndpointDelete,

		Schema: map[string]*schema.Schema{
			"metadata": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"labels": {
							Type:     schema.TypeMap,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"spec": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"interfacename": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},
					},
				},
			},
		},
	}
}

func resourceHostEndpointCreate(d *schema.ResourceData, m interface{}) error {
	calicoClient := m.(config).Client
	hostEndpointInterface := calicoClient.HostEndpoints()

	hostEndpoint, err := createHostEndpoint(d)
	if err != nil {
		return err
	}

	_, err = hostEndpointInterface.Create(ctx, hostEndpoint, opts)
	if err != nil {
		return err
	}

	d.SetId(hostEndpoint.ObjectMeta.Name)
	return resourceHostEndpointRead(d, m)
}

func resourceHostEndpointRead(d *schema.ResourceData, m interface{}) error {
	calicoClient := m.(config).Client
	hostEndpointInterface := calicoClient.HostEndpoints()

	hostEndpointname := dToString(d, "metadata.0.name")

	hostEndpoint, err := hostEndpointInterface.Get(ctx, hostEndpointname, options.GetOptions{})

	if err != nil {
		if _, ok := err.(errors.ErrorResourceDoesNotExist); ok {
			d.SetId("")
			return nil
		} else {
			return err
		}
	}

	d.SetId(hostEndpointname)
	d.Set("metadata.0.name", hostEndpoint.ObjectMeta.Name)

	return nil
}

func resourceHostEndpointUpdate(d *schema.ResourceData, m interface{}) error {
	d.Partial(false)

	calicoClient := m.(config).Client
	hostEndpointInterface := calicoClient.HostEndpoints()

	hostEndpointname := dToString(d, "metadata.0.name")
	backtrackhostEndpoint, err := hostEndpointInterface.Get(ctx, hostEndpointname, options.GetOptions{})
	if err != nil {
		if _, ok := err.(errors.ErrorResourceDoesNotExist); ok {
			d.SetId("")
			return nil
		} else {
			return err
		}
	}

	hostEndpoint, err := createHostEndpoint(d)
	if err != nil {
		return err
	}
	hostEndpoint.ObjectMeta = backtrackhostEndpoint.ObjectMeta

	_, err = hostEndpointInterface.Update(ctx, hostEndpoint, opts)
	if err != nil {
		return err
	}

	return nil
}

func resourceHostEndpointDelete(d *schema.ResourceData, m interface{}) error {
	calicoClient := m.(config).Client
	hostEndpointInterface := calicoClient.HostEndpoints()

	hostEndpointname := dToString(d, "metadata.0.name")

	_, err := hostEndpointInterface.Delete(ctx, hostEndpointname, options.DeleteOptions{})
	if err != nil {
		return err
	}

	return nil
}

func createHostEndpoint(d *schema.ResourceData) (*api.HostEndpoint, error) {
	hostEndpoint := api.NewHostEndpoint()

	if v, ok := d.GetOk("metadata.0.labels"); ok {
		labelsMap := v.(map[string]interface{})
		labels := make(map[string]string, len(labelsMap))
		for k, v := range labelsMap {
			labels[k] = v.(string)
		}
		hostEndpoint.ObjectMeta.Labels = labels
	}

	hostEndpoint.ObjectMeta.Name = dToString(d, "metadata.0.name")
	hostEndpoint.Spec.Node = dToString(d, "spec.0.node")
	hostEndpoint.Spec.InterfaceName = dToString(d, "spec.0.interfacename")
	return hostEndpoint, nil
}
