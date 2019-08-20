#!/bin/bash
set -eux -o pipefail

# super simple just for quick e2e
GO111MODULE=on go mod download
cd example/single_node 
go build
lsof -ti tcp:8080 | xargs kill || true
./single_node -delete -profile m5
./single_node -start -profile m5
export KUBECONFIG=/Users/medmac/.kube/kic-config-m5
kubectl get pods -A
sleep 10
kubectl run hello-minikube --image=k8s.gcr.io/echoserver:1.4 --port=8080
kubectl expose deployment hello-minikube --type=NodePort
sleep 25
kubectl port-forward service/hello-minikube 8080 &
sleep 5
curl http://localhost:8080/
lsof -ti tcp:8080 | xargs kill || true