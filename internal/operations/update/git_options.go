package update

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"os"
)

func getCloneOptions() *git.CloneOptions {
	if isRemoteKeyDefined() {
		return getCloneOptionsWithKey()
	}
	return getCloneOptionsDefault()
}

func getCloneOptionsWithKey() *git.CloneOptions {
	return &git.CloneOptions{
		Auth:          getPublicKeys(),
		URL:           flgs.repo,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", flgs.branch)),
		Progress:      os.Stdout,
	}
}

func getCloneOptionsDefault() *git.CloneOptions {
	return &git.CloneOptions{
		URL:           flgs.repo,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", flgs.branch)),
		Progress:      os.Stdout,
	}
}

func getPushOptions() *git.PushOptions {
	if isRemoteKeyDefined() {
		return getPushOptionsWithKey()
	}
	return getPushOptionsDefault()
}

func getPushOptionsWithKey() *git.PushOptions {
	return &git.PushOptions{
		Auth: getPublicKeys(),
	}
}

func getPushOptionsDefault() *git.PushOptions {
	return &git.PushOptions{}
}

func getPullOptions() *git.PullOptions {
	if isRemoteKeyDefined() {
		return getPullOptionsWithKey()
	}
	return getPullOptionsDefault()
}

func getPullOptionsWithKey() *git.PullOptions {
	return &git.PullOptions{
		Auth: getPublicKeys(),
	}
}

func getPullOptionsDefault() *git.PullOptions {
	return &git.PullOptions{}
}
