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

More informations on configuration [here](https://docs.projectcalico.org/v3.1/reference/calicoctl/resources/felixconfig)