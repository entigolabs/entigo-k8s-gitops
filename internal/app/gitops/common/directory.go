package common

import (
	"fmt"
	"os"
)

const gitOpsWd = "gitops-workdir"

func CdToGitOpsWd() {
	path := fmt.Sprintf("%s/%s", getWd(), gitOpsWd)
	if _, err := os.Stat(path); err != nil {
		if err := os.Mkdir(gitOpsWd, 0755); err != nil {
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

func getWd() string {
	wd, _ := os.Getwd()
	return wd
}
