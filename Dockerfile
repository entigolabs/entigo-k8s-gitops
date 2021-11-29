FROM golang:1.15 as build
COPY . /go/gitops
RUN cd /go/gitops && go build -o bin/gitops -ldflags "-linkmode external -extldflags -static" cmd/gitops/main.go

FROM harbor.fleetcomplete.dev/dockerhub/library/centos:8
RUN yum update -y && \
    yum install -y git unzip jq curl vim python3-pip && pip3 install ruamel.yaml && \
    yum -y clean all && \
    rm -rf /var/tmp/yum-root-* /var/log/anaconda/journal.log && \
    curl -sSL -o /usr/bin/argocd https://github.com/argoproj/argo-cd/releases/download/$(curl --silent "https://api.github.com/repos/argoproj/argo-cd/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')/argocd-linux-amd64 && \
    curl -sSL -o /usr/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl && \
    chmod +x /usr/bin/argocd /usr/bin/kubectl
COPY  --from=build /go/gitops/bin/gitops /usr/bin/gitops
COPY *.py /usr/bin/

