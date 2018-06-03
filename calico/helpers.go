package calico

import (
	"github.com/hashicorp/terraform/helper/schema"
	api "github.com/projectcalico/libcalico-go/lib/apis/v3"
	"github.com/projectcalico/libcalico-go/lib/numorstring"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Generate the metadata of the resource
func dToTypeMeta(d *schema.ResourceData) (meta.ObjectMeta, error) {
	objectMeta := meta.ObjectMeta{}

	objectMeta.Name = dToString(d, "metadata.0.name")

	return objectMeta, nil
}

// Convert field to IPIPMode type
func dToIpIpMode(d *schema.ResourceData, field string) api.IPIPMode {
	f := api.IPIPMode(dToString(d, field))
	return f
}

// Convert field to ASNumber type
func dToAsNumber(d *schema.ResourceData, field string) numorstring.ASNumber {
	f := numorstring.ASNumber(uint(d.Get(field).(int)))
	return f
}

func dToProtoPort(d *schema.ResourceData, field string) map[string]string {
	if v, ok := d.GetOk(field); ok {
		labelMap := v.(map[string]interface{})
		labels := make(map[string]string, len(labelMap))

		for k, v := range labelMap {
			labels[k] = v.(string)
		}
		Labels := labels
		return Labels
	} else {
		return nil
	}
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

// Convert field to map
func dToMap(d *schema.ResourceData, field string) map[string]interface{} {
	f := d.Get(field).(map[string]interface{})

	return f
}