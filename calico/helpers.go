package calico

import (
	"github.com/hashicorp/terraform/helper/schema"
	api "github.com/projectcalico/libcalico-go/lib/apis/v3"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/projectcalico/libcalico-go/lib/numorstring"
)

// Generate the metadata of the resource
func dToTypeMeta(d *schema.ResourceData) (meta.ObjectMeta, error) {
	objectMeta := meta.ObjectMeta{}

	objectMeta.Name = dToString(d,"metadata.0.name")

	return objectMeta, nil
}

// Convert field to IPIPMode type
func dToIpIpMode(d *schema.ResourceData, field string) api.IPIPMode {
	f := api.IPIPMode(d.Get(field).(string))
	return f
}

// Convert field to ASNumber type
func dToAsNumber(d *schema.ResourceData, field string) numorstring.ASNumber {
	f := numorstring.ASNumber(uint(d.Get(field).(int)))
	return f
}

// Convert field to string
func dToString(d *schema.ResourceData, field string) string {
	f := d.Get(field).(string)
	return f
}

// Convert field to bool
func dToBool(d *schema.ResourceData, field string) bool {
	f := d.Get(field).(bool)
	return f
}

