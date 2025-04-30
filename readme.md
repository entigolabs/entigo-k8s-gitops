# GitOps Util

This is [gitops](https://www.gitops.tech/) helper CLI

* [Compiling Source](#compiling-source)
* [How to use](#how-to-use)
* [Examples](#examples)
* [ArgoCD Commands](#argocd-commands)
    * [argocd-get](#argocd-get)
    * [argocd-sync](#argocd-sync)
    * [argocd-update](#argocd-update)
    * [argocd-delete](#argocd-delete)

## Compiling Source

- ```go build -o bin/gitops cmd/gitops/main.go```

## How to Use

This GitOps utility supports 8 commands:

- ```run``` - run update and copy logic
- ```update``` - updates specified images
- ```copy``` - copy from source (default: master) into destination and install changes from ```install.txt```
- ```notify``` - notifies [Jira](https://www.atlassian.com/software/jira) about new deployment. Read more how to use in [docs/jira-notification](./docs/jira-deployment-notification.md)
- ```argocd-get``` - gets an ArgoCD application as json output
- ```argocd-sync``` - syncs an ArgoCD application
- ```argocd-update``` - updates an ArgoCD application, combines argocd-get, git update and argocd-sync
- ```argocd-delete``` - deletes an ArgoCD application

## Examples

Git commands require either the key file or username and password.
Repo url needs to be either ssh or https.

#### Copy Command
```
copy
  --git-repo=<git-repo-url>
  --git-branch=<git-repo-branch-name>
  --git-key-file=<ssh-private-key-location>
  --git-username=<git-username>
  --git-password=<git-password>
  --app-prefix=<path-prefix-to-apply>
  --app-prefix-argo=<path-to-argoapp-configuration>
  --app-prefix-yaml=<path-to-yaml-configuration>
  --app-namespace=<application-namespace-name>
  --app-name=<application-nam>
  --app-dest-branch=<new-branch-name>
```

#### Update Command
##### Example: Update With Application Path Flag
```
./gitops update \
    --git-repo=<repository-url> \
    --git-branch=<repository-branch> \
    --git-key-file=<key-file-location> \
    --git-username=<git-username> \
    --git-password=<git-password> \
    --images=<image:tag,image2:tag> \ 
    --app-path=<app-path>
```

##### Example: Update With Tokenized Path Flags
Tokenized path flags: 
1) ```--app-prefix=<app-prefix>``` 
2) ```--app-namespace=<app-namespace>```
3) ```--app-name=<app-name>```

```
./gitops update \
    --git-repo=<repository-url> \
    --git-branch=<repository-branch> \
    --git-key-file=<key-file-location> \
    --git-username=<git-username> \
    --git-password=<git-password> \
    --images=<image:tag,image2:tag> \ 
    --app-prefix=<app-prefix> \
    --app-namespace=<app-namespace> \
    --app-name=<app-name>
```
**Valuated ```--app-path``` with other than default value will override tokenized path flags**

##### Example: Update With Notification
```
./gitops update
    --git-repo=<repository-url> \
    --git-branch=<repository-branch> \
    --git-key-file=<key-file-location> \
    --git-username=<git-username> \
    --git-password=<git-password> \
    --images=<image:tag,image2:tag> \ 
    --app-path=<app-path>
    --notify-api-url=<api-endpoint>
    --notify-env=<deployment-env-name>
    --notify-registry-uri=<img-registry-uri>
    --notify-auth-token=<api-auth-token>
```

#### Notify Command
##### Example: Basic Usage
```
./gitops notify
    --url=<api-endpoint>
    --deployment-env=<environment-name>
    --old-img=<image:tag>
    --new-img=<image:tag>
    --registry-uri=<rigistry-uri>
    --auth-token=<key=value>
```

## ArgoCD Commands

### argocd-get

Gets an ArgoCD application as json in standard output.

OPTIONS:
* app-name - **Required**, name of the ArgoCD application [$APP_NAME]
* server - **Required**, server tcp address with port [$ARGO_CD_SERVER]
* auth-token - **Required**, authentication token [$ARGO_CD_TOKEN]
* insecure - insecure connection (default: **false**) [$ARGO_CD_INSECURE]
* timeout - timeout for single ArgoCD operation (default: **300**) [$ARGO_CD_TIMEOUT]
* grpc-web - Enables gRPC-web protocol. Useful if Argo CD server is behind proxy which does not support HTTP2 (default: **false**) [$ARGO_CD_GRPC_WEB]
* refresh - Refresh application data when retrieving (default: **true**) [$ARGO_CD_REFRESH]

Minimal example

```./gitops argocd-get --app-name='application-name' --server='localhost:443' --auth-token='auth-token'```

Full example

```./gitops argocd-get --app-name='application-name' --server='localhost:443' --auth-token='auth-token' --timeout=300 --insecure=false --grpc-web=false --refresh=true```

### argocd-sync

Syncs an ArgoCD application and waits for it to reach synced and healthy status.

OPTIONS:
* app-name - **Required**, name of the ArgoCD application [$APP_NAME]
* server - **Required**, server tcp address with port [$ARGO_CD_SERVER]
* auth-token - **Required**, authentication token [$ARGO_CD_TOKEN]
* insecure - insecure connection (default: **false**) [$ARGO_CD_INSECURE]
* timeout - timeout for single ArgoCD operation (default: **300**) [$ARGO_CD_TIMEOUT]
* grpc-web - Enables gRPC-web protocol. Useful if Argo CD server is behind proxy which does not support HTTP2 (default: **false**) [$ARGO_CD_GRPC_WEB]
* async - don't wait for sync to complete (default: **false**) [$ARGO_CD_ASYNC]
* wait-failure - Fail the command when waiting for the sync to complete exceeds the timeout (default: **true**) [$ARGO_CD_WAIT_FAILURE]

Minimal example

```./gitops argocd-sync --app-name='application-name' --server='localhost:443' --auth-token='auth-token'```

Full example

```./gitops argocd-sync --app-name='application-name' --server='localhost:443' --auth-token='auth-token' --timeout=300 --insecure=false --grpc-web=false --refresh=true --async=false --wait-failure=true```

### argocd-update

Updates an ArgoCD application. Combines argocd-get, git update and argocd-sync commands. Either Git Key or Git Username and Password must be provided.

OPTIONS:
* app-name - **Required**, name of the ArgoCD application [$APP_NAME]
* server - **Required**, server tcp address with port [$ARGO_CD_SERVER]
* auth-token - **Required**, authentication token [$ARGO_CD_TOKEN]
* git-key-file - SSH private key location [$GIT_KEY_FILE]
* git-username - git username [$GIT_USERNAME]
* git-password - git password [$GIT_PASSWORD]
* images [-i] - **Required**, images with tags, comma separated list [$IMAGES]
* insecure - insecure connection (default: **false**) [$ARGO_CD_INSECURE]
* timeout - timeout for single ArgoCD operation (default: **300**) [$ARGO_CD_TIMEOUT]
* grpc-web - Enables gRPC-web protocol. Useful if Argo CD server is behind proxy which does not support HTTP2 (default: **false**) [$ARGO_CD_GRPC_WEB]
* git-strict-host-key-checking - strict host key checking (default: **false**) [$GIT_STRICT_HOST_KEY_CHECKING]
* git-push - push changes (default: **true**) [$GIT_PUSH]
* git-author-name value - Git author name (default: **jenkins**) [$GIT_AUTHOR_NAME]
* git-author-email value - Git author email (default: **jenkins@localhost**) [$GIT_AUTHOR_EMAIL]
* keep-registry [-k] - keeps registry part of the changeable image (default: **false**) [$KEEP_REGISTRY]
* deployment-strategy [-s] - change deployment strategy (RollingUpdate | Recreate) (default: if not defined then strategy will remain unchanged) [$DEPLOYMENT-STRATEGY]
* recursive - updates directories and their contents recursively (default: **false**) [$RECURSIVE]
* refresh - Refresh application data when retrieving (default: **true**) [$ARGO_CD_REFRESH]
* async - don't wait for argoCD sync to complete (default: **false**) [$ARGO_CD_ASYNC]
* wait-failure - Fail the command when waiting for the sync to complete exceeds the timeout (default: **true**) [$ARGO_CD_WAIT_FAILURE]

Minimal example

```./gitops argocd-update --app-name='application-name' --server='localhost:443' --auth-token='auth-token' --git-token-file='file' --images='image:1'```

### argocd-delete

Deletes an ArgoCD application.

OPTIONS:
* app-name - **Required**, name of the ArgoCD application [$APP_NAME]
* server - **Required**, server tcp address with port [$ARGO_CD_SERVER]
* auth-token - **Required**, authentication token [$ARGO_CD_TOKEN]
* insecure - insecure connection (default: **false**) [$ARGO_CD_INSECURE]
* timeout - timeout for single ArgoCD operation (default: **300**) [$ARGO_CD_TIMEOUT]
* grpc-web - Enables gRPC-web protocol. Useful if Argo CD server is behind proxy which does not support HTTP2 (default: **false**) [$ARGO_CD_GRPC_WEB]
* cascade - perform a cascaded deletion of all application resources (default: **true**) [$ARGO_CD_CASCADE]

Minimal example

```./gitops argocd-delete --app-name='application-name' --server='localhost:443' --auth-token='auth-token'```

Full example

```./gitops argocd-delete --app-name='application-name' --server='localhost:443' --auth-token='auth-token' --timeout=300 --insecure=false --grpc-web=false --cascade=true```