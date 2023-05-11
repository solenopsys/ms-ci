#!/bin/sh


gitRegistry='git.solenopsys.org'
gitRepoName='alexstorm-hsm-ci'
dockerFilePatch='/mnt/c/dev/sources/goprojects/alexstorm-hsm-ci/cic/jobs/go_build_x.Dockerfile'
buildForArch='linux/amd64,linux/arm64'
dockerRegistry='registry.solenopsys.org'
targetImageName='alexstorm-hsm-ci'
imageVersion='latest'

#git clone https://$gitRegistry/$gitRepoName
cd $gitRepoName

docker run \
    -it \
    --rm \
    --security-opt seccomp=unconfined \
    --security-opt apparmor=unconfined \
    -e BUILDKITD_FLAGS=--oci-worker-no-process-sandbox \
    -v /mnt/c/dev/sources/goprojects/alexstorm-hsm-ci:/tmp/work \
    --entrypoint buildctl-daemonless.sh \
    moby/buildkit:master-rootless \
        build \
        --frontend \
        $dockerFilePatch \
        --local context=/tmp/work \
        --local dockerfile=$dockerFilePatch