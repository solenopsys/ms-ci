FROM alpine/git AS CLONE_SOURCE
RUN mkdir /stuff
WORKDIR /stuff
RUN git clone https://git.alexstorm.solenopsys.org/alexstorm-hsm-ci
WORKDIR /stuff/alexstorm-hsm-ci

FROM --platform=$BUILDPLATFORM  golang:1.19-buster
COPY --from=CLONE_SOURCE /stuff/alexstorm-hsm-ci /app
ARG TARGETARCH
WORKDIR /app
RUN go mod download
RUN  GOOS=linux GOARCH=$TARGETARCH go build  -o /go-binary /app/cmd/app/main.go

CMD [ "/go-binary" ]
