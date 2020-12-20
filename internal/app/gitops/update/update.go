package update

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/git"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/update/operation"
)

var repository = new(git.Repository)

func Run(flags *common.Flags) {
	addFlagsToWorkingRepo(flags)
	cloneOrPull()
	updateImages()
	applyChanges()
	pushOnDemand()
	logEndMessage()
}

func addFlagsToWorkingRepo(flags *common.Flags) {
	repository.GitFlags = flags.Git
	repository.Images = flags.Images
	repository.AppPath = flags.App.Path
	repository.KeepRegistry = flags.KeepRegistry
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

func applyChanges() {
	repository.OpenGitOpsRepo()
	repository.Add()
	repository.CommitIfModified()
}

func resetAndUpdate() {
	common.RmGitOpsWorkDir()
	cloneOrPull()
	updateImages()
	applyChanges()
	pushOnDemand()
}

func pushOnDemand() {
	if repository.GitFlags.Push {
		if err := repository.Push(); err != nil {
			resetAndUpdate()
		}
	} else {
		common.Logger.Println("commit(s) were chosen not to be pushed")
	}
}

func updateImages() {
	common.CdToAppDir(repository.Repo, repository.AppPath)
	updater := operation.Updater{Images: repository.Images, KeepRegistry: repository.KeepRegistry}
	updater.UpdateImages()
}

func logEndMessage() {
	if repository.GitFlags.Push {
		url := common.GetRemoteRepoWebUrl(repository.Repo)
		common.Logger.Println(fmt.Sprintf("repository url: %s", url))
	}
}
