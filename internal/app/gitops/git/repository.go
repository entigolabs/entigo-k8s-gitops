package git

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"time"
)

type Repository struct {
	*git.Repository
	common.GitFlags
	Images       string
	AppPath      string
	KeepRegistry bool
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

func (r *Repository) Push() error {
	pushOptions := r.getPushOptions()
	err := r.Repository.Push(pushOptions)
	if err != nil {
		return handlePushErr(err)
	} else {
		common.Logger.Println("git push was successful")
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

func (r *Repository) Add() {
	wt := r.getWorkTree()
	_, err := wt.Add(".")
	if err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
}

func (r *Repository) gitCommit() {
	cfg := r.getGitConfig()
	wt := r.getWorkTree()
	// TODO getAppName() is workname -> if(featureBranch == 'master') and if not ???
	commitMessage := fmt.Sprintf("Updated image(s) %s in %s", r.Images, getAppName(r.AppPath))
	commit, commitErr := wt.Commit(commitMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name:  cfg.User.Name,
			Email: cfg.User.Email,
			When:  time.Now(),
		},
	})
	if commitErr != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: commitErr})
	}
	_, commitObjErr := r.Repository.CommitObject(commit)
	if commitObjErr != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: commitObjErr})
	}
	common.Logger.Println("git commit was successful")
}

func (r *Repository) getGitConfig() *config.Config {
	cfg, err := r.Repository.Config()
	if err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
	return cfg
}

func (r *Repository) CommitIfModified() {
	if r.isRepoModified() {
		r.gitCommit()
	} else {
		common.Logger.Println("nothing to commit, working tree clean")
	}
}

func (r *Repository) isRepoModified() bool {
	wt := r.getWorkTree()
	status := getGitStatus(wt)
	if !status.IsClean() {
		return true
	}
	return false
}

func getGitStatus(wt *git.Worktree) git.Status {
	status, err := wt.Status()
	if err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
	return status
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
