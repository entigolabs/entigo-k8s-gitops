package copy

import "github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"

func cdToCopiedBranch(path string) {
	if err := common.ChangeDir(path); err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
}

func cdToArgoApp(path string) {
	if err := common.ChangeDir(path); err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
}
