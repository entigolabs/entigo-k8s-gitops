package util

import (
	"fmt"
	"os"
)

func cdToDeployment() {
	path := fmt.Sprintf("%s/%s", getwd(), deploymentDir)
	if _, err := os.Stat(path); err != nil {
		os.Mkdir(deploymentDir, 0755)
	}
	if err := changeDir(path); err != nil {
		logger.Println(&prefixedError{err})
		os.Exit(1)
	}
}

func cdToWorkPath(vars deploymentVariables) {
	path := fmt.Sprintf("%s/%s", getwd(), vars.WorkPath)
	if err := changeDir(path); err != nil {
		logger.Println(&prefixedError{err})
		os.Exit(1)
	}
}
