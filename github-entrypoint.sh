#!/bin/bash



echo "$INPUT_GIT_KEY" > /key-file
chmod 600 /key-file
cat /key-file | sed 's/./& /g'

gitops update \
    --git-repo=$INPUT_GIT_REPO \
    --git-branch=$INPUT_GIT_BRANCH \
    --git-key-file=/key-file \
    --images=$INPUT_IMAGES \
    --app-path=$INPUT_APP_PATH

