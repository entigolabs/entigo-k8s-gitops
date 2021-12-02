package copy

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/git"
	configInstaller "github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/installer"
	"github.com/otiai10/copy"
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

func resetAndUpdate(flags *common.Flags, workingRepo *git.Repository) {
	common.RmGitOpsWd()
	cloneOrPull(flags, workingRepo)
	copyMasterToNewBranch(flags)
	installViaFile(flags)
	installArgoApp(flags)
	applyChanges(workingRepo)
	pushOnDemand(flags, workingRepo)
}

func installArgoApp(flags *common.Flags) {
	common.CdToRepoRoot(flags.Git.Repo)
	cdToArgoApp(flags.ComposeArgoPath())
	if err := copy.Copy("master.yaml", fmt.Sprintf("%s.yaml", flags.App.Branch)); err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
	installer := configInstaller.Installer{Command: common.CopyCmd, DeploymentStrategy: common.UnspecifiedStrategy}
	input := composeArgoAppInstallInput(flags)
	installer.Install(input)
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
	installer := configInstaller.Installer{Command: common.CopyCmd, DeploymentStrategy: common.UnspecifiedStrategy}
	installTxt := string(common.GetFileInput(installFile))
	installTxtVars := *initInstallTxtVariables(flags.App.Branch, flags.App.Name, flags.App.Domain)
	installTxt = installTxtVars.specifyInstallVariables(installTxt)
	input := composeInstallInput(installTxt)
	installer.Install(input)
}

func logEndMessage(workingRepo *git.Repository) {
	if workingRepo.GitFlags.Push {
		url := common.GetRemoteRepoWebUrl(workingRepo.Repo)
		common.Logger.Println(fmt.Sprintf("repository url: %s", url))
	}
}
