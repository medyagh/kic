ARG KUBE_VER
FROM kindest/node:$KUBE_VER
RUN apt-get update && apt-get install -y dnsutils
RUN echo "KIC! $KUBE_VER $TRAVIS_COMMIT" > "/kic.txt"