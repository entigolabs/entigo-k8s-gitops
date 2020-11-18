package util

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"os"
	"regexp"
	"strings"
	"time"
)

func getRepository(url string, branch string, sshKey string) {
	repo := gitClone(url, branch, sshKey)
	repo = configRepo(repo)
}

func gitClone(url string, branch string, sshKey string) *git.Repository {
	fmt.Printf("Git clone %s\n", url)
	repo, err := git.PlainClone("", false, getCloneOptions(url, branch, sshKey))
	if err != nil {
		logger.Println(&prefixedError{err})
		os.Exit(1)
	}
	fmt.Println("Git clone was successful!")
	return repo
}

func isRemoteKeyDefined(sshKey string) bool {
	if _, err := os.Stat(sshKey); err != nil {
		fmt.Printf("Coldn't use SSH key defined via flag. %s\n", err)
		fmt.Println("Using SSH key defined in pipeline.")
		return false
	}
	fmt.Printf("Using SSH key defined in %s\n", sshKey)
	return true
}

func getPublicKeys(sshKey string) *ssh.PublicKeys {

	publicKeys, err := ssh.NewPublicKeysFromFile("git", sshKey, "")
	if err != nil {
		logger.Println(&prefixedError{err})
		os.Exit(1)
	}
	return publicKeys
}

func configRepo(repo *git.Repository) *git.Repository {
	cfg, err := repo.Config()

	if err != nil {
		logger.Println(&prefixedError{err})
		os.Exit(1)
	}
	cfg.User.Name = "jenkins"
	cfg.User.Email = "jenkins@localhost"
	repo.SetConfig(cfg)

	return repo
}

func commit(vars deploymentVariables) {
	repo := openDeploymentRepo()
	wt := getWorkTree(repo)
	gitAdd(wt)
	exitIfUnmodified(wt)
	gitCommit(vars, repo, wt)
}

func push(sshKey string) {
	repo := openDeploymentRepo()
	gitPush(repo, sshKey)
}

func exitIfUnmodified(wt *git.Worktree) {
	status := getGitStatus(wt)
	if status.IsClean() {
		fmt.Println("All files are unmodified. Nothing to commit.")
		os.Exit(0)
	}
	fmt.Print(fmt.Sprintf("Git status: %s", status))
}

func gitCommit(vars deploymentVariables, repo *git.Repository, wt *git.Worktree) {
	cfg := getGitConfig(repo)
	commitMessage := fmt.Sprintf("Updated image(s) %s in %s", vars.ImageNameFlag, vars.WorkName)
	commit, commitErr := wt.Commit(commitMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name:  cfg.User.Name,
			Email: cfg.User.Email,
			When:  time.Now(),
		},
	})
	if commitErr != nil {
		logger.Println(&prefixedError{commitErr})
		os.Exit(1)
	}
	_, commitObjErr := repo.CommitObject(commit)
	if commitObjErr != nil {
		logger.Println(&prefixedError{commitObjErr})
		os.Exit(1)
	}
	fmt.Println("Git commit was successful!")
}

func gitPush(repo *git.Repository, sshKey string) {
	pushErr := repo.Push(getPushOptions(sshKey))
	if pushErr != nil {
		logger.Println(&prefixedError{pushErr})
		os.Exit(1)
	}
	fmt.Println("Git push was successful!")
}

func getGitConfig(repo *git.Repository) *config.Config {
	cfg, cfgErr := repo.Config()
	if cfgErr != nil {
		logger.Println(&prefixedError{cfgErr})
		os.Exit(1)
	}
	return cfg
}

func getGitStatus(wt *git.Worktree) git.Status {
	status, statusErr := wt.Status()
	if statusErr != nil {
		logger.Println(&prefixedError{statusErr})
		os.Exit(1)
	}
	return status
}

func gitAdd(wt *git.Worktree) {
	_, addErr := wt.Add(".")
	if addErr != nil {
		logger.Println(&prefixedError{addErr})
		os.Exit(1)
	}
}

func getWorkTree(repo *git.Repository) *git.Worktree {
	wt, wtErr := repo.Worktree()
	if wtErr != nil {
		logger.Println(&prefixedError{wtErr})
		os.Exit(1)
	}
	return wt
}

func openDeploymentRepo() *git.Repository {
	repoPath := fmt.Sprintf("%s/%s", rootPath, deploymentDir)
	repo, openErr := git.PlainOpen(repoPath)
	if openErr != nil {
		logger.Println(&prefixedError{openErr})
		os.Exit(1)
	}
	return repo
}

func composeFeatureBranchName(branchName string) string {
	strings.Replace(branchName, "%2F", "/", -1)
	strings.Replace(branchName, "_", "-", -1)
	branchTokens := strings.Split(branchName, "/")
	branchName = branchTokens[len(branchTokens)-1]

	if hasValidCharactersForBranchName(branchName) {
		if len(branchName) > 32 {
			fmt.Println("Too long branch name, trimming it to 32 characters.")
			branchName = branchName[0:32]
		}
	} else {
		fmt.Println("Detected illegal characters" + branchName + ", falling back to md5 sum name.")
		branchName = genMD5Hash(branchName)
	}
	return strings.ToLower(branchName)
}

func hasValidCharactersForBranchName(s string) bool {
	matched, err := regexp.MatchString("^[A-Za-z0-9-]+$", s)
	if err != nil {
		logger.Println(&prefixedError{err})
		os.Exit(1)
	}
	return matched
}
