package git

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/go-git/go-git/v5"
)

var Repository *git.Repository = new(git.Repository)

func Clone(gitFlags common.GitFlags) {
	common.Logger.Println(fmt.Sprintf("git clone %s", gitFlags.Repo))
	repo, err := git.PlainClone(getRepositoryName(gitFlags.Repo), false, getCloneOptions(gitFlags))
	if err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
	common.Logger.Println("git clone was successful")
	Repository = repo
}

func Pull(gitFlags common.GitFlags, repo *git.Repository) (*git.Repository, error) {
	wt := getWorkTree(repo)
	err := wt.Pull(getPullOptions(gitFlags))
	if err != nil {
		return nil, handlePullErr(err)
	} else {
		common.Logger.Println("git pull was successful")
	}
	return repo, nil
}

func OpenGitOpsRepo(repoSshUrl string) {
	repo, err := git.PlainOpen(getRepositoryRootPath(repoSshUrl))
	if err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
	Repository = repo
}

func getWorkTree(repo *git.Repository) *git.Worktree {
	wt, err := repo.Worktree()
	if err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
	return wt
}

func ConfigRepo(repo *git.Repository) *git.Repository {
	cfg, err := repo.Config()
	if err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
	cfg.User.Name = "jenkins"
	cfg.User.Email = "jenkins@localhost"
	err = repo.SetConfig(cfg)
	if err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
	return repo
}
