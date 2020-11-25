package update

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/util"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	goGitSsh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"golang.org/x/crypto/ssh"
	"os"
	"strings"
	"time"
)

func gitClone() *git.Repository {
	util.Logger.Println(fmt.Sprintf("git clone %s", flgs.repo))
	repo, err := git.PlainClone(getRepositoryName(), false, getCloneOptions())
	if err != nil {
		util.Logger.Println(&util.PrefixedError{Reason: err})
		os.Exit(1)
	}
	util.Logger.Println("git clone was successful")
	return repo
}

func gitAdd(repo *git.Repository) {
	wt := getWorkTree(repo)
	_, err := wt.Add(".")
	if err != nil {
		util.Logger.Println(&util.PrefixedError{Reason: err})
		os.Exit(1)
	}
}

func gitCommit(repo *git.Repository) {
	cfg := getGitConfig(repo)
	wt := getWorkTree(repo)
	commitMessage := fmt.Sprintf("Updated image(s) %s in %s", flgs.images, util.GetAppName(flgs.appPath))
	commit, commitErr := wt.Commit(commitMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name:  cfg.User.Name,
			Email: cfg.User.Email,
			When:  time.Now(),
		},
	})
	if commitErr != nil {
		util.Logger.Println(&util.PrefixedError{Reason: commitErr})
		os.Exit(1)
	}
	_, commitObjErr := repo.CommitObject(commit)
	if commitObjErr != nil {
		util.Logger.Println(&util.PrefixedError{Reason: commitObjErr})
		os.Exit(1)
	}
	util.Logger.Println("git commit was successful")
}

func gitPush(repo *git.Repository) {
	err := repo.Push(getPushOptions())
	if err != nil {
		util.Logger.Println(&util.PrefixedError{Reason: err})
		os.Exit(1)
	}
	util.Logger.Println("git push was successful")
}

func getWorkTree(repo *git.Repository) *git.Worktree {
	wt, err := repo.Worktree()
	if err != nil {
		util.Logger.Println(&util.PrefixedError{Reason: err})
		os.Exit(1)
	}
	return wt
}

func getGitConfig(repo *git.Repository) *config.Config {
	cfg, err := repo.Config()
	if err != nil {
		util.Logger.Println(&util.PrefixedError{Reason: err})
		os.Exit(1)
	}
	return cfg
}

func configRepo(repo *git.Repository) *git.Repository {
	cfg, err := repo.Config()
	if err != nil {
		util.Logger.Println(&util.PrefixedError{Reason: err})
		os.Exit(1)
	}
	cfg.User.Name = "jenkins"
	cfg.User.Email = "jenkins@localhost"
	err = repo.SetConfig(cfg)
	if err != nil {
		util.Logger.Println(&util.PrefixedError{Reason: err})
		os.Exit(1)
	}
	return repo
}

func exitIfUnmodified(repo *git.Repository) {
	wt := getWorkTree(repo)
	status := getGitStatus(wt)
	if status.IsClean() {
		util.Logger.Println("all files are unmodified, nothing to commit")
		os.Exit(0)
	}
	util.Logger.Print(fmt.Sprintf("git status: %s", status))
}

func getGitStatus(wt *git.Worktree) git.Status {
	status, err := wt.Status()
	if err != nil {
		util.Logger.Println(&util.PrefixedError{Reason: err})
		os.Exit(1)
	}
	return status
}

func openGitOpsRepo() *git.Repository {
	repo, err := git.PlainOpen(getRepositoryRootPath())
	if err != nil {
		util.Logger.Println(&util.PrefixedError{Reason: err})
		os.Exit(1)
	}
	return repo
}

func isRemoteKeyDefined() bool {
	if _, err := os.Stat(flgs.keyFile); err != nil {
		util.Logger.Println(fmt.Sprintf("coldn't use SSH key defined via flag. %s", err))
		util.Logger.Println("using SSH key defined in pipeline")
		return false
	}
	util.Logger.Println(fmt.Sprintf("using SSH key defined in %s", flgs.keyFile))
	return true
}

func getPublicKeys() *goGitSsh.PublicKeys {
	if flgs.strictHostKeyChecking {
		return getPublicKeysDefault()
	}
	return getPublicKeysNonStrict()
}

func getPublicKeysDefault() *goGitSsh.PublicKeys {
	publicKeys, err := goGitSsh.NewPublicKeysFromFile("git", flgs.keyFile, "")
	if err != nil {
		util.Logger.Println(&util.PrefixedError{Reason: err})
		os.Exit(1)
	}
	return publicKeys
}

func getPublicKeysNonStrict() *goGitSsh.PublicKeys {
	publicKeys := getPublicKeysDefault()
	publicKeys.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	return publicKeys
}

func getRepositoryName() string {
	tokens := strings.Split(flgs.repo, "/")
	lastToken := tokens[len(tokens)-1]
	return lastToken[:len(lastToken)-4]
}

func getRepositoryRootPath() string {
	return fmt.Sprintf("%s/%s/%s", util.RootPath, GitOpsWd, getRepositoryName())
}
