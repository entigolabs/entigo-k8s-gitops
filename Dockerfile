FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS build
COPY . /go/gitops
WORKDIR /go/gitops
RUN go get -d ./...
ARG GITHUB_SHA=main
ARG VERSION=latest
ARG TARGETARCH
ARG TARGETOS

RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    GOOS="$TARGETOS" GOARCH="$TARGETARCH" go build -o /tmp/bin/gitops -ldflags \
     "-X github.com/entigolabs/entigo-infralib-agent/common.version=${VERSION} \
      -X github.com/entigolabs/entigo-infralib-agent/common.buildDate=$(date -u +'%Y-%m-%dT%H:%M:%SZ') \
      -X github.com/entigolabs/entigo-infralib-agent/common.gitCommit=${GITHUB_SHA} \
      -X github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common.gitTreeState=clean \
               -extldflags -static -s -w" cmd/gitops/main.go

FROM alpine:3
RUN apk --no-cache add bash
ENTRYPOINT ["bash", "-c"]
COPY github-entrypoint.sh /github-entrypoint.sh
COPY --from=build /tmp/bin/gitops /usr/bin/gitops