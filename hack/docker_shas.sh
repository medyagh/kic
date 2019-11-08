#!/bin/bash
# set -eux -o pipefail


### to be used in automated PRs to genrate code for pkg/image

for ver in v1.11.10 v1.12.8 v1.12.9 v1.12.10 v1.13.6 v1.13.7 v1.14.3 v1.15.0 v1.15.3 v1.16.1 v1.16.2
do
    docker pull medyagh/kic:$ver | grep Digest: | awk -v x=$ver '{print "\tcase \""x"\":\n\t\treturn \"medyagh/kic:"x"@"$2"\", nil"}'
done


