package update

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/git"
	configInstaller "github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/installer"
	"io/ioutil"
	"strings"
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
	images := strings.Split(workingRepo.Images, ",")
	installer := configInstaller.Installer{Command: common.UpdateCmd, AppBranch: "placeholder", AppName: "placeholder", KeepRegistry: workingRepo.KeepRegistry}
	installInput := getInstallInput(images)
	installer.Install(installInput)
}

func getInstallInput(images []string) string {
	yamlNamesArray, err := readDirFiltered(".", ".yaml")
	if err != nil {
		common.Logger.Println(common.PrefixedError{Reason: err})
	}
	return composeInstallInput(images, yamlNamesArray)
}

func composeInstallInput(images []string, yamlNamesArray []string) string {
	installInput := ""
	for _, image := range images {
		editLine := fmt.Sprintf("edit %s %s %s", strings.Join(yamlNamesArray, ","), getInstallLocations(), image)
		installInput += fmt.Sprintf("%s\n", editLine)
	}
	return installInput
}

func getInstallLocations() string {
	deploymentStatefulSetDaemonSetJobLocation := "spec.template.spec.containers.*.image"
	podLocation := "spec.containers.*.image"
	cronJobLocation := "spec.jobTemplate.spec.template.spec.containers.*.image"
	return strings.Join([]string{deploymentStatefulSetDaemonSetJobLocation, podLocation, cronJobLocation}, ",")
}

func readDirFiltered(path string, suffix string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, file := range files {
		if file.IsDir() == false && strings.HasSuffix(file.Name(), suffix) {
			result = append(result, file.Name())
		}
	}
	return result, err
}

func logEndMessage(workingRepo *git.Repository) {
	if workingRepo.GitFlags.Push {
		url := common.GetRemoteRepoWebUrl(workingRepo.Repo)
		common.Logger.Println(fmt.Sprintf("repository url: %s", url))
	}
}