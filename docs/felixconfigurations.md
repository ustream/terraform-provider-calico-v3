# Felix Configurations

A Felix configuration resource represents Felix configuration options for the cluster.

You can declare a Felix Configuration with this configuration : 

```hcl
resource "calico_felixconfiguration" "default" {
  metadata {
    name = "node.test"
  }
  spec {
    node = "global"
    peer_ip = "192.168.0.5"
    as_number = 62523
  }
}
```

## metadata
|**Field**|**Description**|**Accepted Values**|**Schema**|**Default**|
|---------|---------------|-------------------|----------|-----------|
|name|Unique name to describe this resource instance. Required.|Alphanumeric string with optional ., _, or -.|string|

* Calico automatically creates a resource named `default` containing the global default configuration settings for Felix. You can use calicoctl to view and edit these settings
* The resources with the name `node.<nodename>` contain the node-specific overrides, and will be applied to the `node <nodename>`. When deleting a node the FelixConfiguration resource associated with the node will also be deleted.


## spec
|**Field**|**Description**|**Accepted Values**|**Schema**|**Default**|
|---------|---------------|-------------------|----------|-----------|
|chain_insert_mode|Controls whether Felix hooks the kernel’s top-level iptables chains by inserting a rule at the top of the chain or by appending a rule at the bottom. Insert is the safe default since it prevents Calico’s rules from being bypassed. If you switch to Append mode, be sure that the other rules in the chains signal acceptance by falling through to the Calico rules, otherwise the Calico policy will be bypassed.|Insert, Append|string|Insert|
|default_endpoint_to_host_action|This parameter controls what happens to traffic that goes from a workload endpoint to the host itself (after the traffic hits the endpoint egress policy). By default Calico blocks traffic from workload endpoints to the host itself with an iptables “DROP” action. If you want to allow some or all traffic from endpoint to host, set this parameter to Return or Accept. Use Return if you have your own rules in the iptables “INPUT” chain; Calico will insert its rules at the top of that chain, then Return packets to the “INPUT” chain once it has completed processing workload endpoint egress policy. Use Accept to unconditionally accept packets from workloads after processing workload endpoint egress policy.|Drop, Return, Accept|string|Drop|
|failsafe_inbound_host_ports|UDP/TCP protocol/port pairs that Felix will allow incoming traffic to host endpoints on irrespective of the security policy. This is useful to avoid accidentally cutting off a host with incorrect configuration. The default value allows SSH access, etcd, BGP and DHCP.| |List of ProtoPort|- protocol: tcp<br>port: 22<br>- protocol: udp<br>port: 68<br>- protocol: tcp<br>port: 179<br>- protocol: tcp<br>port: 2379<br>- protocol: tcp<br>port: 2380<br>- protocol: tcp<br>port: 6666<br>- protocol: tcp<br>port: 6667|
|failsafe_outbound_host_ports|UDP/TCP protocol/port pairs that Felix will allow outgoing traffic from host endpoints to irrespective of the security policy. This is useful to avoid accidentally cutting off a host with incorrect configuration. The default value opens etcd’s standard ports to ensure that Felix does not get cut off from etcd as well as allowing DHCP and DNS.| |List of ProtoPort|- protocol: udp<br>port: 53<br>- protocol: udp<br>port: 67<br>- protocol: tcp<br>port: 179<br>- protocol: tcp<br>port: 2379<br>- protocol: tcp<br>port: 2380<br>- protocol: tcp<br>port: 6666<br>- protocol: tcp<br>port: 6667<br>|
|ignore_loose_rpf|Set to true to allow Felix to run on systems with loose reverse path forwarding (RPF). Warning: Calico relies on “strict” RPF checking being enabled to prevent workloads, such as VMs and privileged containers, from spoofing their IP addresses and impersonating other workloads (or hosts). Only enable this flag if you need to run with “loose” RPF and you either trust your workloads or have another mechanism in place to prevent spoofing.|boolean|boolean|false|
|interface_exclude|A comma-separated list of interface names that should be excluded when Felix is resolving host endpoints. The default value ensures that Felix ignores Kubernetes’ internal kube-ipvs0 device.|string|string|kube-ipvs0|
|interface_prefix|The interface name prefix that identifies workload endpoints and so distinguishes them from host endpoint interfaces. Note: in environments other than bare metal, the orchestrators configure this appropriately. For example our Kubernetes and Docker integrations set the ‘cali’ value, and our OpenStack integration sets the ‘tap’ value.|string|string|cali|
|ipip_enabled|Whether Felix should configure an IPinIP interface on the host. Set automatically to true by calico/node or calicoctl when you create an IPIP-enabled pool.|boolean|false|| 
|ipip_mtu|The MTU to set on the tunnel device. See Configuring MTU|int|int|1440|
|ipsets_refresh_interval|Period, in seconds, at which Felix re-checks the IP sets in the dataplane to ensure that no other process has accidentally broken Calico’s rules. Set to 0 to disable IP sets refresh. Note: the default for this value is lower than the other refresh intervals as a workaround for a Linux kernel bug that was fixed in kernel version 4.11. If you are using v4.11 or greater you may want to set this to a higher value to reduce Felix CPU usage.|int|int|10|
|iptables_filter_allow_action|This parameter controls what happens to traffic that is accepted by a Felix policy chain in the iptables filter table (i.e. a normal policy chain). The default will immediately Accept the traffic. Use Return to send the traffic back up to the system chains for further processing.|Accept, Return|string|Accept|
|iptables_lock_file_path|Location of the iptables lock file. You may need to change this if the lock file is not in its standard location (for example if you have mapped it into Felix’s container at a different path).|string|string|/run/xtables.lock|
|iptables_lock_probe_interval|Time, in milliseconds, that Felix will wait between attempts to acquire the iptables lock if it is not available. Lower values make Felix more responsive when the lock is contended, but use more CPU.|int|int|50|
|iptables_lock_timeout|Time, in seconds, that Felix will wait for the iptables lock, or 0, to disable. To use this feature, Felix must share the iptables lock file with all other processes that also take the lock. When running Felix inside a container, this requires the /run directory of the host to be mounted into the calico/node or calico/felix container.|int|int|0 (Disabled)|
|iptables_mangle_allow_action|This parameter controls what happens to traffic that is accepted by a Felix policy chain in the iptables mangle table (i.e. a pre-DNAT policy chain). The default will immediately Accept the traffic. Use Return to send the traffic back up to the system chains for further processing.|Accept, Return|string|Accept|
|iptables_mark_mask|Mask that Felix selects its IPTables Mark bits from. Should be a 32 bit hexadecimal number with at least 8 bits set, none of which clash with any other mark bits in use on the system.|netmask|netmask|0xff000000|
|iptables_post_write_check_interval|Period, in seconds, after Felix has done a write to the dataplane that it schedules an extra read back in order to check the write was not clobbered by another process. This should only occur if another application on the system doesn’t respect the iptables lock.|int|int|1|
|iptables_refresh_interval|Period, in seconds, at which Felix re-checks all iptables state to ensure that no other process has accidentally broken Calico’s rules. Set to 0 to disable iptables refresh.|int|int|90|
|ipv6_support|IPv6 support for Felix|true, false|boolean|true|
|log_file_path|The full path to the Felix log. Set to "" to disable file logging.|string|string|/var/log/calico/felix.log|
|log_prefix|The log prefix that Felix uses when rendering LOG rules.|string|string|calico-packet|
|log_severity_file|The log severity above which logs are sent to the log file.|Same as logSeveritySys|string|Info|
|log_severity_screen|The log severity above which logs are sent to the stdout.|Same as LogSeveritySys|string|Info|
|log_severity_sys|The log severity above which logs are sent to the syslog. Set to "" for no logging to syslog.|Debug, Info, Warning, Error, Fatal|string|Info|
|max_ipset_size|Maximum size for the ipsets used by Felix to implement tags. Should be set to a number that is greater than the maximum number of IP addresses that are ever expected in a tag.|int|int|1048576|
|metadata_addr|The IP address or domain name of the server that can answer VM queries for cloud-init metadata. In OpenStack, this corresponds to the machine running nova-api (or in Ubuntu, nova-api-metadata). A value of none (case insensitive) means that Felix should not set up any NAT rule for the metadata path.|IPv4, hostname, none|string|127.0.0.1|
|metadata_port|The port of the metadata server. This, combined with global.MetadataAddr (if not ‘None’), is used to set up a NAT rule, from 169.254.169.254:80 to MetadataAddr:MetadataPort. In most cases this should not need to be changed.|int|int|8775|
|prometheus_go_metrics_enabled|Set to false to disable Go runtime metrics collection, which the Prometheus client does by default. This reduces the number of metrics reported, reducing Prometheus load.|boolean|boolean|true|
|prometheus_metrics_enabled|Set to true to enable the experimental Prometheus metrics server in Felix.|boolean|boolean|false|
|prometheus_metrics_port|Experimental: TCP port that the Prometheus metrics server should bind to.|int|int|9091|
|prometheus_process_metrics_enabled|Set to false to disable process metrics collection, which the Prometheus client does by default. This reduces the number of metrics reported, reducing Prometheus load.|boolean|boolean|true|
|reporting_interval|Interval at which Felix reports its status into the datastore or 0 to disable. Must be non-zero in OpenStack deployments.|int|int|30|
|reporting_ttl|Time-to-live setting for process-wide status reports.|int|int|90|
|route_refresh_interval|Period, in seconds, at which Felix re-checks the routes in the dataplane to ensure that no other process has accidentally broken Calico’s rules. Set to 0 to disable route refresh.|int|int|90|
|usage_reporting_enabled|Reports anonymous Calico version number and cluster size to projectcalico.org. Logs warnings returned by the usage server. For example, if a significant security vulnerability has been discovered in the version of Calico being used.|boolean|boolean|true|
|usage_reporting_initial_delay|Minimum initial delay before first usage report, in seconds.|int|int|300|
|usage_reporting_interval|The interval at which Felix does usage reports, in seconds. The default is 1 day.|int|int|86400|

### Proto Port
|**Field**|**Description**|**Accepted Values**|**Schema**|**Default**|
|---------|---------------|-------------------|----------|-----------|
|port|The exact port match|0-65535|int|
|protocol|The protocol match|tcp, udp|string|

More informations [here](https://docs.projectcalico.org/v3.1/reference/calicoctl/resources/felixconfig)