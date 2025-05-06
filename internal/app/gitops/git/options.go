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
	"strings"
)

func (r *Repository) getCloneOptions() *git.CloneOptions {
	reference := plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", r.GitFlags.Branch))
	if r.GitFlags.Branch == string(plumbing.HEAD) {
		reference = plumbing.HEAD
	}
	return &git.CloneOptions{
		Auth:          r.getAuth(),
		URL:           r.getRepo(),
		ReferenceName: reference,
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
	if r.KeyFile != "" && r.Username != "" {
		if strings.HasPrefix(r.Repo, "git@") {
			return r.getPublicKeys()
		} else if strings.HasPrefix(r.Repo, "http") {
			return r.getBasicAuth()
		}
	} else if r.KeyFile != "" {
		return r.getPublicKeys()
	} else if r.Username != "" {
		return r.getBasicAuth()
	}
	return nil
}

func (r *Repository) getRepo() string {
	if r.KeyFile != "" && r.Username != "" {
		return r.Repo
	} else if r.KeyFile != "" && strings.HasPrefix(r.Repo, "http") {
		common.Logger.Println(&common.Warning{Reason: errors.New("SSH key file is defined, but repo URL is HTTP, converting URL to SSH")})
		return convertToSSH(r.Repo)
	} else if r.Username != "" && strings.HasPrefix(r.Repo, "git@") {
		common.Logger.Println(&common.Warning{Reason: errors.New("username is defined, but repo URL is SSH, converting URL to HTTP")})
		return convertToHTTP(r.Repo)
	}
	return r.Repo
}

func convertToSSH(repo string) string {
	repo = strings.Replace(repo, "https://", "", 1)
	repo = strings.Replace(repo, "http://", "", 1)
	parts := strings.Split(repo, "/")
	if len(parts) >= 2 {
		return fmt.Sprintf("git@%s:%s/%s", parts[0], parts[1], strings.Join(parts[2:], "/"))
	}
	return repo
}

func convertToHTTP(repo string) string {
	repo = strings.Replace(repo, "git@", "", 1)
	repo = strings.Replace(repo, ":", "/", 1)
	return fmt.Sprintf("https://%s", repo)
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
