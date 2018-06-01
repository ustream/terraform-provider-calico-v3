# BGP Peers
A BGP configuration resource (BGPConfiguration) represents BGP specific configuration options for the cluster.

You can declare a BGP Peer with this configuration : 

```hcl
resource "calico_bgpconfiguration" "default" {
  metadata {
    name = "Config1"
  }
  spec {
    log_severity_screen = "Warning"
    node_to_node_mesh_enabled = false
    as_number = 62523
  }
}
```

## metadata

|**Field**|**Description**|**Accepted Values**|**Schema**|
|---------|---------------|-------------------|----------|
|name|The name of this IPPool resource. Required.|Alphanumeric string with optional ., _, or -.|string|

* The resource with the name default has a specific meaning - this contains the BGP global default configuration.
* The resources with the name `node.<nodename>` contain the node-specific overrides, and will be applied to the node <nodename>. When deleting a node the FelixConfiguration resource associated with the node will also be deleted.
  

## spec

|**Field**|**Description**|**Accepted Values**|**Schema**|**Default**|
|---------|---------------|-------------------|----------|-----------|
|log_severity_screen|Global log level|Debug, Info, Warning, Error, Fatal|string|Info|
|node_to_node_mesh_enabled|Full BGP node-to-node mesh|true, false|string|true|
|as_number|The AS Number of this peer.|A valid AS Number, may be specified in dotted notation.|integer/string|64512|

More informations on configuration [here](https://docs.projectcalico.org/v3.1/reference/calicoctl/resources/bgpconfig)