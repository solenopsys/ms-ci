FROM alpine/git AS CLONE_SOURCE
RUN mkdir /stuff
WORKDIR /stuff
RUN git clone https://$gitRegistry/$gitRepoName
WORKDIR /stuff/$gitRepoName

FROM registry.solenopsys.org/alexstorm-buildx-job AS BUILD_IMAGE
RUN mkdir /stuff
WORKDIR /stuff
COPY --from=CLONE_SOURCE /stuff/$gitRepoName /stuff
RUN docker buildx build  --platform $buildForArch -t $dockerRegistry/$targetImageName:$imageVersion --push .
