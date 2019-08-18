FROM kindest/node:${KUBE_VER}
RUN echo "KIC!" > "/kic.txt"