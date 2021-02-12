package update

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/git"
	configInstaller "github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/installer"
	"strings"
)

func Run(flags *common.Flags) {
	repo := initWorkingRepo(flags)
	cloneOrPull(repo)
	updateImages(repo)
	updateDeploymentStrategy(repo)
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
	return repository
}

func updateImages(workingRepo *git.Repository) {
	cdToAppDir(workingRepo.Repo, workingRepo.AppFlags.Path)
	images := strings.Split(workingRepo.Images, ",")
	imgUpdateInput := installInput{installType: imageUpdateInstType, changeData: images}
	installer := configInstaller.Installer{Command: common.UpdateCmd, KeepRegistry: workingRepo.KeepRegistry, DeploymentStrategy: common.UnspecifiedStrategy}
	input := getInstallInput(imgUpdateInput)
	installer.Install(input)
}

func updateDeploymentStrategy(repo *git.Repository) {
	if repo.DeploymentStrategy == common.UnspecifiedStrategy {
		return
	}
	cdToAppDir(repo.Repo, repo.AppFlags.Path)
	strategyAsStr := common.ConvDeploymentStrategyToStr(repo.DeploymentStrategy)
	strategyUpdateInput := installInput{installType: strategyUpdateInstType, changeData: []string{strategyAsStr}}
	installer := configInstaller.Installer{Command: common.UpdateCmd, DeploymentStrategy: repo.DeploymentStrategy}
	input := getInstallInput(strategyUpdateInput)
	installer.Install(input)
}

func resetAndUpdate(workingRepo *git.Repository) {
	common.RmGitOpsWd()
	cloneOrPull(workingRepo)
	updateImages(workingRepo)
	updateDeploymentStrategy(workingRepo)
	applyChanges(workingRepo)
	pushOnDemand(workingRepo)
}

func logEndMessage(workingRepo *git.Repository) {
	if workingRepo.GitFlags.Push {
		url := common.GetRemoteRepoWebUrl(workingRepo.Repo)
		common.Logger.Println(fmt.Sprintf("repository url: %s", url))
	}
}
