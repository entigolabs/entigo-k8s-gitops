package util

import (
	"os"
)

func validateFlags(flags Flags) bool {
	areParamsValid := doFlagsValidation(flags.ImageName, flags.DeploymentEnvName)
	if areParamsValid == false {
		os.Exit(1)
	}
	return areParamsValid
}

func doFlagsValidation(imageName string, deployEnvName string) bool {
	areParamsValid := true
	if err := validateImageName(imageName); err != nil {
		areParamsValid = false
	}
	if err := validateDeploymentEnvName(deployEnvName); err != nil {
		areParamsValid = false
	}
	return areParamsValid
}

func validateImageName(appName string) error {
	if appName == "" {
		err := &argumentError{Flag_ImageName, appName}
		logger.Println(&prefixedError{err})
		return err
	}
	return nil
}

func validateDeploymentEnvName(deployEnvName string) error {
	if _, err := isDeploymentEnvNameValid(deployEnvName); err != nil {
		logger.Println(&prefixedError{err})
		return err
	}
	return nil
}

func isDeploymentEnvNameValid(envName string) (bool, error) {
	switch envName {
	case string(devEntigo):
		return true, nil
	case string(devDatel):
		return true, nil
	case string(testKK): // TODO rm - testing purposes only
		return true, nil
	case string(testDatel):
		return true, nil
	case string(demoDatel):
		return true, nil
	default:
		return false, &argumentError{Flag_DeployEnv, envName}
	}
}
