package copy

import "github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"

func cdToRepoRoot(repo string) {
	repoRootPath := common.GetRepositoryRootPath(repo)
	if err := common.ChangeDir(repoRootPath); err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
}

func cdToCopiedBranch(path string) {
	if err := common.ChangeDir(path); err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
}
