# Gitops Util

This is [gitops](https://www.gitops.tech/) helper utility.

## How to use

To run the binary execute ```./gitops``` command with necessary flags.

### Flag definitions 
* Flag_K8sNs - kubernetes Namespace
* img-name - image name to be changed
* argo-app - Argo app name
* branch - Git branch name
* ssh-key - SSH key for gitops
* deploy-env - deployment environment (default env is dev)

<!-- TODO remove this -->
Example command to execute
```./gitops -k8s-ns=ns-test -img-name=nginx:1.9.3 -deploy-env=kk-test -argo-app=app2```

## Project Specific Commands

- ```go build -o bin/gitops cmd/gitops/main.go``` - generate binary
- ```go test ./...``` - run tests