package git

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/installer"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp/capability"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"regexp"
	"strings"
	"time"
)

var (
	azureHTTPSPattern = regexp.MustCompile(`^https://(?:[^@]+@)?dev\.azure\.com/[^/]+/[^/]+/_git/[^/]+`)
	azureSSHPattern   = regexp.MustCompile(`^git@ssh\.dev\.azure\.com:v3/[^/]+/[^/]+/[^/]+`)
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
	common.Logger.Printf("git clone %s\n", r.Repo)
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
	commitMessage, msgErr := r.getCommitMessage()
	if msgErr != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: msgErr})
	}
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
	return "", fmt.Errorf("unsupported command '%v' for commit messafe", r.Command)
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
	return !status.IsClean()
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

func (r *Repository) Init() {
	if !isAzureDevOps(r.Repo) {
		return
	}
	// Azure DevOps requires capabilities multi_ack / multi_ack_detailed which go-git does not support by default.
	// We need to set UnsupportedCapabilities to avoid errors when dealing with Azure DevOps.
	// See: https://github.com/go-git/go-git/blob/master/_examples/azure_devops/main.go
	common.Logger.Println("Detected Azure DevOps repository, disabling unsupported capabilities")
	transport.UnsupportedCapabilities = []capability.Capability{
		capability.ThinPack,
	}
}

func (r *Repository) Close() {
	if !isAzureDevOps(r.Repo) {
		return
	}
	// To support processes with multiple repositories, re-enable the default capabilities
	transport.UnsupportedCapabilities = []capability.Capability{
		capability.MultiACK,
		capability.MultiACKDetailed,
		capability.ThinPack,
	}
}

func isAzureDevOps(url string) bool {
	return azureHTTPSPattern.MatchString(url) || azureSSHPattern.MatchString(url)
}
