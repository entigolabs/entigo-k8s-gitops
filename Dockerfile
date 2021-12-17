FROM golang:1.15 as build
COPY . /go/gitops
ARG PACKAGE="github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
WORKDIR /go/gitops
RUN set -xe && \
    VERSION=$(cat ./VERSION) && \
    BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ') && \
    GIT_COMMIT=$(git rev-parse HEAD) && \
    GIT_TREE_STATE=$(if [ -z "`git status --porcelain`" ]; then echo "clean" ; else echo "dirty"; fi) && \
    LDFLAGS="-X ${PACKAGE}.version=${VERSION} \
  -X ${PACKAGE}.buildDate=${BUILD_DATE} \
  -X ${PACKAGE}.gitCommit=${GIT_COMMIT} \
  -X ${PACKAGE}.gitTreeState=${GIT_TREE_STATE}" && \
    go build \
    -o bin/gitops \
    -ldflags "${LDFLAGS} -linkmode external -extldflags -static " \
    cmd/gitops/main.go

FROM alpine:3
COPY  --from=build /go/gitops/bin/gitops /usr/bin/gitops
ENTRYPOINT ["sh"]