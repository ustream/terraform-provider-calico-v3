package calico

import (
	"time"

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

// Convert field to int
func dToInt(d *schema.ResourceData, field string) int {
	f := d.Get(field).(int)
	return f
}

// Convert field to uint32
func dToUint32(d *schema.ResourceData, field string) uint32 {
	f := uint32(d.Get(field).(uint))
	return f
}

// Convert field to float64
func dToFloat64(d *schema.ResourceData, field string) float64 {
	f := float64(d.Get(field).(float64))
	return f
}

// Convert field to duration
func dToDuration(d *schema.ResourceData, field string) meta.Duration {
	f := time.Duration(dToInt(d, field))
	f2 := meta.Duration{f}

	return f2
}

func entityRuleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"nets": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"not_nets": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"selector": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"not_selector": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"ports": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"not_ports": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func ruleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"rule": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"protocol": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"not_protocol": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"icmp": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:     schema.TypeInt,
										Optional: true,
									},
									"code": &schema.Schema{
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
						"not_icmp": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:     schema.TypeInt,
										Optional: true,
									},
									"code": &schema.Schema{
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
						"source": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem:     entityRuleSchema(),
						},
						"destination": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem:     entityRuleSchema(),
						},
					},
				},
			},
		},
	}
}
