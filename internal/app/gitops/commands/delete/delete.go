package delete

import (
	"errors"
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/git"
	"os"
)

func Run(flags *common.Flags) {
	repo := initWorkingRepo(flags)
	cloneOrPull(flags, repo)
	common.CdToRepoRoot(flags.Git.Repo)
	deleteAppBranch(flags)
	deleteArgoApp(flags)
	applyChanges(repo)
	pushOnDemand(flags, repo)
	logEndMessage(repo)
}

func initWorkingRepo(flags *common.Flags) *git.Repository {
	repository := new(git.Repository)
	repository.GitFlags = flags.Git
	repository.Images = flags.Images
	repository.AppFlags = flags.App
	repository.KeepRegistry = flags.KeepRegistry
	repository.Command = common.DeleteCmd
	repository.LoggingLevel = common.ConvStrToLoggingLvl(flags.LoggingLevel)
	return repository
}

func resetAndUpdate(flags *common.Flags, workingRepo *git.Repository) {
	common.RmGitOpsWd()
	cloneOrPull(flags, workingRepo)
	deleteAppBranch(flags)
	deleteArgoApp(flags)
	applyChanges(workingRepo)
	pushOnDemand(flags, workingRepo)
}

func deleteArgoApp(flags *common.Flags) {
	argoAppConfPath := fmt.Sprintf("%s/%s.yaml", flags.ComposeArgoPath(), flags.App.Branch)
	if _, err := os.Stat(argoAppConfPath); os.IsNotExist(err) {
		msg := fmt.Sprintf("skiping delation of %s.yaml - %s", flags.App.Branch, err)
		common.Logger.Println(&common.Warning{Reason: errors.New(msg)})
	} else {
		if err := os.RemoveAll(argoAppConfPath); err != nil {
			common.Logger.Println(&common.PrefixedError{Reason: err})
		} else {
			common.Logger.Println(fmt.Sprintf("deleted %s", argoAppConfPath))
		}
	}
}

func deleteAppBranch(flags *common.Flags) {
	appBranchPath := fmt.Sprintf("%s/%s", flags.ComposeYamlPath(), flags.App.Branch)
	if _, err := os.Stat(appBranchPath); os.IsNotExist(err) {
		msg := fmt.Sprintf("skiping delation of %s - %s", flags.App.Branch, err)
		common.Logger.Println(&common.Warning{Reason: errors.New(msg)})
	} else {
		if err := os.RemoveAll(appBranchPath); err != nil {
			common.Logger.Println(&common.PrefixedError{Reason: err})
		} else {
			common.Logger.Println(fmt.Sprintf("deleted %s", appBranchPath))

		}
	}
}

func logEndMessage(workingRepo *git.Repository) {
	if workingRepo.GitFlags.Push {
		url := common.GetRemoteRepoWebUrl(workingRepo.Repo)
		common.Logger.Println(fmt.Sprintf("repository url: %s", url))
	}
}
