package update

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/git"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/update/operation"
)

func Run(flags *common.Flags) {
	fmt.Println("update:", flags)
	cloneOrPull(flags.Git)
	updateImages(flags)
}

func cloneOrPull(gitFlags common.GitFlags) {
	common.CdToGitOpsWd()
	if !git.DoesRepoExist(gitFlags.Repo) {
		cloneAndConfig(gitFlags)
	} else {
		git.OpenGitOpsRepo(gitFlags.Repo)
		if _, err := git.Pull(gitFlags, git.Repository); err != nil {
			resetAndUpdate()
		}
	}
}

func cloneAndConfig(gitFlags common.GitFlags) {
	git.Clone(gitFlags)
	git.ConfigRepo(git.Repository)
}

func resetAndUpdate() {
	common.Logger.Fatal("resetAndUpdate ->>> TODO")
	//reset()
	//cloneOrPull()
	//updateImages()
	//applyChanges()
	//pushOnDemand()
}

func updateImages(flags *common.Flags) {
	common.CdToAppDir(flags.Git.Repo, flags.App.Path)
	updater := operation.Updater{Images: flags.Images, KeepRegistry: flags.KeepRegistry}
	updater.UpdateImages()
}
