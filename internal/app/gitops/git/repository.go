package git

import (
	"errors"
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/installer"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"strings"
	"time"
)

// TODO refactor Repository
type Repository struct {
	*git.Repository
	common.GitFlags
	common.AppFlags
	common.DeploymentNotificationFlags
	Images             string
	KeepRegistry       bool
	Command            common.Command
	LoggingLevel       common.LoggingLevel
	DeploymentStrategy common.DeploymentStrategy
	Recursive          bool
	InstallHistory     []installer.InstallHistory
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
	common.Logger.Println(fmt.Sprintf("pushOptions RemoteName: %s, RemoteURL: %s", pushOptions.RemoteName, pushOptions.RemoteURL))
	err := r.Repository.Push(pushOptions)
	if err != nil {
		return handlePushErr(err)
	} else {
		common.Logger.Println("git push was successful")
	}
	return nil
}

func (r *Repository) OpenGitOpsRepo() {
	common.Logger.Println("In OpenGitOpsRepo method")
	repo, err := git.PlainOpen(common.GetRepositoryRootPath(r.Repo))
	if err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
	r.Repository = repo
}

func (r *Repository) Add() {
	common.Logger.Println("In repository.go Add method")
	wt := r.getWorkTree()
	common.Logger.Println("finished getWorkTree method")
	_, err := wt.Add(".")
	if err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	} else {
		common.Logger.Println("finished Add method without error")
	}
}

func (r *Repository) gitCommit() {
	common.Logger.Println("In gitCommit method")
	cfg := r.getGitConfig()
	common.Logger.Println(fmt.Sprintf("cfg.Core.Worktree: %s", cfg.Core.Worktree))
	// common.Logger.Println(fmt.Sprintf("cfg.Core.Worktree: %s", cfg.URLs[0].Name))
	wt := r.getWorkTree()
	commitMessage, msgErr := r.getCommitMessage()
	common.Logger.Println(fmt.Sprintf("commitMessage: %s", commitMessage))
	if msgErr != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: msgErr})
	}
	common.Logger.Println("Starting to produce commit at wt.Commit")
	commit, commitErr := wt.Commit(commitMessage, &git.CommitOptions{
		All: true,
		Author: &object.Signature{
			Name:  cfg.User.Name,
			Email: cfg.User.Email,
			When:  time.Now(),
		},
	})
	if commitErr != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: commitErr})
	}
	common.Logger.Println("Starting to commit at r.Repository.CommitObject")
	_, commitObjErr := r.Repository.CommitObject(commit)
	if commitObjErr != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: commitObjErr})
	}
	common.Logger.Println("git commit was successful")
}

func (r *Repository) getCommitMessage() (string, error) {
	switch r.Command {
	case common.UpdateCmd:
		return fmt.Sprintf("updated image(s) %s in %s", r.Images, r.getAppName()), nil
	case common.CopyCmd:
		return fmt.Sprintf("copied configurations from %s/%s to %s/%s", r.AppFlags.Name, r.AppFlags.SourceBranch, r.getAppName(), r.AppFlags.DestBranch), nil
	case common.DeleteCmd:
		return fmt.Sprintf("deleted %s and %s.yaml", r.AppFlags.DestBranch, r.AppFlags.DestBranch), nil
	}
	return "", errors.New(fmt.Sprintf("unsupported command '%v' for commit messafe", r.Command))
}

func (r *Repository) getGitConfig() *config.Config {
	cfg, err := r.Repository.Config()
	if err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
	return cfg
}

func (r *Repository) CommitIfModified() {
	common.Logger.Println("In CommitIfModified method")
	if r.isRepoModified() {
		common.Logger.Println("Repo is modified, starting with gitCommit method")
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
	common.Logger.Println("In getWorkTree method")
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
	cfg.User.Name = r.GitFlags.AuthorName
	cfg.User.Email = r.GitFlags.AuthorEmail
	err = r.Repository.SetConfig(cfg)
	if err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
}

func (r *Repository) getAppName() string {
	if r.AppFlags.Name == "" {
		pathTokens := strings.Split(r.AppFlags.Path, "/")
		return pathTokens[len(pathTokens)-1]
	}
	return r.AppFlags.Name
}
