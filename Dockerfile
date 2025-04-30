FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS build
COPY . /go/gitops
ARG PACKAGE="github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
WORKDIR /go/gitops
RUN go get -d ./...
ARG TARGETARCH
ARG TARGETOS
RUN set -xe && \
    VERSION=DOCKER && \
    BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ') && \
    GIT_COMMIT=$(git rev-parse HEAD) && \
    GIT_TREE_STATE=$(if [ -z "`git status --porcelain`" ]; then echo "clean" ; else echo "dirty"; fi) && \
    LDFLAGS="-X ${PACKAGE}.version=${VERSION} \
  -X ${PACKAGE}.buildDate=${BUILD_DATE} \
  -X ${PACKAGE}.gitCommit=${GIT_COMMIT} \
  -X ${PACKAGE}.gitTreeState=${GIT_TREE_STATE}" && \
    GOOS="$TARGETOS" GOARCH="$TARGETARCH" go build \
    -o bin/gitops \
    -ldflags "${LDFLAGS} -linkmode external -extldflags -static " \
    cmd/gitops/main.go

FROM alpine:3
RUN apk add bash
ENTRYPOINT ["bash", "-c"]
COPY github-entrypoint.sh /github-entrypoint.sh
COPY --from=build /go/gitops/bin/gitops /usr/bin/gitops
