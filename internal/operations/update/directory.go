package update

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/util"
	"os"
)

func cdToGitOpsWd() {
	path := fmt.Sprintf("%s/%s", util.GetWd(), GitOpsWd)
	if _, err := os.Stat(path); err != nil {
		if err := os.Mkdir(GitOpsWd, 0755); err != nil {
			util.Logger.Println(&util.PrefixedError{Reason: err})
		}
	}
	if err := util.ChangeDir(path); err != nil {
		util.Logger.Println(&util.PrefixedError{Reason: err})
		os.Exit(1)
	}
}

func cdToAppDir() {
	path := fmt.Sprintf("%s/%s", getRepositoryRootPath(), flgs.appPath)
	if err := util.ChangeDir(path); err != nil {
		util.Logger.Println(&util.PrefixedError{Reason: err})
		os.Exit(1)
	}
}
