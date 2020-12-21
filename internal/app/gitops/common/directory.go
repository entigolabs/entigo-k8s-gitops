package common

import (
	"fmt"
	"os"
)

const GitOpsWd = "gitops-workdir"

var RootPath = GetWd()

func CdToGitOpsWd() {
	path := fmt.Sprintf("%s/%s", RootPath, GitOpsWd)
	if _, err := os.Stat(path); err != nil {
		if err := os.Mkdir(GitOpsWd, 0755); err != nil {
			Logger.Println(&PrefixedError{Reason: err})
		}
	}
	if err := changeDir(path); err != nil {
		Logger.Fatal(&PrefixedError{Reason: err})
	}
}

func CdToRepoRoot(repo string) {
	repoRootPath := GetRepositoryRootPath(repo)
	if err := changeDir(repoRootPath); err != nil {
		Logger.Fatal(&PrefixedError{Reason: err})
	}
}

func CdToAppDir(repo string, appPath string) {
	path := fmt.Sprintf("%s/%s", GetRepositoryRootPath(repo), appPath)
	if err := changeDir(path); err != nil {
		Logger.Fatal(&PrefixedError{Reason: err})
	}
}

func changeDir(path string) error {
	if err := os.Chdir(path); err != nil {
		return err
	}
	Logger.Println(fmt.Sprintf("changed directory to: %s", path))
	return nil
}

func GetWd() string {
	wd, _ := os.Getwd()
	return wd
}

func RmGitOpsWorkDir() {
	if err := changeDir(RootPath); err != nil {
		Logger.Fatal(&PrefixedError{Reason: err})
	}
	if err := os.RemoveAll(fmt.Sprintf("%s/%s", RootPath, GitOpsWd)); err != nil {
		Logger.Fatal(&PrefixedError{Reason: err})
	}
	Logger.Println("gitops-workdir successfully cleaned")
}
