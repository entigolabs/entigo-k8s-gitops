# GitOps Util

This is [gitops](https://www.gitops.tech/) helper CLI

## Compiling Source

- ```go build -o bin/gitops cmd/gitops/main.go```

## How to Use

This GitOps utility supports 5 commands:

- ```run``` - run update and copy logic
- ```update``` - updates specified images
- ```copy``` - copies from master branch to specified branch
- ```argocd-sync``` - syncs an ArgoCD application
- ```argocd-delete``` - deletes an ArgoCD application

## Examples

#### Update Command Example
##### Example With Application Path Flag
```
./gitops update \
    --git-repo=<repository-ssh-url> \
    --git-branch=<repository-branch> \
    --git-key-file=<key-file-location> \
    --images=<image:tag,image2:tag> \ 
    --app-path=<app-path>
```

##### Example With Tokenized Path Flags 

Tokenized path flags: 
1) ```--app-prefix=<app-prefix>``` 
2) ```--app-namespace=<app-namespace>```
3) ```--app-name=<app-name>```

```
./gitops update \
    --git-repo=<repository-ssh-url> \
    --git-branch=<repository-branch> \
    --git-key-file=<key-file-location> \
    --images=<image:tag,image2:tag> \ 
    --app-prefix=<app-prefix> \
    --app-namespace=<app-namespace> \
    --app-name=<app-name>
```
**Valuated ```--app-path``` with other than default value will override tokenized path flags** 