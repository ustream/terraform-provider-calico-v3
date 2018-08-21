package calico

import (
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	api "github.com/projectcalico/libcalico-go/lib/apis/v3"
	"github.com/projectcalico/libcalico-go/lib/errors"
	caliconet "github.com/projectcalico/libcalico-go/lib/net"
	"github.com/projectcalico/libcalico-go/lib/numorstring"
	"github.com/projectcalico/libcalico-go/lib/options"
)

func resourceGlobalNetworkPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceGlobalNetworkPolicyCreate,
		Read:   resourceGlobalNetworkPolicyRead,
		Update: resourceGlobalNetworkPolicyUpdate,
		Delete: resourceGlobalNetworkPolicyDelete,

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
							ForceNew: false,
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
						"order": {
							Type:     schema.TypeFloat,
							Optional: true,
							ForceNew: true,
						},
						"track": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"apply_on_forward": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"selector": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"prednat": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"ingress": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     ruleSchema(),
						},
						"egress": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     ruleSchema(),
						},
					},
				},
			},
		},
	}
}

func resourceGlobalNetworkPolicyCreate(d *schema.ResourceData, m interface{}) error {
	calicoClient := m.(config).Client
	globalNetworkPolicyInterface := calicoClient.GlobalNetworkPolicies()

	globalNetworkPolicy, err := createGlobalNetworkPolicy(d)
	if err != nil {
		return err
	}

	_, err = globalNetworkPolicyInterface.Create(ctx, globalNetworkPolicy, opts)
	if err != nil {
		return err
	}

	d.SetId(globalNetworkPolicy.ObjectMeta.Name)
	return resourceGlobalNetworkPolicyRead(d, m)
}

func resourceGlobalNetworkPolicyRead(d *schema.ResourceData, m interface{}) error {
	calicoClient := m.(config).Client
	globalNetworkPolicyInterface := calicoClient.GlobalNetworkPolicies()

	nameglobalNetworkPolicy := dToString(d, "metadata.0.name")

	globalNetworkPolicy, err := globalNetworkPolicyInterface.Get(ctx, nameglobalNetworkPolicy, options.GetOptions{})

	if err != nil {
		if _, ok := err.(errors.ErrorResourceDoesNotExist); ok {
			d.SetId("")
			return nil
		} else {
			return err
		}
	}

	d.SetId(nameglobalNetworkPolicy)
	d.Set("metadata.0.name", globalNetworkPolicy.ObjectMeta.Name)

	return nil
}

func resourceGlobalNetworkPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	d.Partial(false)

	calicoClient := m.(config).Client
	globalNetworkPolicyInterface := calicoClient.GlobalNetworkPolicies()

	nameglobalNetworkPolicy := dToString(d, "metadata.0.name")
	backtrackglobalNetworkPolicy, err := globalNetworkPolicyInterface.Get(ctx, nameglobalNetworkPolicy, options.GetOptions{})
	if err != nil {
		if _, ok := err.(errors.ErrorResourceDoesNotExist); ok {
			d.SetId("")
			return nil
		} else {
			return err
		}
	}

	globalNetworkPolicy, err := createGlobalNetworkPolicy(d)
	if err != nil {
		return err
	}
	globalNetworkPolicy.ObjectMeta = backtrackglobalNetworkPolicy.ObjectMeta

	_, err = globalNetworkPolicyInterface.Update(ctx, globalNetworkPolicy, opts)
	if err != nil {
		return err
	}

	return nil
}

func resourceGlobalNetworkPolicyDelete(d *schema.ResourceData, m interface{}) error {
	calicoClient := m.(config).Client
	globalNetworkPolicyInterface := calicoClient.GlobalNetworkPolicies()

	nameglobalNetworkPolicy := dToString(d, "metadata.0.name")

	_, err := globalNetworkPolicyInterface.Delete(ctx, nameglobalNetworkPolicy, options.DeleteOptions{})
	if err != nil {
		return err
	}

	return nil
}

func createGlobalNetworkPolicy(d *schema.ResourceData) (*api.GlobalNetworkPolicy, error) {
	spec, err := dToGlobalNetworkPolicySpec(d)
	if err != nil {
		return nil, err
	}

	objectMeta, err := dToTypeMeta(d)
	if err != nil {
		return nil, err
	}

	newGlobalNetworkPolicy := api.NewGlobalNetworkPolicy()
	newGlobalNetworkPolicy.ObjectMeta = objectMeta
	newGlobalNetworkPolicy.ObjectMeta.Name = dToString(d, "metadata.0.name")
	newGlobalNetworkPolicy.Spec = spec
	return newGlobalNetworkPolicy, nil
}

