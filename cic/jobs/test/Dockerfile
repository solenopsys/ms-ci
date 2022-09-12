FROM alpine/git AS CLONE_SOURCE
RUN mkdir /stuff
WORKDIR /stuff
RUN git clone https://git.alexstorm.solenopsys.org/alexstorm-hsm-ci
WORKDIR /stuff/alexstorm-hsm-ci

FROM buildkit AS BUILD_IMAGE
RUN mkdir /stuff
WORKDIR /stuff
COPY --from=CLONE_SOURCE /stuff/alexstorm-hsm-ci /stuff
RUN docker buildx build  --platform linux/amd64,linux/arm64 -t registry.alexstorm.solenopsys.org/alexstorm-hsm-ci:latest --push .
