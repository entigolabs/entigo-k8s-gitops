package copy

import (
	"errors"
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	copyInstaller "github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/copy/installer"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/git"
	"github.com/otiai10/copy"
)

func Run(flags *common.Flags) {
	if flags.LoggingLevel == "dev" { // TODO tm logging
		common.Logger.Println(&common.Warning{Reason: errors.New(fmt.Sprintf("copy:", flags))})
	}

	repo := initWorkingRepo(flags)
	cloneOrPull(repo)
	copyMasterToNewBranch(flags)
	installer := copyInstaller.Installer{GitBranch: flags.Git.Branch, AppName: flags.App.Name}
	installer.Install()
	// TODO impl install.txt related logic
}

func copyMasterToNewBranch(flags *common.Flags) {
	cdToRepoRoot(flags.Git.Repo)
	sourceDir := fmt.Sprintf("%s/master", flags.ComposeYamlPath())
	destinationDir := fmt.Sprintf("%s/%s", flags.ComposeYamlPath(), flags.Git.Branch)
	if err := copy.Copy(sourceDir, destinationDir); err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
	cdToCopiedBranch(destinationDir)
}

func initWorkingRepo(flags *common.Flags) *git.Repository {
	repository := new(git.Repository)
	repository.GitFlags = flags.Git
	repository.Images = flags.Images
	repository.AppPath = flags.App.Path
	repository.KeepRegistry = flags.KeepRegistry
	return repository
}

func cloneOrPull(workingRepo *git.Repository) {
	common.CdToGitOpsWd()
	if !workingRepo.DoesRepoExist() {
		cloneAndConfig(workingRepo)
	} else {
		workingRepo.OpenGitOpsRepo()
		if err := workingRepo.Pull(); err != nil {
			common.Logger.Fatal("workingRepo.Pull error") // TODO implement correct logic
			//resetAndUpdate(workingRepo)
		}
	}
}

func cloneAndConfig(workingRepo *git.Repository) {
	workingRepo.Clone()
	workingRepo.ConfigRepo()
}
