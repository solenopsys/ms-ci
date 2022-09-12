docker  run -e gitRegistry='git.alexstorm.solenopsys.org' \
            -e gitRepoName='alexstorm-hsm-ci' \
            -e dockerFilePatch='./cic/jobs/go_build_x.Dockerfile' \
            -e buildForArch='linux/amd64,linux/arm64' \
            -e dockerRegistry='registry.alexstorm.solenopsys.org' \
            -e targetImageName='alexstorm-hsm-ci' \
            -e imageVersion='latest' \
            --name build_test registry.alexstorm.solenopsys.org/alexstorm-buildx-job





