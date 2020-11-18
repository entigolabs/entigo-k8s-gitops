package util

import (
	"fmt"
	"runtime"
	"strings"
)

const (
	Flag_K8sNs     string = "k8s-ns"
	Flag_ImageName        = "img-name"
	Flag_ArgoApp          = "argo-app"
	Flag_Branch           = "branch"
	Flag_SshKey           = "ssh-key"
	Flag_DeployEnv        = "deploy-env"
)

func DeployArgo(flags Flags) {
	deployVars := flagsSetup(flags)
	cdToDeployment()
	deployVars = repositorySetup(deployVars)
	multibranchStep(deployVars)
	everyBranchStep(deployVars)
	fmt.Println(fmt.Sprintf("Repository url: %s", getWebUrl(deployVars.RepositoryUrl)))
}

func multibranchStep(deployVars deploymentVariables) {
	if !isMultibranchWorthy(deployVars) {
		return
	}
	fmt.Printf("WorkPath (%s) is suitable for multibranch step!\n", deployVars.WorkPath)
	// TODO implement multibranch logic
}

func everyBranchStep(vars deploymentVariables) {
	cdToWorkPath(vars)
	changeImages(vars)
	commit(vars)
	push(vars.SshKeyFlag)
}

func changeImages(vars deploymentVariables) {
	imgNameTokens := strings.Split(vars.ImageNameFlag, ",")
	for _, img := range imgNameTokens {
		changeSpecificImage(img)
	}
}

func changeSpecificImage(image string) {
	switch os := runtime.GOOS; os {
	case "darwin":
		changeImageOsx(image)
	default:
		changeImageDefault(image)
	}
}
