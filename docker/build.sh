#!/bin/sh

REGISTRY=registry.klogsolenopsys.org
NAME=alexstorm-buildx-job
ARCHS="linux/amd64,linux/arm64"
docker  buildx build  --platform ${ARCHS} -f docker/job-build-container.Dockerfile -t ${REGISTRY}/${NAME}:latest   --push .





