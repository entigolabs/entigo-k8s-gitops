package git

import (
	"errors"
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"io"
	"os"
)

func (r *Repository) getCloneOptions() *git.CloneOptions {
	return &git.CloneOptions{
		Auth:          r.getAuth(),
		URL:           r.Repo,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", r.GitFlags.Branch)),
		Progress:      r.getProgressWriter(),
	}
}

func (r *Repository) getPullOptions() *git.PullOptions {
	return &git.PullOptions{
		Auth: r.getAuth(),
	}
}

func (r *Repository) getPushOptions() *git.PushOptions {
	return &git.PushOptions{
		Auth: r.getAuth(),
	}
}

func (r *Repository) getAuth() transport.AuthMethod {
	if r.isRemoteKeyDefined(r.KeyFile) {
		return r.getPublicKeys()
	} else if r.Username != "" {
		return r.getBasicAuth()
	}
	return nil
}

func (r *Repository) getProgressWriter() io.Writer {
	switch r.LoggingLevel {
	case common.DevLoggingLvl:
		return os.Stdout
	case common.ProdLoggingLvl:
		return io.Discard
	default:
		msg := fmt.Sprintf("unsupported logging level: %v", r.LoggingLevel)
		common.Logger.Fatal(&common.PrefixedError{Reason: errors.New(msg)})
	}
	return os.Stdout
}
