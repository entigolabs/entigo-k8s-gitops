#!/bin/bash


tocuh /dev/shm/key-file && chmod 600 /dev/shm/key-file && echo "$INPUT_GIT_KEY" > /dev/shm/key-file && \
gitops update \
    --git-repo=$INPUT_GIT_REPO \
    --git-branch=$INPUT_GIT_BRANCH \
    --git-key-file=/dev/shm/key-file \
    --images=$INPUT_IMAGES \
    --app-path=$INPUT_APP_PATH

