package util

import "fmt"

type deployEnvName string

type deploymentVariables struct {
	DeploymentEnvNameFlag deployEnvName
	ImageNameFlag         string
	K8sNamespaceFlag      string
	GitBranchFlag         string
	SshKeyFlag            string
	FeatureBranch         string
	MasterPath            string
	WorkPath              string
	WorkName              string
	YamlRootPath          string
	YamlPath              string
	ArgoRootPath          string
	ArgoPath              string
	ArgoWork              string
	ArgoMaster            string
	ArgoAppNameFlag       string
	ArgoCDHost            string
	ArgoCDCredentials     string
	RepositoryUrl         string
	RepositoryBranch      string
	RepositoryCredentials string
	RootPath              string
	syncArgoApps          bool
}

type Flags = struct {
	K8sNamespace      string
	ImageName         string
	ArgoAppName       string
	GitBranch         string
	SshKey            string
	DeploymentEnvName string
}

type warning struct {
	reason error
}

func (w *warning) Error() string {
	warning := fmt.Sprintf("[warning] %s", w.reason)
	return fmt.Sprintf("\x1b[36;1m%s\x1b[0m", warning)
}

type prefixedError struct {
	reason error
}

func (pe *prefixedError) Error() string {
	error := fmt.Sprintf("[error] %s", pe.reason)
	return fmt.Sprintf("\x1b[31;1m%s\x1b[0m", error)
}

type argumentError struct {
	flagName string
	argument string
}

func (ae *argumentError) Error() string {
	error := fmt.Sprintf("Unsupported argument (%s) value: %s", ae.flagName, ae.argument)
	return fmt.Sprintf("\x1b[31;1m%s\x1b[0m", error)
}
