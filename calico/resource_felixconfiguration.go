package calico

import (
	"github.com/hashicorp/terraform/helper/schema"
	api "github.com/projectcalico/libcalico-go/lib/apis/v3"
	"github.com/projectcalico/libcalico-go/lib/errors"
	"github.com/projectcalico/libcalico-go/lib/options"
)

func resourceCalicoFelixConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceCalicoFelixConfigurationCreate,
		Read:   resourceCalicoFelixConfigurationRead,
		Update: resourceCalicoFelixConfigurationUpdate,
		Delete: resourceCalicoFelixConfigurationDelete,

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
					},
				},
			},
			"spec": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"chain_insert_mode": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "Insert",
						},
						"default_endpoint_to_host_action": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "Drop",
						},
						"failsafe_inbound_host_ports": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"port": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"protocol": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "tcp",
									},
								},
							},
						},
						"failsafe_outbound_host_ports": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"port": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"protocol": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "tcp",
									},
								},
							},
						},
						"ignore_loose_rpf": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"interface_exclude": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "kube-ipvs0",
						},
						"interface_prefix": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "cali",
						},
						"ipip_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"ipip_mtu": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1440,
						},
						"ipsets_refresh_interval": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  10,
						},
						"iptables_filter_allow_action": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "Accept",
						},
						"iptables_lock_file_path": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "/run/xtables.lock",
						},
						"iptables_lock_probe_interval": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  50,
						},
						"iptables_lock_timeout": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						"iptables_mangle_allow_action": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "Accept",
						},
						"iptables_mark_mask": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "0xff000000",
						},
						"iptables_post_write_check_interval": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},
						"iptables_refresh_interval": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  90,
						},
						"ipv6_support": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"log_file_path": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "/var/log/calico/felix.log",
						},
						"log_prefix": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "calico-packet",
						},
						"log_severity_file": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "Info",
						},
						"log_severity_screen": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "Info",
						},
						"log_severity_sys": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "Info",
						},
						"max_ipset_size": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1048576,
						},
						"metadata_addr": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "127.0.0.1",
						},
						"metadata_port": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  8775,
						},
						"prometheus_go_metrics_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"prometheus_metrics_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"prometheus_metrics_port": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  9091,
						},
						"prometheus_process_metrics_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"reporting_interval": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  30,
						},
						"reporting_ttl": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  90,
						},
						"reporting_refresh_interval": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  90,
						},
						"usage_reporting_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"usage_reporting_initial_delay": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  300,
						},
						"usage_reporting_interval": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  86400,
						},
					},
				},
			},
		},
	}
}

// dToFelixConfigurationSpec return the spec of the FelixConfiguration
func dToFelixConfigurationSpec(d *schema.ResourceData) (api.FelixConfigurationSpec, error) {
	//TODO spec.FailsafeInboundHostPorts and spec.FailsafeOutboundHostPorts
	ignoreLooseRPF := dToBool(d, "spec.0.ignore_loose_rpf")
	ipipEnabled := dToBool(d, "spec.0.ipip_enabled")
	ipipMTU := dToInt(d, "spec.0.ipip_mtu")
	ipsetsRefreshInterval := dToDuration(d, "spec.0.ipsets_refresh_interval")
	iptablesLockProbeInterval := dToDuration(d, "spec.0.iptables_lock_probe_interval")
	iptablesLockTimeout := dToDuration(d, "spec.0.iptables_lock_timeout")
	iptablesMarkMask := dToUint32(d, "spec.0.iptables_mark_mask")
	iptablesPostWriteCheckInterval := dToDuration(d, "spec.0.iptables_post_write_check_interval")
	iptablesRefreshInterval := dToDuration(d, "spec.0.iptables_refresh_interval")
	ipv6Support := dToBool(d, "spec.0.ipv6_support")
	maxIpsetSize := dToInt(d, "spec.0.max_ipset_size")
	metadataPort := dToInt(d, "spec.0.metadata_port")
	prometheusGoMetricsEnabled := dToBool(d, "spec.0.prometheus_go_metrics_enabled")
	prometheusMetricsEnabled := dToBool(d, "spec.0.prometheus_metrics_enabled")
	prometheusMetricsPort := dToInt(d, "spec.0.prometheus_metrics_port")
	prometheusProcessMetricsEnabled := dToBool(d, "spec.0.prometheus_process_metrics_enabled")
	reportingInterval := dToDuration(d, "spec.0.reporting_interval")
	reportingTTL := dToDuration(d, "spec.0.reporting_ttl")
	routeRefreshInterval := dToDuration(d, "spec.0.route_refresh_interval")
	usageReportingEnabled := dToBool(d, "spec.0.usage_reporting_enabled")
	usageReportingInitialDelay := dToDuration(d, "spec.0.usage_reporting_initial_delay")
	usageReportingInterval := dToDuration(d, "spec.0.usage_reporting_interval")

	spec := api.FelixConfigurationSpec{}
	spec.ChainInsertMode = dToString(d, "spec.0.chain_insert_mode")
	spec.DefaultEndpointToHostAction = dToString(d, "spec.0.default_endpoint_to_host_action")
	spec.IgnoreLooseRPF = &ignoreLooseRPF
	spec.InterfaceExclude = dToString(d, "spec.0.interface_exclude")
	spec.InterfacePrefix = dToString(d, "spec.0.interface_prefix")
	spec.IPIPEnabled = &ipipEnabled
	spec.IPIPMTU = &ipipMTU
	spec.IpsetsRefreshInterval = &ipsetsRefreshInterval
	spec.IptablesFilterAllowAction = dToString(d, "spec.0.iptables_filter_allow_action")
	spec.IptablesLockFilePath = dToString(d, "spec.0.iptables_lock_file_path")
	spec.IptablesLockProbeInterval = &iptablesLockProbeInterval
	spec.IptablesLockTimeout = &iptablesLockTimeout
	spec.IptablesMangleAllowAction = dToString(d, "spec.0.iptables_mangle_allow_action")
	spec.IptablesMarkMask = &iptablesMarkMask
	spec.IptablesPostWriteCheckInterval = &iptablesPostWriteCheckInterval
	spec.IptablesRefreshInterval = &iptablesRefreshInterval
	spec.IPv6Support = &ipv6Support
	spec.LogFilePath = dToString(d, "spec.0.log_file_path")
	spec.LogPrefix = dToString(d, "spec.0.log_prefix")
	spec.LogSeverityFile = dToString(d, "spec.0.log_severity_file")
	spec.LogSeverityScreen = dToString(d, "spec.0.log_severity_screen")
	spec.LogSeveritySys = dToString(d, "spec.0.log_severity_sys")
	spec.MaxIpsetSize = &maxIpsetSize
	spec.MetadataAddr = dToString(d, "spec.0.metadata_addr")
	spec.MetadataPort = &metadataPort
	spec.PrometheusGoMetricsEnabled = &prometheusGoMetricsEnabled
	spec.PrometheusMetricsEnabled = &prometheusMetricsEnabled
	spec.PrometheusMetricsPort = &prometheusMetricsPort
	spec.PrometheusProcessMetricsEnabled = &prometheusProcessMetricsEnabled
	spec.ReportingInterval = &reportingInterval
	spec.ReportingTTL = &reportingTTL
	spec.RouteRefreshInterval = &routeRefreshInterval
	spec.UsageReportingEnabled = &usageReportingEnabled
	spec.UsageReportingInitialDelay = &usageReportingInitialDelay
	spec.UsageReportingInterval = &usageReportingInterval

	return spec, nil
}

