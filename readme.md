# GitOps Util

This is [gitops](https://www.gitops.tech/) helper CLI

* [Compiling Source](#compiling-source)
* [How to use](#how-to-use)
* [Examples](#examples)
* [ArgoCD Commands](#argocd-commands)
    * [argocd-sync](#argocd-sync)
    * [argocd-delete](#argocd-delete)

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

## ArgoCD Commands

### argocd-sync

Syncs an ArgoCD application and waits for it to reach synced and healthy status.

OPTIONS:
* app-name - **Required**, name of the ArgoCD application [$APP_NAME]
* server - **Required**, server tcp address with port [$ARGO_CD_SERVER]
* auth-token - **Required**, authentication token [$ARGO_CD_TOKEN]
* insecure - insecure connection (default: **false**) [$ARGO_CD_INSECURE]
* timeout - timeout for single ArgoCD operation (default: **300**) [$ARGO_CD_TIMEOUT]
* async - don't wait for sync to complete (default: **false**) [$ARGO_CD_ASYNC]

Minimal example

```./gitops argocd-sync --app-name='application-name' --server='localhost:443' --auth-token='auth-token'```

Full example

```./gitops argocd-sync --app-name='application-name' --server='localhost:443' --auth-token='auth-token' --timeout=300 insecure=false async=false```

### argocd-delete

Deletes an ArgoCD application.

OPTIONS:
* app-name - **Required**, name of the ArgoCD application [$APP_NAME]
* server - **Required**, server tcp address with port [$ARGO_CD_SERVER]
* auth-token - **Required**, authentication token [$ARGO_CD_TOKEN]
* insecure - insecure connection (default: **false**) [$ARGO_CD_INSECURE]
* timeout - timeout for single ArgoCD operation (default: **300**) [$ARGO_CD_TIMEOUT]
* cascade - perform a cascaded deletion of all application resources (default: **true**) [$ARGO_CD_CASCADE]

Minimal example

```./gitops argocd-delete --app-name='application-name' --server='localhost:443' --auth-token='auth-token'```

Full example

```./gitops argocd-sync --app-name='application-name' --server='localhost:443' --auth-token='auth-token' --timeout=300 insecure=false cascade=true```