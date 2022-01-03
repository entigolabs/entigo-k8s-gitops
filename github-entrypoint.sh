#!/bin/bash
touch /dev/shm/key-file && chmod 600 /dev/shm/key-file && echo "$INPUT_GIT_KEY" > /dev/shm/key-file
set -x
if [ "$INPUT_COMMAND" == "update" ]
then
gitops update \
    --git-repo="$INPUT_GIT_REPO" \
    --git-branch="$INPUT_GIT_BRANCH" \
    --git-key-file="/dev/shm/key-file" \
    --images="$INPUT_IMAGES" \
    --app-path="$INPUT_APP_PATH"

elif [ "$INPUT_COMMAND" == "argocd-update" ]
then
gitops argocd-update \
    --git-key-file="/dev/shm/key-file" \
    --images="$INPUT_IMAGES" \
    --app-name="${INPUT_APP_NAME}" \
    --server="${INPUT_ARGOCD_SERVER}" \
    --auth-token="${INPUT_ARGOCD_TOKEN}"

else
  echo "Unknown command: $INPUT_COMMAND"
  exit 100
fi
