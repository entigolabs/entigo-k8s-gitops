package git

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/go-git/go-git/v5"
)

type Repository struct {
	*git.Repository
	common.GitFlags
}

func (r *Repository) Clone() {
	common.Logger.Println(fmt.Sprintf("git clone %s", r.Repo))
	repo, err := git.PlainClone(common.GetRepositoryName(r.Repo), false, r.getCloneOptions())
	if err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
	common.Logger.Println("git clone was successful")
	r.Repository = repo
}

func (r *Repository) Pull() error {
	wt := r.getWorkTree()
	err := wt.Pull(r.getPullOptions())
	if err != nil {
		return handlePullErr(err)
	} else {
		common.Logger.Println("git pull was successful")
	}
	return nil
}

func (r *Repository) OpenGitOpsRepo() {
	repo, err := git.PlainOpen(common.GetRepositoryRootPath(r.Repo))
	if err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
	r.Repository = repo
}

func (r *Repository) getWorkTree() *git.Worktree {
	wt, err := r.Repository.Worktree()
	if err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
	return wt
}

func (r *Repository) ConfigRepo() {
	cfg, err := r.Repository.Config()
	if err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
	cfg.User.Name = "jenkins"
	cfg.User.Email = "jenkins@localhost"
	err = r.Repository.SetConfig(cfg)
	if err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
}
