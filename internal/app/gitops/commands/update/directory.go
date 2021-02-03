package update

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
)

func cdToAppDir(repo string, appPath string) {
	path := fmt.Sprintf("%s/%s", common.GetRepositoryRootPath(repo), appPath)
	if err := common.ChangeDir(path); err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
}
