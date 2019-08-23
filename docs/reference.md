# Reference

### config files

- /etc/kubernetes/bootstrap-kubelet.conf
- /var/lib/kubelet/config.yaml
- /etc/kubernetes/admin.conf
- /etc/kubernetes/controller-manager.conf
- /etc/kubernetes/kubelet.conf
- /etc/kubernetes/scheduler.conf


### manifests
- /etc/kubernetes/manifests/etcd.yaml
- /etc/kubernetes/manifests/kube-apiserver.yaml
- /etc/kubernetes/manifests/kube-controller-manager


### certs
- /etc/kubernetes/pki/apiserver-etcd-client.crt
- /etc/kubernetes/pki/apiserver-etcd-client.key
- /etc/kubernetes/pki/apiserver-kubelet-client.crt
- /etc/kubernetes/pki/apiserver-kubelet-client.key
- /etc/kubernetes/pki/apiserver.crt
- /etc/kubernetes/pki/apiserver.key
- /etc/kubernetes/pki/ca.crt
- /etc/kubernetes/pki/ca.key
- /etc/kubernetes/pki/front-proxy-ca.crt
- /etc/kubernetes/pki/front-proxy-ca.key
- /etc/kubernetes/pki/front-proxy-client.crt
- /etc/kubernetes/pki/front-proxy-client.key
- /etc/kubernetes/pki/sa.key
- /etc/kubernetes/pki/sa.pub


### etcd certs
- /etc/kubernetes/pki/etcd/ca.crt
- /etc/kubernetes/pki/etcd/ca.key
- /etc/kubernetes/pki/etcd/healthcheck-client.crt
- /etc/kubernetes/pki/etcd/healthcheck-client.key
- /etc/kubernetes/pki/etcd/peer.crt
- /etc/kubernetes/pki/etcd/peer.key
- /etc/kubernetes/pki/etcd/server.crt
- /etc/kubernetes/pki/etcd/server.key


### containerd socket
--container-runtime=remote --container-runtime-endpoint=/run/containerd/
--container-runtime-endpoint=/run/containerd/containerd.sock 

### ps aux :

- ```kubelet --bootstrap-kubeconfig=/etc/kubernetes/bootstrap-kubelet.conf --kubeconfig=/etc/kubernetes/kubelet.conf --config=/var/lib/kubelet/config.yaml --container-runtime=remote --container-runtime-endpoint=/run/containerd/containerd.sock --fail-swap-on=false --node-ip= --fail-swap-on=false```

- ```kube-scheduler --bind-address=127.0.0.1 --kubeconfig=/etc/kubernetes/scheduler.conf --leader-elect=true```

- ```kube-controller-manager --allocate-node-cidrs=true --authentication-kubeconfig=/etc/kubernetes/controller-manager.conf --authorization-kubeconfig=/etc/kubernetes/controller-manager.conf --bind-address=127.0.0.1 --client-ca-file=/etc/kubernetes/pki/ca.crt --cluster-cidr=10.244.0.0/16 --cluster-signing-cert-file=/etc/kubernetes/pki/ca.crt --cluster-signing-key-file=/etc/kubernetes/pki/ca.key --controllers=*,bootstrapsigner,tokencleaner --enable-hostpath-provisioner=true --kubeconfig=/etc/kubernetes/controller-manager.conf --leader-elect=true --node-cidr-mask-size=24 --requestheader-client-ca-file=/etc/kubernetes/pki/front-proxy-ca.crt --root-ca-file=/etc/kubernetes/pki/ca.crt --service-account-private-key-file=/etc/kubernetes/pki/sa.key --use-service-account-credentials=true```

- ```etcd --advertise-client-urls=https://172.17.0.2:2379 --cert-file=/etc/kubernetes/pki/etcd/server.crt --client-cert-auth=true --data-dir=/var/lib/etcd --initial-advertise-peer-urls=https://172.17.0.2:2380 --initial-cluster=m5control-plane=https://172.17.0.2:2380 --key-file=/etc/kubernetes/pki/etcd/server.key --listen-client-urls=https://127.0.0.1:2379,https://172.17.0.2:2379 --listen-peer-urls=https://172.17.0.2:2380 --name=m5control-plane --peer-cert-file=/etc/kubernetes/pki/etcd/peer.crt --peer-client-cert-auth=true --peer-key-file=/etc/kubernetes/pki/etcd/peer.key --peer-trusted-ca-file=/etc/kubernetes/pki/etcd/ca.crt --snapshot-count=10000 --trusted-ca-file=/etc/kubernetes/pki/etcd/ca.crt```



- ```kube-apiserver --advertise-address=172.17.0.2 --allow-privileged=true --authorization-mode=Node,RBAC --client-ca-file=/etc/kubernetes/pki/ca.crt --enable-admission-plugins=NodeRestriction --enable-bootstrap-token-auth=true --etcd-cafile=/etc/kubernetes/pki/etcd/ca.crt --etcd-certfile=/etc/kubernetes/pki/apiserver-etcd-client.crt --etcd-keyfile=/etc/kubernetes/pki/apiserver-etcd-client.key --etcd-servers=https://127.0.0.1:2379 --insecure-port=0 --kubelet-client-certificate=/etc/kubernetes/pki/apiserver-kubelet-client.crt --kubelet-client-key=/etc/kubernetes/pki/apiserver-kubelet-client.key --kubelet-preferred-address-types=InternalIP,ExternalIP,Hostname --proxy-client-cert-file=/etc/kubernetes/pki/front-proxy-client.crt --proxy-client-key-file=/etc/kubernetes/pki/front-proxy-client.key --requestheader-allowed-names=front-proxy-client --requestheader-client-ca-file=/etc/kubernetes/pki/front-proxy-ca.crt --requestheader-extra-headers-prefix=X-Remote-Extra- --requestheader-group-headers=X-Remote-Group --requestheader-username-headers=X-Remote-User --secure-port=6443 --service-account-key-file=/etc/kubernetes/pki/sa.pub --service-cluster-ip-range=10.96.0.0/12 --tls-cert-file=/etc/kubernetes/pki/apiserver.crt --tls-private-key-file=/etc/kubernetes/pki/apiserver.key```
