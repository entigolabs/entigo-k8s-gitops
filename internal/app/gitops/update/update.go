package update

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/git"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/update/operation"
)

var repository = new(git.Repository)

func Run(flags *common.Flags) {
	fmt.Println("update:", flags)
	repository.GitFlags = flags.Git
	cloneOrPull()
	updateImages(flags)
	//applyChanges() TODO impl this or test if this is good idea to move git into struct
}

func cloneOrPull() {
	common.CdToGitOpsWd()
	if !repository.DoesRepoExist() {
		cloneAndConfig()
	} else {
		repository.OpenGitOpsRepo()
		if err := repository.Pull(); err != nil {
			resetAndUpdate()
		}
	}
}

func cloneAndConfig() {
	repository.Clone()
	repository.ConfigRepo()
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
