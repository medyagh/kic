# Quick Start 
```
    ./single_node -start -profile m5
	kubectl get pods -A
	kubectl run hello-minikube --image=k8s.gcr.io/echoserver:1.4 --port=8080
	kubectl expose deployment hello-minikube --type=NodePort
 	kubectl port-forward service/hello-minikube 8080
	curl http://localhost:8080/
    ./single_node -delete -profile m5


```