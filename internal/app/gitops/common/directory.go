package common

import (
	"fmt"
	"os"
)

const GitOpsWd = "gitops-workdir"

var RootPath = GetWd()

func CdToRepoRoot(repo string) {
	repoRootPath := GetRepositoryRootPath(repo)
	if err := ChangeDir(repoRootPath); err != nil {
		Logger.Fatal(&PrefixedError{Reason: err})
	}
}

func CdToGitOpsWd() {
	path := fmt.Sprintf("%s/%s", RootPath, GitOpsWd)
	if _, err := os.Stat(path); err != nil {
		if err := os.Mkdir(GitOpsWd, 0755); err != nil {
			Logger.Println(&PrefixedError{Reason: err})
		}
	}
	if err := ChangeDir(path); err != nil {
		Logger.Fatal(&PrefixedError{Reason: err})
	}
}

func ChangeDir(path string) error {
	if err := os.Chdir(path); err != nil {
		return err
	}
	Logger.Printf("changed directory to: %s\n", path)
	return nil
}

func GetWd() string {
	wd, _ := os.Getwd()
	return wd
}

func RmGitOpsWd() {
	if err := ChangeDir(RootPath); err != nil {
		Logger.Fatal(&PrefixedError{Reason: err})
	}
	if err := os.RemoveAll(fmt.Sprintf("%s/%s", RootPath, GitOpsWd)); err != nil {
		Logger.Fatal(&PrefixedError{Reason: err})
	}
	Logger.Println("gitops-workdir successfully cleaned")
}
