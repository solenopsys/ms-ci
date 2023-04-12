FROM alpine/git AS CLONE_SOURCE
ARG REPO_NAME
ARG USER_INFO
ARG GIT_HOST
RUN mkdir /sources
RUN git clone https://$USER_INFO$GIT_HOST/$REPO_NAME /sources

FROM --platform=$BUILDPLATFORM  golang:buster
ARG TARGETARCH
ARG GO_MAIN_FILE
COPY --from=CLONE_SOURCE /sources /sources
WORKDIR /sources
RUN go mod download
ENV GOOS=linux
ENV GOARCH=$TARGETARCH
RUN go build -o /go-binary $GO_MAIN_FILE
CMD [ "/go-binary" ]
