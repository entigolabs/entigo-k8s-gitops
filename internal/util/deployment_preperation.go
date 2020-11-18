package util

import (
	"fmt"
	"os"
)

func flagsSetup(flags Flags) deploymentVariables {
	validateFlags(flags)
	var deployVars deploymentVariables
	deployVars = assignFlags(flags, deployVars)
	deployVars = chooseDeploymentEnv(deployVars.DeploymentEnvNameFlag, deployVars)
	return deployVars
}

func repositorySetup(deployVars deploymentVariables) deploymentVariables {
	getRepository(deployVars.RepositoryUrl, deployVars.RepositoryBranch, deployVars.SshKeyFlag)
	deployVars = defineAdditionalVars(deployVars)

	validateRepository(deployVars)
	return deployVars
}

func defineAdditionalVars(deployVars deploymentVariables) deploymentVariables {
	deployVars = setYamlRootPath(deployVars)
	deployVars.ArgoRootPath = deployVars.RootPath + "/argoapps"
	deployVars = defineYamlPath(deployVars)
	deployVars = defineArgoPath(deployVars)
	deployVars = defineWorkPathAndWorkName(deployVars)
	return deployVars
}

func defineYamlPath(vars deploymentVariables) deploymentVariables {
	vars.YamlPath = fmt.Sprintf("%s/%s", vars.YamlRootPath, vars.ArgoAppNameFlag)
	if vars.K8sNamespaceFlag != "" {
		vars.YamlPath = fmt.Sprintf("%s/%s/%s", vars.YamlRootPath, vars.K8sNamespaceFlag, vars.ArgoAppNameFlag)
	}
	return vars
}

func defineArgoPath(vars deploymentVariables) deploymentVariables {
	vars.ArgoPath = fmt.Sprintf("%s/%s", vars.ArgoRootPath, vars.ArgoAppNameFlag)
	if vars.K8sNamespaceFlag != "" {
		vars.ArgoPath = fmt.Sprintf("%s/%s/%s", vars.ArgoRootPath, vars.K8sNamespaceFlag, vars.ArgoAppNameFlag)
	}
	return vars
}

func defineWorkPathAndWorkName(deployVars deploymentVariables) deploymentVariables {
	if deployVars.GitBranchFlag == "" {
		deployVars.WorkPath = deployVars.YamlPath
		deployVars.WorkName = deployVars.ArgoAppNameFlag
	} else {
		deployVars = defineWhenBranchHasName(deployVars)
	}
	return deployVars
}

func defineWhenBranchHasName(deployVars deploymentVariables) deploymentVariables {
	deployVars.FeatureBranch = composeFeatureBranchName(deployVars.GitBranchFlag)
	deployVars.WorkPath = deployVars.YamlPath + "/" + deployVars.FeatureBranch
	deployVars.MasterPath = deployVars.YamlPath + "/master"
	deployVars.ArgoWork = deployVars.ArgoPath + "/" + deployVars.FeatureBranch + ".yaml"
	deployVars.ArgoMaster = deployVars.ArgoPath + "/master.yaml"

	if deployVars.FeatureBranch == "master" {
		deployVars.WorkName = deployVars.ArgoAppNameFlag
	} else {
		deployVars.WorkName = deployVars.ArgoAppNameFlag + "-" + deployVars.FeatureBranch
	}
	return deployVars
}

func setYamlRootPath(deployVars deploymentVariables) deploymentVariables {
	path := deployVars.RootPath + "/yaml"
	if _, err := os.Stat(path); err == nil {
		deployVars.YamlRootPath = path
	} else {
		deployVars.YamlRootPath = deployVars.RootPath
	}
	return deployVars
}

func assignFlags(flags Flags, deployVars deploymentVariables) deploymentVariables {
	deployEnvName := convertStrToDeployEnvName(flags.DeploymentEnvName)
	deployVars.DeploymentEnvNameFlag = deployEnvName
	deployVars.K8sNamespaceFlag = flags.K8sNamespace
	deployVars.ImageNameFlag = flags.ImageName
	deployVars.GitBranchFlag = flags.GitBranch
	deployVars.SshKeyFlag = flags.SshKey

	if flags.ArgoAppName == "" {
		deployVars.ArgoAppNameFlag = flags.ImageName
	} else {
		deployVars.ArgoAppNameFlag = flags.ArgoAppName
	}

	return deployVars
}
