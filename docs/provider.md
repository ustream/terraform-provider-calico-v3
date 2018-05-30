# Provider

You can configure this provider with two type of backend. You can configure it with env var (in italic in this document).

## ETCD
```hcl
provider "calico" {
  backend_type="etcdv3"
  etcd_endpoints="https://127.0.0.1:2379"
  etcd_ca_cert_file="/etc/ssl/etcd/ca.crt"
  etcd_cert_file="/etc/ssl/etcd/etcd.crt"
  etcd_key_file="/etc/ssl/etcd/etcd.key"
}
```

|**Field**|**Environment variables**|**Description**|
|---------|-------------------------|---------------|
|backend_type|CALICO_BACKEND_TYPE|Type of the backend (etcdv3 or kubernetes)|
|etcd_endpoints|CALICO_ETCD_ENDPOINTS|List of ETCD endpoints, separate by ","|  
|etcd_ca_cert_file|CALICO_ETCD_CA_CERT_FILE|Cert of the CA (Only if SSL is enabled on your ETCD Server)| 
|etcd_cert_file|CALICO_ETCD_CERT_FILE|Client cert file to connect to ETCD (Only if SSL is enabled on your ETCD Server)|
|etcd_key_file|CALICO_ETCD_KEY_FILE|Key file to connect to ETCD (Only if SSL is enabled on your ETCD Server)| 

## Kubernetes
```hcl
provider "calico" {
  backend_type="kubernetes"
  kube_config="~/.kube/config"
  kube_apiserver="https://127.0.0.1"
  kube_token="/var/serviceaccounts/secret"
  kube_ca_cert_file="/etc/ssl/kubernetes/ca.crt"
  kube_cert_file="/etc/ssl/kubernetes/kubernetes.crt"
  kube_key_file="/etc/ssl/kubernetes/kubernetes.key"
}
```

|**Field**|**Environment variables**|**Description**|
|---------|-------------------------|---------------|
|backend_type|CALICO_BACKEND_TYPE|Type of the backend (etcdv3 or kubernetes)|
|kube_config|CALICO_KUBE_CONFIG|Kube config to connect to Kubernetes server| 
|kube_api_server|CALICO_KUBE_API_SERVER|APIServer to fetch (Optional if the apiserver is fill in the kube config)|
|kube_token|CALICO_KUBE_TOKEN|Token to connect to the APIServer|
|kube_ca_cert_file|CALICO_KUBE_CA_FILE|Cert of the CA (Only if SSL is enabled on your Kubernetes)|
|kube_cert_file|CALICO_KUBE_CERT_FILE|Client cert file to connect to APIServer (Only if SSL is enabled on your Kubernetes)|
|kube_key_file|CALICO_KUBE_KEY_FILE|Key file to connect to APIServer (Only if SSL is enabled on your Kubernetes)|