func dToGlobalNetworkPolicySpec(d *schema.ResourceData) (api.GlobalNetworkPolicySpec, error) {
	spec := api.GlobalNetworkPolicySpec{}

	spec.PreDNAT = dToBool(d, "spec.0.prednat")
	spec.Selector = dToString(d, "spec.0.selector")
	order := dToFloat64(d, "spec.0.order")
	spec.Order = &order
	spec.DoNotTrack = dToBool(d, "spec.0.track")
	spec.ApplyOnForward = dToBool(d, "spec.0.apply_on_forward")
	if v, ok := d.GetOk("spec.0.ingress.0.rule.#"); ok {
		ingressRules := make([]api.Rule, v.(int))
		for i := range ingressRules {
			ruleMap := d.Get("spec.0.ingress.0.rule." + strconv.Itoa(i)).(map[string]interface{})
			rule, err := ruleMapToRule(ruleMap)
			if err != nil {
				return spec, err
			}
			spec.Ingress = append(spec.Ingress, rule)
			ingressRules[i] = rule
		}
	}
	if v, ok := d.GetOk("spec.0.egress.0.rule.#"); ok {
		egressRules := make([]api.Rule, v.(int))
		for i := range egressRules {
			ruleMap := d.Get("spec.0.egress.0.rule." + strconv.Itoa(i)).(map[string]interface{})
			rule, err := ruleMapToRule(ruleMap)
			if err != nil {
				return spec, err
			}
			spec.Egress = append(spec.Egress, rule)
			egressRules[i] = rule
		}
	}
	if len(spec.Ingress) > 0 {
		spec.Types = append(spec.Types, "Ingress")
	}
	if len(spec.Egress) > 0 {
		spec.Types = append(spec.Types, "Egress")
	}

	return spec, nil
}

func ruleMapToRule(ruleMap map[string]interface{}) (api.Rule, error) {
	rule := api.Rule{}

	if val, ok := ruleMap["action"]; ok {
		switch val.(string) {
		case "allow":
			rule.Action = "Allow"
		case "Allow":
			rule.Action = "Allow"
		case "deny":
			rule.Action = "Deny"
		case "Deny":
			rule.Action = "Deny"
		case "log":
			rule.Action = "Log"
		case "Log":
			rule.Action = "Log"
		case "pass":
			rule.Action = "Pass"
		case "Pass":
			rule.Action = "Pass"

		}
	}
	if val, ok := ruleMap["protocol"]; ok {
		if len(val.(string)) > 0 {
			protocol := numorstring.ProtocolFromString(val.(string))
			rule.Protocol = &protocol
		}
	}
	if val, ok := ruleMap["not_protocol"]; ok {
		if len(val.(string)) > 0 {
			notProtocol := numorstring.ProtocolFromString(val.(string))
			rule.NotProtocol = &notProtocol
		}
	}
	if val, ok := ruleMap["icmp"]; ok {
		icmpList := val.([]interface{})
		if len(icmpList) > 0 {
			for _, v := range icmpList {
				icmpMap := v.(map[string]interface{})
				icmpType := icmpMap["type"].(int)
				icmpCode := icmpMap["code"].(int)
				icmp := api.ICMPFields{
					Type: &icmpType,
					Code: &icmpCode,
				}
				rule.ICMP = &icmp
			}
		}
	}
	if val, ok := ruleMap["not_icmp"]; ok {
		icmpList := val.([]interface{})
		if len(icmpList) > 0 {
			for _, v := range icmpList {
				icmpMap := v.(map[string]interface{})
				icmpType := icmpMap["type"].(int)
				icmpCode := icmpMap["code"].(int)
				icmp := api.ICMPFields{
					Type: &icmpType,
					Code: &icmpCode,
				}
				rule.NotICMP = &icmp
			}
		}
	}
	if val, ok := ruleMap["source"]; ok {
		sourceList := val.([]interface{})

		if len(sourceList) > 0 {
			srcEntityRules, err := srcDstListToEntityRule(sourceList)
			if err != nil {
				return rule, err
			}
			rule.Source = srcEntityRules
		}
	}

	if val, ok := ruleMap["destination"]; ok {
		destinationList := val.([]interface{})

		if len(destinationList) > 0 {
			destEntityRules, err := srcDstListToEntityRule(destinationList)
			if err != nil {
				return rule, err
			}
			rule.Destination = destEntityRules
		}
	}

	return rule, nil
}

func srcDstListToEntityRule(srcDstList []interface{}) (api.EntityRule, error) {
	entityRule := api.EntityRule{}
	resourceRuleMap := srcDstList[0].(map[string]interface{})

	for _, n := range resourceRuleMap["nets"].([]interface{}) {
		_, _, err := caliconet.ParseCIDR(n.(string))
		if err != nil {
			return entityRule, err
		}
		entityRule.Nets = append(entityRule.Nets, n.(string))
	}
	for _, n := range resourceRuleMap["not_nets"].([]interface{}) {
		_, _, err := caliconet.ParseCIDR(n.(string))
		if err != nil {
			return entityRule, err
		}
		entityRule.NotNets = append(entityRule.NotNets, n.(string))
	}
	if v, ok := resourceRuleMap["selector"]; ok {
		entityRule.Selector = v.(string)
	}
	if v, ok := resourceRuleMap["not_selector"]; ok {
		entityRule.NotSelector = v.(string)
	}
	if v, ok := resourceRuleMap["ports"]; ok {
		if resourcePortList, ok := v.([]interface{}); ok {
			portList, err := toPortList(resourcePortList)
			if err != nil {
				return entityRule, err
			}
			entityRule.Ports = portList
		}
	}
	if v, ok := resourceRuleMap["not_ports"]; ok {
		if resourcePortList, ok := v.([]interface{}); ok {
			portList, err := toPortList(resourcePortList)
			if err != nil {
				return entityRule, err
			}
			entityRule.NotPorts = portList
		}
	}
	return entityRule, nil
}

func toPortList(resourcePortList []interface{}) ([]numorstring.Port, error) {
	portList := make([]numorstring.Port, len(resourcePortList))

	for i, v := range resourcePortList {
		p, err := numorstring.PortFromString(v.(string))
		if err != nil {
			return portList, err
		}
		portList[i] = p
	}
	return portList, nil
}
