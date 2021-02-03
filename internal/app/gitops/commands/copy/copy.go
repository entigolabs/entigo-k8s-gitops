package copy

import (
	"fmt"
	configInstaller "github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/commands/copy/installer"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/git"
	"github.com/otiai10/copy"
	"strings"
)

const installFile = "install.txt"

func Run(flags *common.Flags) {
	repo := initWorkingRepo(flags)
	cloneOrPull(flags, repo)
	copyMasterToNewBranch(flags)
	installViaFile(flags)
	installArgoApp(flags)
	applyChanges(repo)
	pushOnDemand(flags, repo)
	logEndMessage(repo)
}

func installArgoApp(flags *common.Flags) { // todo refactor
	common.CdToRepoRoot(flags.Git.Repo)
	cdToArgoApp(flags.ComposeArgoPath())
	if err := copy.Copy("master.yaml", fmt.Sprintf("%s.yaml", flags.App.Branch)); err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
	installer := configInstaller.Installer{AppBranch: flags.App.Branch, AppName: flags.App.Name}
	installInput := getArgoAppInstallInput(flags)
	installer.Install(installInput)
}

func getArgoAppInstallInput(flags *common.Flags) string {
	destinationDir := fmt.Sprintf("%s/%s", flags.ComposeYamlPath(), flags.App.Branch)
	editName := fmt.Sprintf("edit %s.yaml metadata.name %s-%s", flags.App.Branch, flags.App.Name, flags.App.Branch)
	editPath := fmt.Sprintf("edit %s.yaml spec.source.path %s", flags.App.Branch, destinationDir)
	editNamespace := fmt.Sprintf("edit %s.yaml spec.destination.namespace %s", flags.App.Branch, flags.App.Namespace)
	return strings.Join([]string{editName, editPath, editNamespace}, "\n")
}

func copyMasterToNewBranch(flags *common.Flags) {
	common.CdToRepoRoot(flags.Git.Repo)
	sourceDir := fmt.Sprintf("%s/master", flags.ComposeYamlPath())
	destinationDir := fmt.Sprintf("%s/%s", flags.ComposeYamlPath(), flags.App.Branch)
	if err := copy.Copy(sourceDir, destinationDir); err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
	cdToCopiedBranch(destinationDir)
}

func installViaFile(flags *common.Flags) {
	installer := configInstaller.Installer{AppBranch: flags.App.Branch, AppName: flags.App.Name}
	installInput := string(common.GetFileInput(installFile))
	installer.Install(installInput)
}

func initWorkingRepo(flags *common.Flags) *git.Repository {
	repository := new(git.Repository)
	repository.GitFlags = flags.Git
	repository.Images = flags.Images
	repository.AppFlags = flags.App
	repository.KeepRegistry = flags.KeepRegistry
	repository.Command = common.CopyCmd
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

func logEndMessage(workingRepo *git.Repository) {
	if workingRepo.GitFlags.Push {
		url := common.GetRemoteRepoWebUrl(workingRepo.Repo)
		common.Logger.Println(fmt.Sprintf("repository url: %s", url))
	}
}

func resetAndUpdate(flags *common.Flags, workingRepo *git.Repository) {
	common.RmGitOpsWd()
	cloneOrPull(flags, workingRepo)
	copyMasterToNewBranch(flags)
	installViaFile(flags)
	installArgoApp(flags)
	applyChanges(workingRepo)
	pushOnDemand(flags, workingRepo)
}
