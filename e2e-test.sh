#!/bin/bash
set -eux -o pipefail

 if ! kubectl &>/dev/null; then
    echo "WARNING: No kubectl installation found in your enviroment."
    curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl
    chmod +x ./kubectl
    sudo mv ./kubectl /usr/local/bin/kubectl
fi


 
# super simple just for quick e2e
GO111MODULE=on go mod download
cd example/single_node 
GO111MODULE=on go build
lsof -ti tcp:8080 | xargs kill || true
./single_node -delete -profile m5
./single_node -start -profile m5
export KUBECONFIG=$HOME/.kube/kic-config-m5
sleep 3 
kubectl wait deployment -l k8s-app=kube-dns --for condition=available --timeout=120s -n kube-system
kubectl get pods -A
kubectl wait pod -l component=kube-scheduler --for condition=Initialized --timeout=100s -n kube-system
kubectl wait pod -l component=kube-scheduler --for condition=ContainersReady --timeout=100s -n kube-system
kubectl get pods -A
kubectl wait pod -l component=kube-apiserver --for condition=Initialized --timeout=100s -n kube-system || true
kubectl wait pod -l component=kube-apiserver --for condition=ContainersReady --timeout=100s -n kube-system || true
kubectl wait pod -l component=kube-apiserver --for condition=PodScheduled --timeout=100s -n kube-system || true
kubectl wait pod -l component=kube-apiserver --for condition=Ready --timeout=100s -n kube-system || true
kubectl wait pod -l component=etcd --for condition=Ready --timeout=100s -n kube-system || true

kubectl run hello-minikube --image=k8s.gcr.io/echoserver:1.4 --port=8080
kubectl expose deployment hello-minikube --type=NodePort
kubectl wait deployment -l run=hello-minikube --for condition=available --timeout=100s
kubectl port-forward service/hello-minikube 8080 &
sleep 3
curl http://localhost:8080/
lsof -ti tcp:8080 | xargs kill || true