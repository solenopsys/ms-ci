

FROM alpine/git AS CLONE_SOURCE
ARG REPO_NAME
ARG USER_INFO
ARG GIT_HOST
RUN mkdir /sources
RUN git clone https://$USER_INFO$GIT_HOST/$REPO_NAME /sources


FROM --platform=$BUILDPLATFORM  golang:1.19-buster
COPY --from=CLONE_SOURCE /stuff/solenopsys-hsm-ci /app
ARG TARGETARCH
WORKDIR /app
RUN go mod download
RUN  GOOS=linux GOARCH=$TARGETARCH go build  -o /go-binary /app/cmd/app/main.go

CMD [ "/go-binary" ]
