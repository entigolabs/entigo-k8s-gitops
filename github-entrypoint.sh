#!/bin/bash
set -x

GIT_AUTH_ARGS=""

if [ -n "$INPUT_GIT_KEY" ]; then
  touch /dev/shm/key-file && chmod 600 /dev/shm/key-file && echo "$INPUT_GIT_KEY" > /dev/shm/key-file
  GIT_AUTH_ARGS+="--git-key-file=/dev/shm/key-file "
fi

if [ -n "$INPUT_GIT_USERNAME" ] && [ -n "$INPUT_GIT_PASSWORD" ]; then
  GIT_AUTH_ARGS+="--git-username=$INPUT_GIT_USERNAME --git-password=$INPUT_GIT_PASSWORD "
fi

if [ -z "$GIT_AUTH_ARGS" ]; then
  echo "Error: At least one authentication method must be provided with INPUT_GIT_KEY and/or INPUT_GIT_USERNAME."
  exit 1
fi

if [ "$INPUT_COMMAND" == "update" ]
then
gitops update \
    --git-repo="$INPUT_GIT_REPO" \
    --git-branch="$INPUT_GIT_BRANCH" \
    $GIT_AUTH_ARGS \
    --images="$INPUT_IMAGES" \
    --app-path="$INPUT_APP_PATH"

elif [ "$INPUT_COMMAND" == "argocd-update" ]
then
gitops argocd-update \
    $GIT_AUTH_ARGS \
    --images="$INPUT_IMAGES" \
    --app-name="${INPUT_APP_NAME}" \
    --server="${INPUT_ARGOCD_SERVER}" \
    --auth-token="${INPUT_ARGOCD_TOKEN}"

else
  echo "Unknown command: $INPUT_COMMAND"
  exit 100
fi
