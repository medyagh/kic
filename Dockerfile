ARG KUBE_VER
FROM kindest/node:$KUBE_VER
RUN echo "KIC! $KUBE_VER $TRAVIS_COMMIT" > "/kic.txt"