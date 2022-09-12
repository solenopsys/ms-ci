https://github.com/GoogleContainerTools/kaniko

https://github.com/moby/buildkit

buildctl build \
--frontend dockerfile.v0 \
--opt platform=linux/amd64,linux/arm64 \
--output type=image,name=docker.io/username/image,push=true \
...


https://github.com/tektoncd/catalog/blob/main/task/buildkit-daemonless/0.1/buildkit-daemonless.yaml