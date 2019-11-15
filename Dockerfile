ARG KUBE_VER
FROM kindest/node:$KUBE_VER
RUN apt-get update && apt-get install -y \
  sudo \
  dnsutils \
  && rm -rf /var/lib/apt/lists/*
RUN echo "kic! ${TRAVIS_COMMIT}-${KUBE_VER}" > "/kic.txt"