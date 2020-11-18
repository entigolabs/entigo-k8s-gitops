package main

import (
	"flag"
	"github.com/entigolabs/entigo-k8s-gitops/internal/util"
)

func main() {
	flags := getFlags()
	util.DeployArgo(flags)
}

func getFlags() util.Flags {
	k8sNsFlag, ImageNameFlag, argoAppFlag, branchFlag, sshKeyFlag, deployEnvFlag := parseFlags()

	return util.Flags{
		K8sNamespace:      *k8sNsFlag,
		ImageName:         *ImageNameFlag,
		ArgoAppName:       *argoAppFlag,
		GitBranch:         *branchFlag,
		SshKey:            *sshKeyFlag,
		DeploymentEnvName: *deployEnvFlag,
	}
}

func parseFlags() (*string, *string, *string, *string, *string, *string) {
	k8sNsFlag := flag.String(util.Flag_K8sNs, "", "Kubernetes Namespace")
	ImageNameFlag := flag.String(util.Flag_ImageName, "", "Kubernetes application name")
	argoAppFlag := flag.String(util.Flag_ArgoApp, "", "Argo CD application name")
	branchFlag := flag.String(util.Flag_Branch, "", "Git branch")
	sshKeyFlag := flag.String(util.Flag_SshKey, "", "SSH key")
	deployEnvFlag := flag.String(util.Flag_DeployEnv, "dev", "Deployment environment name")
	flag.Parse()
	return k8sNsFlag, ImageNameFlag, argoAppFlag, branchFlag, sshKeyFlag, deployEnvFlag
}
