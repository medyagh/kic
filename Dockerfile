ARG KUBE_VER
FROM kindest/node:$KUBE_VER
RUN echo "KIC!" > "/kic.txt"