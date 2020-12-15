package update

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/git"
)

func Run(flags *common.Flags) {
	fmt.Println("update:", flags)
	cloneOrPull(flags.Git)

}

func cloneOrPull(gitFlags common.GitFlags) {
	common.CdToGitOpsWd()
	if !git.DoesRepoExist(gitFlags.Repo) {
		cloneAndConfig(gitFlags)
	} else {
		//openedRepo := openGitOpsRepo()
		//gitPull(openedRepo)
	}
}

func cloneAndConfig(gitFlags common.GitFlags) {
	git.Clone(gitFlags)
	//repo = configRepo(repo)
}
