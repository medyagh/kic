#!/bin/bash
set -eux -o pipefail

 if ! kubectl &>/dev/null; then
    echo "WARNING: No kubectl installation found in your enviroment."
    curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl
    chmod +x ./kubectl
    sudo mv ./kubectl /usr/local/bin/kubectl
fi


 # Build Example
make out/e2e

# clean up the previous runs (if any)
lsof -ti tcp:8080 | xargs kill || true
./out/e2e -delete -profile m5

# start a cluster
./out/e2e -start -profile m5 -cpu 2 -memory 2000m
export KUBECONFIG=$HOME/.kube/kic-config-m5

# wait for things to be up print out pods -A for logs to see
sleep 1 
kubectl wait deployment -l k8s-app=kube-dns --for condition=available --timeout=300s -n kube-system || true
kubectl get pods -A || true
kubectl wait pod -l component=kube-scheduler --for condition=Initialized --timeout=100s -n kube-system || true
kubectl wait pod -l component=kube-scheduler --for condition=ContainersReady --timeout=100s -n kube-system || true
kubectl get pods -A || true
kubectl wait pod -l component=kube-apiserver --for condition=Initialized --timeout=100s -n kube-system || true
kubectl wait pod -l component=kube-apiserver --for condition=ContainersReady --timeout=100s -n kube-system || true
kubectl wait pod -l component=kube-apiserver --for condition=PodScheduled --timeout=100s -n kube-system || true
kubectl wait pod -l component=kube-apiserver --for condition=Ready --timeout=100s -n kube-system || true
kubectl get pods -A || true
kubectl wait pod -l component=etcd --for condition=Ready --timeout=100s -n kube-system || true
kubectl get pods -A || true

# deploy an example app
# make a service for it
# check if the service is accessiable.
kubectl run hello-minikube --image=k8s.gcr.io/echoserver:1.4 --port=8080
kubectl expose deployment hello-minikube --type=NodePort
kubectl wait deployment -l run=hello-minikube --for condition=available --timeout=100s
kubectl port-forward service/hello-minikube 8080 &
sleep 3
curl http://localhost:8080/


# test config file content and perm on the node
docker exec m5control-plane cat /kic/kubeadm.conf | grep  apiServerEndpoint
docker exec m5control-plane stat -c '%a' kic/kubeadm.conf | grep 644
# todo check if it is right permission

# check that container was creatred for control-plane
docker ps || grep "m5control-plane"


# pulling an image to load to a new kic cluster
docker images || true
docker pull busybox
docker images || true

## Creatre a second cluster test and load an image to it.
docker save busybox > busybox.tar
./out/e2e -start -profile m6 -image ./busybox.tar


# delete our cluster in the end
./out/e2e -delete -profile m5
./out/e2e -delete -profile m6

lsof -ti tcp:8080 | xargs kill || true

