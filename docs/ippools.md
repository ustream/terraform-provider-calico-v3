# IP Pools
An IP pool resource (IPPool) represents a collection of IP addresses from which Calico expects endpoint IPs to be assigned.

You can declare an ippool with this configuration : 

```hcl
resource "calico_ippool" "default" {
  metadata {
    name = "ippool-192.168.0.0"
  }
  spec {
    cidr = "192.168.0.0/24"
    ipip_mode = "Never"
    nat_outgoing = true
    disabled = false
  }
}
```

## metadata

|**Field**|**Description**|**Accepted Values**|**Schema**|
|---------|---------------|-------------------|----------|
|name|The name of this IPPool resource. Required.|Alphanumeric string with optional ., _, or -.|string|
  
  
## spec

|**Field**|**Description**|**Accepted Values**|**Schema**|
|---------|---------------|-------------------|----------|
|cidr|IP range to use for this pool.|A valid IPv4 or IPv6 CIDR. Subnet length must be /26 or less for IPv4 and /122 or less for IPv6. Must not overlap with the Link Local range 169.254.0.0/16 or fe80::/10.|string|| |
|ipip_mode|The IPIP mode defining when IPIP will be used.|Always, CrossSubnet, Never|string|Never|
|nat_outgoing|When enabled, packets sent from Calico networked containers in this pool to destinations outside of this pool will be masqueraded.|true, false|boolean|false|
|disabled|When set to true, Calico IPAM will not assign addresses from this pool.|true, false|boolean|false|

More informations on configuration [here](https://docs.projectcalico.org/v3.1/reference/calicoctl/resources/ippool).