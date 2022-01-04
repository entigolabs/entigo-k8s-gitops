#!/usr/bin/env bash
REALPATH=$(which realpath)
if [ -z $REALPATH ]; then
realpath() {
    [[ $1 = /* ]] && echo "$1" || echo "$PWD/${1#./}"
}
fi
SCRIPT_PATH=$(realpath $(dirname "$0"))
PACKAGE="github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
VERSION=$(cat ${SCRIPT_PATH}/VERSION)
BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT=$(git rev-parse HEAD)
GIT_TREE_STATE=$(if [ -z "`git status --porcelain`" ]; then echo "clean" ; else echo "dirty"; fi)

LDFLAGS="-X ${PACKAGE}.version=${VERSION} \
  -X ${PACKAGE}.buildDate=${BUILD_DATE} \
  -X ${PACKAGE}.gitCommit=${GIT_COMMIT} \
  -X ${PACKAGE}.gitTreeState=${GIT_TREE_STATE}"

go mod tidy && \
go build -o bin/gitops -ldflags "${LDFLAGS} -linkmode external -extldflags -static" cmd/gitops/main.go