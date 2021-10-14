# extremity-ingress-controller
External Kubernetes Ingress Controller build on ha-proxy-ingress-controller


#  PROJECT UNDER ACTIVE DEVELOPMENT


## Purpose

The purpose of this project is to make external ingresses easier and more secure


## Current architecture 

1. Starts BIRD
2. Starts extremety
    1. Adds a BGPConfiguration  
```
apiVersion: crd.projectcalico.org/v1
kind: BGPConfiguration
metadata:
 name: default
spec:
 logSeverityScreen: Info
 nodeToNodeMeshEnabled: true
 asNumber: 65000
```
1. extremety queries the golang Kubernetes API for a list of the nodes/watch for added nodes
    1. Adds a BIRD .conf for newely added nodes
    2. Triggers a BIRD configuration file (to update the config)


## Example docker run commmand

```
docker run \
    --cap-add=NET_ADMIN \
    --network kind \
    -v '/mnt/user/k8s/config':'/root/.kube/config':'rw' \
    d0d921358269 --external \
    --configmap=default/haproxy-kubernetes-ingress \
    --program=/usr/sbin/haproxy \
    --disable-ipv6 \
    --ipv4-bind-address=0.0.0.0
```


## Current Supported Matrix

### Supported Kubernetes CNI
| Name | Backend | Datastore | Status |
| :---: | :---: | | :---: | | :---: | 
| Calico | BIRD | Kubernetes | Pre-Alpha |


### Supported Ingresses
| Name | Status |
| :---: | | :---: |
| haproxy-kubernetes-ingress | Pre-Alpha |


### Supported Service-Mesh
| Name | Status |
| :---: | :---: |
| Linkerd | Planned |




