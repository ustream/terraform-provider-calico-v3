# BGP Peers

You can declare a BGP Peer with this configuration : 

```hcl
resource "calico_bgppeer" "default" {
  metadata {
    name = "Router1"
  }
  spec {
    node = "global"
    peer_ip = "192.168.0.5"
    as_number = "62523"
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
|node|If specified, the scope is node level, otherwise the scope is global.|The hostname of the node to which this peer applies.|string| 
|peerIP|The IP address of this peer.|Valid IPv4 or IPv6 address.|string| 
|asNumber|The AS Number of this peer.|A valid AS Number, may be specified in dotted notation.|integer/string| 