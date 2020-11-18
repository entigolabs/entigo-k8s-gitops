package util

import (
	"fmt"
	"os"
)

func validateRepository(vars deploymentVariables) {
	validateNamespace(vars)
	validateArgoAppPath(vars)
}

func validateNamespace(vars deploymentVariables) {
	nsPath := fmt.Sprintf("%s/%s", vars.YamlRootPath, vars.K8sNamespaceFlag)
	if _, err := os.Stat(nsPath); err != nil {
		logger.Println(&prefixedError{err})
		os.Exit(1)
	}
}

func validateArgoAppPath(vars deploymentVariables) {
	argoAppPath := vars.YamlPath
	if _, err := os.Stat(argoAppPath); err != nil {
		logger.Println(&prefixedError{err})
		os.Exit(1)
	}
}

func isMultibranchWorthy(vars deploymentVariables) bool {
	if _, err := os.Stat(vars.WorkPath); err != nil {
		return true
	}
	return false
}
