package calico

import (
	"github.com/hashicorp/terraform/helper/schema"
	api "github.com/projectcalico/libcalico-go/lib/apis/v3"
	"github.com/projectcalico/libcalico-go/lib/errors"
	"github.com/projectcalico/libcalico-go/lib/options"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

func resourceCalicoBgpConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceCalicoBgpConfigurationCreate,
		Read:   resourceCalicoBgpConfigurationRead,
		Update: resourceCalicoBgpConfigurationUpdate,
		Delete: resourceCalicoBgpConfigurationDelete,

		Schema: map[string]*schema.Schema{
			"metadata": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ForceNew: false,
						},
					},
				},
			},
			"spec": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"log_severity_screen": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"node_to_node_mesh_enabled": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
						"as_number": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

// dToBgpConfigurationSpec return the spec of the BgpConfiguration
func dToBgpConfigurationSpec(d *schema.ResourceData) (api.BGPConfigurationSpec, error) {
	spec := api.BGPConfigurationSpec{}

	logSeverityScreen := d.Get("spec.0.log_severity_screen").(string)
	spec.LogSeverityScreen = logSeverityScreen

	nodeToNodeMeshEnabled := d.Get("spec.0.node_to_node_mesh_enabled").(bool)
	spec.NodeToNodeMeshEnabled = &nodeToNodeMeshEnabled

	//TODO: Reactivate this field
	//asNumber := d.Get("spec.0.as_number").(numorstring.ASNumber)
	//spec.ASNumber = &asNumber

	return spec, nil
}

// dToBgpConfigurationSpec return the metadata of the BgpConfiguration
func dToBgpConfigurationTypeMeta(d *schema.ResourceData) (meta.ObjectMeta, error) {
	objectMeta := meta.ObjectMeta{}

	objectMeta.Name = d.Get("metadata.0.name").(string)

	return objectMeta, nil
}

// resourceCalicoBgpConfigurationCreate create a new BgpConfiguration in Calico
func resourceCalicoBgpConfigurationCreate(d *schema.ResourceData, m interface{}) error {
	calicoClient := m.(config).Client
	BgpConfigurationInterface := calicoClient.BGPConfigurations()

	BgpConfiguration, err := createBgpConfigurationApiRequest(d)
	if err != nil {
		return err
	}

	_, err = BgpConfigurationInterface.Create(ctx, BgpConfiguration, opts)
	if err != nil {
		return err
	}

	d.SetId(BgpConfiguration.ObjectMeta.Name)
	return resourceCalicoBgpConfigurationRead(d, m)
}

// resourceCalicoBgpConfigurationRead get a specific BgpConfiguration
func resourceCalicoBgpConfigurationRead(d *schema.ResourceData, m interface{}) error {
	calicoClient := m.(config).Client
	BgpConfigurationInterface := calicoClient.BGPConfigurations()

	nameBgpConfiguration := d.Get("metadata.0.name").(string)

	BgpConfiguration, err := BgpConfigurationInterface.Get(ctx, nameBgpConfiguration, options.GetOptions{})
	log.Printf("Obj: %+v", d)

	// Handle endpoint does not exist
	if err != nil {
		if _, ok := err.(errors.ErrorResourceDoesNotExist); ok {
			d.SetId("")
			return nil
		} else {
			return err
		}
	}

	d.SetId(nameBgpConfiguration)
	d.Set("metadata.0.name", BgpConfiguration.ObjectMeta.Name)
	d.Set("spec.0.log_severity_screen", BgpConfiguration.Spec.LogSeverityScreen)
	d.Set("spec.0.node_to_node_mesh_enabled", BgpConfiguration.Spec.NodeToNodeMeshEnabled)
	d.Set("spec.0.as_number", BgpConfiguration.Spec.ASNumber)

	return nil
}

// resourceCalicoBgpConfigurationUpdate update an BgpConfiguration in Calico
func resourceCalicoBgpConfigurationUpdate(d *schema.ResourceData, m interface{}) error {
	d.Partial(false)

	calicoClient := m.(config).Client
	BgpConfigurationInterface := calicoClient.BGPConfigurations()

	BgpConfiguration, err := createBgpConfigurationApiRequest(d)
	if err != nil {
		return err
	}

	_, err = BgpConfigurationInterface.Update(ctx, BgpConfiguration, opts)
	if err != nil {
		return err
	}

	return nil
}

// resourceCalicoBgpConfigurationDelete delete an BgpConfiguration in Calico
func resourceCalicoBgpConfigurationDelete(d *schema.ResourceData, m interface{}) error {
	calicoClient := m.(config).Client
	BgpConfigurationInterface := calicoClient.BGPConfigurations()

	nameBgpConfiguration := d.Get("metadata.0.name").(string)

	_, err := BgpConfigurationInterface.Delete(ctx, nameBgpConfiguration, options.DeleteOptions{})
	if err != nil {
		return err
	}

	return nil
}

// createApiRequest prepare the request of creation and update
func createBgpConfigurationApiRequest(d *schema.ResourceData) (*api.BGPConfiguration, error) {
	// Set Spec to BgpConfiguration Spec
	spec, err := dToBgpConfigurationSpec(d)
	if err != nil {
		return nil, err
	}

	// Set Metadata to Kubernetes Metadata
	objectMeta, err := dToBgpConfigurationTypeMeta(d)
	if err != nil {
		return nil, err
	}

	// Create a new BGP Configuration, with TypeMeta filled in
	// Then, fill the metadata and the spec
	newBgpConfiguration := api.NewBGPConfiguration()
	newBgpConfiguration.ObjectMeta = objectMeta
	newBgpConfiguration.Spec = spec

	return newBgpConfiguration, nil
}
