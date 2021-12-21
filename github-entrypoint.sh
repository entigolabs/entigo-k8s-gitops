#!/bin/bash

echo "-----BEGIN RSA PRIVATE KEY-----" > /key-file
chmod 600 /key-file
echo $INPUT_GIT_KEY >> /key-file
echo "-----END RSA PRIVATE KEY-----"  >> /key-file

gitops update \
    --git-repo=$INPUT_GIT_REPO \
    --git-branch=$INPUT_GIT_BRANCH \
    --git-key-file=/key-file \
    --images=$INPUT_IMAGES \
    --app-path=$INPUT_APP_PATH

