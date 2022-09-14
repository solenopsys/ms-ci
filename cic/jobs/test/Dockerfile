

FROM alpine/git AS CLONE_SOURCE
ARG REPO_NAME
RUN echo "REPO_NAME: $REPO_NAME"
RUN mkdir /stuff
WORKDIR /stuff
RUN git clone https://admin:root@git.alexstorm.solenopsys.org/$REPO_NAME
WORKDIR /stuff/solenopsys-hsm-ci

FROM --platform=$BUILDPLATFORM  golang:1.19-buster
COPY --from=CLONE_SOURCE /stuff/solenopsys-hsm-ci /app
ARG TARGETARCH
WORKDIR /app
RUN go mod download
RUN  GOOS=linux GOARCH=$TARGETARCH go build  -o /go-binary /app/cmd/app/main.go

CMD [ "/go-binary" ]
