package update

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"os"
)

func getCloneOptions(url string, branch string, sshKey string) *git.CloneOptions {
	if isRemoteKeyDefined(sshKey) {
		return getCloneOptionsWithKey(url, branch, sshKey)
	}
	return getCloneOptionsDefault(url, branch)
}

func getCloneOptionsWithKey(url string, branch string, sshKey string) *git.CloneOptions {
	return &git.CloneOptions{
		Auth:          getPublicKeys(sshKey),
		URL:           url,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
		Progress:      os.Stdout,
	}
}

func getCloneOptionsDefault(url string, branch string) *git.CloneOptions {
	return &git.CloneOptions{
		URL:           url,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
		Progress:      os.Stdout,
	}
}

func getPushOptions(sshKey string) *git.PushOptions {
	if isRemoteKeyDefined(sshKey) {
		return getPushOptionsWithKey(sshKey)
	}
	return getPushOptionsDefault()
}

func getPushOptionsWithKey(sshKey string) *git.PushOptions {
	return &git.PushOptions{
		Auth: getPublicKeys(sshKey),
	}
}

func getPushOptionsDefault() *git.PushOptions {
	return &git.PushOptions{}
}
