ARG KUBE_VER
FROM kindest/node:$KUBE_VER
RUN apt-get update
RUN apt-get install dnsutils
RUN echo "KIC! $KUBE_VER $TRAVIS_COMMIT" > "/kic.txt"