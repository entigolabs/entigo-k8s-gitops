# Gitops Util

This is [gitops](https://www.gitops.tech/) helper utility.

## How to Use

Execute ```go build -o bin/gitops cmd/gitops/main.go``` to generate binary.
Execute ```./gitops update``` command with necessary flags to update images.

### Supported Flags for Update Operation

* ```--repo``` - Git repository SSH URL
* ```--branch``` repository branch
* ```--key-file``` - SSH private key location
* ```--images``` - images with tags
* ```--app-path``` path to application folder