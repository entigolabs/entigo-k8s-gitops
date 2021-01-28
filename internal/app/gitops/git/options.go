package git

import (
	"errors"
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"io"
	"io/ioutil"
	"os"
)

func (r Repository) getCloneOptions() *git.CloneOptions {
	if isRemoteKeyDefined(r.KeyFile) {
		return r.getCloneOptionsWithKey()
	}
	return r.getCloneOptionsDefault()
}

func (r Repository) getCloneOptionsWithKey() *git.CloneOptions {
	return &git.CloneOptions{
		Auth:          r.getPublicKeys(),
		URL:           r.Repo,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", r.GitFlags.Branch)),
		Progress:      r.getProgressWriter(),
	}
}

func (r Repository) getPullOptions() *git.PullOptions {
	if isRemoteKeyDefined(r.KeyFile) {
		return r.getPullOptionsWithKey()
	}
	return getPullOptionsDefault()
}

func (r Repository) getPullOptionsWithKey() *git.PullOptions {
	return &git.PullOptions{
		Auth: r.getPublicKeys(),
	}
}

func getPullOptionsDefault() *git.PullOptions {
	return &git.PullOptions{}
}

func (r Repository) getPushOptions() *git.PushOptions {
	if isRemoteKeyDefined(r.KeyFile) {
		return r.getPushOptionsWithKey()
	}
	return getPushOptionsDefault()
}

func (r Repository) getPushOptionsWithKey() *git.PushOptions {
	return &git.PushOptions{
		Auth: r.getPublicKeys(),
	}
}

func getPushOptionsDefault() *git.PushOptions {
	return &git.PushOptions{}
}

func (r Repository) getCloneOptionsDefault() *git.CloneOptions {
	return &git.CloneOptions{
		URL:           r.Repo,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", r.GitFlags.Branch)),
		Progress:      r.getProgressWriter(),
	}
}

func (r Repository) getProgressWriter() io.Writer {
	switch r.LoggingLevel {
	case "dev":
		return os.Stdout
	case "prod":
		return ioutil.Discard
	default:
		msg := fmt.Sprintf("unsupported logger level: %s", r.LoggingLevel)
		common.Logger.Fatal(&common.PrefixedError{Reason: errors.New(msg)})
	}
	return os.Stdout
}
