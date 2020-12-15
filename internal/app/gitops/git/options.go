package git

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"os"
)

func getCloneOptions(gitFlags common.GitFlags) *git.CloneOptions {
	if isRemoteKeyDefined(gitFlags.KeyFile) {
		return getCloneOptionsWithKey(gitFlags)
	}
	return getCloneOptionsDefault(gitFlags)
}

func getCloneOptionsWithKey(gitFlags common.GitFlags) *git.CloneOptions {
	return &git.CloneOptions{
		Auth:          getPublicKeys(gitFlags),
		URL:           gitFlags.Repo,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", gitFlags.Branch)),
		Progress:      os.Stdout,
	}
}

func getCloneOptionsDefault(gitFlags common.GitFlags) *git.CloneOptions {
	return &git.CloneOptions{
		URL:           gitFlags.Repo,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", gitFlags.Branch)),
		Progress:      os.Stdout,
	}
}

//func getPushOptions() *git.PushOptions {
//	if isRemoteKeyDefined() {
//		return getPushOptionsWithKey()
//	}
//	return getPushOptionsDefault()
//}
//
//func getPushOptionsWithKey() *git.PushOptions {
//	return &git.PushOptions{
//		Auth: getPublicKeys(),
//	}
//}
//
//func getPushOptionsDefault() *git.PushOptions {
//	return &git.PushOptions{}
//}
//
//func getPullOptions() *git.PullOptions {
//	if isRemoteKeyDefined() {
//		return getPullOptionsWithKey()
//	}
//	return getPullOptionsDefault()
//}
//
//func getPullOptionsWithKey() *git.PullOptions {
//	return &git.PullOptions{
//		Auth: getPublicKeys(),
//	}
//}
//
//func getPullOptionsDefault() *git.PullOptions {
//	return &git.PullOptions{}
//}
