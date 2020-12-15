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
