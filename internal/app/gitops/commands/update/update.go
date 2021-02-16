package update

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/git"
	configInstaller "github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/installer"
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
	repository.DeploymentStrategy = common.ConvStrToDeploymentStrategy(flags.DeploymentStrategy)
	repository.Recursive = flags.Recursive
	return repository
}

func updateImages(repo *git.Repository) {
	cdToAppDir(repo.Repo, repo.AppFlags.Path)
	input := getInstallInput(repo)
	installer := configInstaller.Installer{Command: common.UpdateCmd, KeepRegistry: repo.KeepRegistry, DeploymentStrategy: repo.DeploymentStrategy}
	installer.Install(input)
}

func resetAndUpdate(workingRepo *git.Repository) {
	common.RmGitOpsWd()
	cloneOrPull(workingRepo)
	updateImages(workingRepo)
	applyChanges(workingRepo)
	pushOnDemand(workingRepo)
}

func logEndMessage(workingRepo *git.Repository) {
	if workingRepo.GitFlags.Push {
		url := common.GetRemoteRepoWebUrl(workingRepo.Repo)
		common.Logger.Println(fmt.Sprintf("repository url: %s", url))
	}
}
