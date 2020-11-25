# Gitops Util

This is [gitops](https://www.gitops.tech/) helper utility.

## How to Use

Execute ```go build -o bin/gitops cmd/gitops/main.go``` to generate binary.
Execute ```./gitops update``` command with necessary flags to update images.

## Examples

#### Update Operation Example
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
1)  ```--prefix=<prefix>``` 
2) ```--app-namespace=<ns>```
3) ```--app-name=<app-name>```

```)
./gitops update \
    --repo=<repository-ssh-url> \
    --branch=<repository-branch> \
    --key-file=<key-file-location> \
    --images=<image:tag,image2:tag> \ 
    --prefix=<prefix> \
    --app-namespace=<ns> \
    --app-name=<app-name>
```
**Valuated ```--app-path``` with other than default value will override tokenized path flags** 

### Supported Flags for Update Operation

* ```--repo``` - Git repository SSH URL
    * default value is ```""```
* ```--branch``` repository branch
    * default value is ```""```
* ```--key-file``` - SSH private key location
    * default value is ```""```
* ```--strict-host-key-checking``` - strict host key checking boolean
    * default value is ```false```
    * if ```true``` then ```known_hosts``` file will be searched from these [default locations](https://github.com/src-d/go-git/blob/master/plumbing/transport/ssh/auth_method.go#L273):
        * ```homeDirPath + "/.ssh/known_hosts"```
        * ```"/etc/ssh/ssh_known_hosts"```
* ```--push``` - git push boolean
    * default value is ```true```
* ```--images``` - images with tags
    * default value is ```""```
* ```--app-path``` path to application folder
    * default value is ```""```
* ```--prefix``` - path prefix to apply
    * default value is ```""```
* ```--app-namespace``` - application namespace
    * default value is ```""```
* ```--app-name``` - application name
    * default value is ```""```