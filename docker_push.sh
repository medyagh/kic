#!/bin/bash
set -eux -o pipefail
docker build --tag medyagh/kic:v1.15.0 --build-arg KUBE_VER=v1.15.0 .
docker build --tag medyagh/kic:v1.14.3 --build-arg KUBE_VER=v1.14.3 .
docker build --tag medyagh/kic:v1.13.7 --build-arg KUBE_VER=v1.13.7 .
docker build --tag medyagh/kic:v1.12.9 --build-arg KUBE_VER=v1.12.9 .
docker build --tag medyagh/kic:v1.11.10 --build-arg KUBE_VER=v1.11.10 .
docker push medyagh/kic:v1.15.0
docker push medyagh/kic:v1.14.3
docker push medyagh/kic:v1.13.7
docker push medyagh/kic:v1.12.9
docker push medyagh/kic:v1.11.10