#!/bin/sh



git clone https://$gitRegistry/$gitRepoName
cd $gitRepoName
DOCKER_BUILDKIT=1 docker buildx build -f $dockerFilePatch --platform $buildForArch -t $dockerRegistry/$targetImageName:$imageVersion --push .

