# GitOps Util

This is [gitops](https://www.gitops.tech/) helper CLI

## Compiling Source

- ```go build -o bin/gitops cmd/gitops/main.go```

## How to Use

This GitOps utility supports 3 commands:

- ```run``` - run update and copy logic
- ```update``` - updates specified images
- ```copy``` - copies from master branch to specified branch

## Examples

#### Update Command Example
##### Example With Application Path Flag
```
./gitops update \
    --repo=<repository-ssh-url> \
    --branch=<repository-branch> \
    --key-file=<key-file-location> \
    --images=<image:tag,image2:tag> \ 
    --app-path=<app-path>
``` 

##### Example With Tokenized Path Flags 

Tokenized path flags: 
1) ```--prefix=<app-prefix>``` 
2) ```--namespace=<app-namespace>```
3) ```--name=<app-name>```

```)
./gitops update \
    --repo=<repository-ssh-url> \
    --branch=<repository-branch> \
    --key-file=<key-file-location> \
    --images=<image:tag,image2:tag> \ 
    --prefix=<app-prefix> \
    --namespace=<app-namespace> \
    --name=<app-name>
```
**Valuated ```--app-path``` with other than default value will override tokenized path flags** 