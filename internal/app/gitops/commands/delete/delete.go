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

func cloneOrPull(flags *common.Flags, workingRepo *git.Repository) {
	common.CdToGitOpsWd()
	if !workingRepo.DoesRepoExist() {
		cloneAndConfig(workingRepo)
	} else {
		workingRepo.OpenGitOpsRepo()
		if err := workingRepo.Pull(); err != nil {
			resetAndUpdate(flags, workingRepo)
		}
	}
}

func cloneAndConfig(workingRepo *git.Repository) {
	workingRepo.Clone()
	workingRepo.ConfigRepo()
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

func applyChanges(workingRepo *git.Repository) {
	workingRepo.OpenGitOpsRepo()
	workingRepo.Add()
	workingRepo.CommitIfModified()
}

func pushOnDemand(flags *common.Flags, workingRepo *git.Repository) {
	if workingRepo.GitFlags.Push {
		if err := workingRepo.Push(); err != nil {
			resetAndUpdate(flags, workingRepo)
		}
	} else {
		common.Logger.Println("commit(s) were chosen not to be pushed")
	}
}

func resetAndUpdate(flags *common.Flags, workingRepo *git.Repository) {
	common.RmGitOpsWd()
	cloneOrPull(flags, workingRepo)
	deleteAppBranch(flags)
	deleteArgoApp(flags)
	applyChanges(workingRepo)
	pushOnDemand(flags, workingRepo)
}

func logEndMessage(workingRepo *git.Repository) {
	if workingRepo.GitFlags.Push {
		url := common.GetRemoteRepoWebUrl(workingRepo.Repo)
		common.Logger.Println(fmt.Sprintf("repository url: %s", url))
	}
}