// resourceCalicoFelixConfigurationCreate create a new FelixConfiguration in Calico
func resourceCalicoFelixConfigurationCreate(d *schema.ResourceData, m interface{}) error {
	calicoClient := m.(config).Client
	FelixConfigurationInterface := calicoClient.FelixConfigurations()

	FelixConfiguration, err := createFelixConfigurationApiRequest(d)
	if err != nil {
		return err
	}

	_, err = FelixConfigurationInterface.Create(ctx, FelixConfiguration, opts)
	if err != nil {
		return err
	}

	d.SetId(FelixConfiguration.ObjectMeta.Name)
	return resourceCalicoFelixConfigurationRead(d, m)
}

// resourceCalicoFelixConfigurationRead get a specific FelixConfiguration
func resourceCalicoFelixConfigurationRead(d *schema.ResourceData, m interface{}) error {
	calicoClient := m.(config).Client
	FelixConfigurationInterface := calicoClient.FelixConfigurations()

	nameFelixConfiguration := dToString(d, "metadata.0.name")

	FelixConfiguration, err := FelixConfigurationInterface.Get(ctx, nameFelixConfiguration, options.GetOptions{})

	// Handle endpoint does not exist
	if err != nil {
		if _, ok := err.(errors.ErrorResourceDoesNotExist); ok {
			d.SetId("")
			return nil
		} else {
			return err
		}
	}

	d.SetId(nameFelixConfiguration)
	d.Set("metadata.0.name", FelixConfiguration.ObjectMeta.Name)

	return nil
}

// resourceCalicoFelixConfigurationUpdate update an FelixConfiguration in Calico
func resourceCalicoFelixConfigurationUpdate(d *schema.ResourceData, m interface{}) error {
	d.Partial(false)

	calicoClient := m.(config).Client
	FelixConfigurationInterface := calicoClient.FelixConfigurations()

	FelixConfiguration, err := createFelixConfigurationApiRequest(d)
	if err != nil {
		return err
	}

	_, err = FelixConfigurationInterface.Update(ctx, FelixConfiguration, opts)
	if err != nil {
		return err
	}

	return nil
}

// resourceCalicoFelixConfigurationDelete delete an FelixConfiguration in Calico
func resourceCalicoFelixConfigurationDelete(d *schema.ResourceData, m interface{}) error {
	calicoClient := m.(config).Client
	FelixConfigurationInterface := calicoClient.FelixConfigurations()

	nameFelixConfiguration := dToString(d, "metadata.0.name")

	_, err := FelixConfigurationInterface.Delete(ctx, nameFelixConfiguration, options.DeleteOptions{})
	if err != nil {
		return err
	}

	return nil
}

// createApiRequest prepare the request of creation and update
func createFelixConfigurationApiRequest(d *schema.ResourceData) (*api.FelixConfiguration, error) {
	// Set Spec to FelixConfiguration Spec
	spec, err := dToFelixConfigurationSpec(d)
	if err != nil {
		return nil, err
	}

	// Set Metadata to Kubernetes Metadata
	objectMeta, err := dToTypeMeta(d)
	if err != nil {
		return nil, err
	}

	// Create a new BGP Configuration, with TypeMeta filled in
	// Then, fill the metadata and the spec
	newFelixConfiguration := api.NewFelixConfiguration()
	newFelixConfiguration.ObjectMeta = objectMeta
	newFelixConfiguration.Spec = spec

	return newFelixConfiguration, nil
}
