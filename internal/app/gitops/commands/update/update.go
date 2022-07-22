package update

import (
	"fmt"
	notifyCmd "github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/commands/notify"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/git"
	configInstaller "github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/installer"
	"strings"
)

func Run(flags *common.Flags) {
	repo := initWorkingRepo(flags)
	cloneOrPull(repo)
	updateImages(repo)
	applyChanges(repo)
	pushOnDemand(repo)
	logRepoUrl(repo)
	notifyOnDemand(repo)
}

func initWorkingRepo(flags *common.Flags) *git.Repository {
	repository := new(git.Repository)
	repository.GitFlags = flags.Git
	repository.Images = flags.Images
	repository.AppFlags = flags.App
	repository.DeploymentNotificationFlags = flags.Notification
	repository.KeepRegistry = flags.KeepRegistry
	repository.Command = common.UpdateCmd
	repository.LoggingLevel = common.ConvStrToLoggingLvl(flags.LoggingLevel)
	repository.DeploymentStrategy = common.ConvStrToDeploymentStrategy(flags.DeploymentStrategy)
	repository.Recursive = flags.Recursive
	repository.InstallHistory = []configInstaller.InstallHistory{}
	return repository
}

func updateImages(repo *git.Repository) {
	cdToAppDir(repo.Repo, repo.AppFlags.Path)
	input := composeInstallInput(strings.Split(repo.Images, ","), repo.Recursive)
	installer := configInstaller.Installer{Command: common.UpdateCmd, KeepRegistry: repo.KeepRegistry,
		DeploymentStrategy: repo.DeploymentStrategy, InstallHistory: repo.InstallHistory}
	repo.InstallHistory = installer.Install(input)
}

func resetAndUpdate(workingRepo *git.Repository) {
	common.RmGitOpsWd()
	cloneOrPull(workingRepo)
	updateImages(workingRepo)
	applyChanges(workingRepo)
	pushOnDemand(workingRepo)
}

func logRepoUrl(workingRepo *git.Repository) {
	if workingRepo.GitFlags.Push {
		url := common.GetRemoteRepoWebUrl(workingRepo.Repo)
		common.Logger.Println(fmt.Sprintf("repository url: %s", url))
	}
}

func notifyOnDemand(repo *git.Repository) {
	if err := common.ValidateNotificationFlags(repo.DeploymentNotificationFlags); err == nil {
		notifyChanges(repo)
	}
}

func notifyChanges(repo *git.Repository) {
	uniqueChanges := removeDuplicates(repo.InstallHistory)
	for _, historyRecord := range uniqueChanges {
		notify(repo, historyRecord)
	}
}

func notify(repo *git.Repository, historyRecord configInstaller.InstallHistory) {
	notificationFlags := common.DeploymentNotificationFlags{
		URL:         repo.DeploymentNotificationFlags.URL,
		Environment: repo.DeploymentNotificationFlags.Environment,
		OldImage:    historyRecord.OldValue,
		NewImage:    historyRecord.NewValue,
		RegistryUri: repo.DeploymentNotificationFlags.RegistryUri,
		AuthToken:   repo.DeploymentNotificationFlags.AuthToken,
	}
	notifyCmd.RunNotificationRequest(notificationFlags)
}

func removeDuplicates(historyRecord []configInstaller.InstallHistory) []configInstaller.InstallHistory {
	allKeys := make(map[string]bool)
	var list []configInstaller.InstallHistory
	for _, item := range historyRecord {
		// todo -> temporary: choose uniqueness by new value
		if _, value := allKeys[item.NewValue]; !value {
			allKeys[item.NewValue] = true
			list = append(list, item)
		}
	}
	return list
}
