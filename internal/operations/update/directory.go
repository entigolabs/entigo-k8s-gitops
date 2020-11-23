package update

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/util"
	"os"
)

func cdToGitOps() {
	path := fmt.Sprintf("%s/%s", util.GetWd(), GitOpsWd)
	if _, err := os.Stat(path); err != nil {
		os.Mkdir(GitOpsWd, 0755)
	}
	if err := util.ChangeDir(path); err != nil {
		util.Logger.Println(&util.PrefixedError{Reason: err})
		os.Exit(1)
	}
}

func cdToWorkPath() {
	path := fmt.Sprintf("%s/%s", util.GetWd(), flgs.appPath)
	if err := util.ChangeDir(path); err != nil {
		util.Logger.Println(&util.PrefixedError{Reason: err})
		os.Exit(1)
	}
}
