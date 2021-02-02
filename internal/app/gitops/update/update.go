package update

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/git"
	gitOpsUpdated "github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/update/updater"
)

func Run(flags *common.Flags) {
	repo := initWorkingRepo(flags)
	cloneOrPull(repo)
	updateImages(repo)
	applyChanges(repo)
	pushOnDemand(repo)
	logEndMessage(repo)
}

func initWorkingRepo(flags *common.Flags) *git.Repository {
	repository := new(git.Repository)
	repository.GitFlags = flags.Git
	repository.Images = flags.Images
	repository.AppFlags = flags.App
	repository.KeepRegistry = flags.KeepRegistry
	repository.Command = common.UpdateCmd
	repository.LoggingLevel = common.ConvStrToLoggingLvl(flags.LoggingLevel)
	return repository
}

func cloneOrPull(workingRepo *git.Repository) {
	common.CdToGitOpsWd()
	if !workingRepo.DoesRepoExist() {
		cloneAndConfig(workingRepo)
	} else {
		pullAndClone(workingRepo)
	}
}

func cloneAndConfig(workingRepo *git.Repository) {
	workingRepo.Clone()
	workingRepo.ConfigRepo()
}

func pullAndClone(workingRepo *git.Repository) {
	workingRepo.OpenGitOpsRepo()
	if err := workingRepo.Pull(); err != nil {
		resetAndUpdate(workingRepo)
	}
	workingRepo.ConfigRepo()
}

func applyChanges(workingRepo *git.Repository) {
	workingRepo.OpenGitOpsRepo()
	workingRepo.Add()
	workingRepo.CommitIfModified()
}

func resetAndUpdate(workingRepo *git.Repository) {
	common.RmGitOpsWd()
	cloneOrPull(workingRepo)
	updateImages(workingRepo)
	applyChanges(workingRepo)
	pushOnDemand(workingRepo)
}

func pushOnDemand(workingRepo *git.Repository) {
	if workingRepo.GitFlags.Push {
		if err := workingRepo.Push(); err != nil {
			resetAndUpdate(workingRepo)
		}
	} else {
		common.Logger.Println("commit(s) were chosen not to be pushed")
	}
}

func updateImages(workingRepo *git.Repository) {
	cdToAppDir(workingRepo.Repo, workingRepo.AppFlags.Path)
	updater := gitOpsUpdated.Updater{Images: workingRepo.Images, KeepRegistry: workingRepo.KeepRegistry}
	updater.UpdateImages()
}

func logEndMessage(workingRepo *git.Repository) {
	if workingRepo.GitFlags.Push {
		url := common.GetRemoteRepoWebUrl(workingRepo.Repo)
		common.Logger.Println(fmt.Sprintf("repository url: %s", url))
	}
}
