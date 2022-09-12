# syntax=docker/dockerfile:experimental

FROM alpine

ENV DOCKER_BUILDKIT=1
ENV BUILDKIT_PROGRESS=plain
ENV DOCKER_CLI_EXPERIMENTAL=enabled

RUN apk update
RUN apk add docker
RUN apk add git
RUN mkdir -p $HOME/.docker/cli-plugins
RUN if [[ "$TARGETARCH" == "arm64" ]]; then wget -O $HOME/.docker/cli-plugins/docker-buildx  "https://github.com/docker/buildx/releases/download/v0.9.1/buildx-v0.9.1.linux-arm64" ; else  wget -O $HOME/.docker/cli-plugins/docker-buildx  "https://github.com/docker/buildx/releases/download/v0.9.1/buildx-v0.9.1.linux-amd64"; fi
RUN chmod a+x $HOME/.docker/cli-plugins/docker-buildx

COPY ./docker/entrypoint.sh ./

ENTRYPOINT ["sh","entrypoint.sh" ]