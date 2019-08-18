#!/bin/bash
set -eux -o pipefail

docker push medyagh/kic:v1.15.0
docker push medyagh/kic:v1.14.3
docker push medyagh/kic:v1.13.7
docker push medyagh/kic:v1.12.9
docker push medyagh/kic:v1.11.10