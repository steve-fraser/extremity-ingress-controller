docker run \
    --cap-add=NET_ADMIN \
    --network kind \
    -v '/mnt/user/k8s/config':'/root/.kube/config':'rw' \
    unraid-test --external \
    --configmap=default/haproxy-kubernetes-ingress \
    --program=/usr/sbin/haproxy \
    --disable-ipv6 \
    --ipv4-bind-address=0.0.0.0
