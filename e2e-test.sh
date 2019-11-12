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
./out/e2e -remove -profile m5 || true


# start a cluster
echo "Starting a cluster with 2 cpu and 2 GB ram" && ./out/e2e -start -profile m5 -cpu 2 -memory 2000m
export KUBECONFIG=$HOME/.kube/kic-config-m5


# test status command
./out/e2e -status -profile m5 | grep "Running"

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
docker exec m5-control-plane cat /kic/kubeadm.conf | grep  apiServerEndpoint
docker exec m5-control-plane stat -c '%a' kic/kubeadm.conf | grep 644

# check that container was creatred for control-plane
docker ps || grep "m5-control-plane"

# pulling an image to load to a new kic cluster
docker images || true  # list images before
docker pull busybox
docker tag busybox e2e-example-img
docker images || true  # list images after

## Create a second cluster test and load an image to it.
echo "Starting a second cluster" && ./out/e2e -start -profile cluster2 

## load an image from user machine to cluster
echo "Loading image from user machine to cluster" && ./out/e2e -profile cluster2  -image e2e-example-img -load=true
echo "Checking if image is loaded" && docker exec cluster2-control-plane ctr -n k8s.io images ls  | grep e2e-example-img

## copy file from user machine to cluster
touch copy-test.txt
echo "copy test" > copy-test.txt
echo "Copying file from user machine to cluster" && ./out/e2e -profile cluster2 -cp -src=copy-test.txt -dest=/etc/copy-test.txt
echo "Checking if file was copied" && docker exec cluster2-control-plane cat /etc/copy-test.txt | grep "copy test"

## remove file from cluster
echo "Removing file from cluster" && ./out/e2e -profile cluster2 -rm-file -src=/etc/copy-test.txt
echo "Checking if file was removed" && docker exec cluster2-control-plane test ! -f /etc/copy-test.txt

## pause a node
echo "Pausing a node" && ./out/e2e -pause -profile cluster2
echo "Checking if node was paused" && docker inspect --format '{{.State.Status}}' cluster2-control-plane  | grep paused

## stop a node
echo "Stopping a node" && ./out/e2e -stop -profile cluster2
echo "Checking if node was stopped" && docker inspect --format '{{.State.Status}}' cluster2-control-plane  | grep exited

# remove our cluster in the end
./out/e2e -remove -profile cluster2

lsof -ti tcp:8080 | xargs kill || true
