#!/bin/bash
# set -eux -o pipefail

for v in v1.11.10 v1.12.8 v1.12.9 v1.12.10 v1.13.6 v1.13.7 v1.14.3 v1.15.0 v1.15.3 v1.16.1 v1.16.2
do
    docker build --tag medyagh/kic:$v --build-arg KUBE_VER=$v .
    docker push medyagh/kic:$v
done
