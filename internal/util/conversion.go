package util

import (
	"errors"
	"os"
)

func doConvertStrToDeployEnvName(str string) (deployEnvName, error) {
	switch str {
	case string(devEntigo):
		return devEntigo, nil
	case string(devDatel):
		return devDatel, nil
	case string(testKK): // TODO rm - testing purposes only
		return testKK, nil
	case string(testDatel):
		return testDatel, nil
	case string(demoDatel):
		return demoDatel, nil
	}
	return "", &prefixedError{errors.New("Could not parse string to DeployEnvName")}
}

func convertStrToDeployEnvName(str string) deployEnvName {
	var result deployEnvName
	if conversionRes, err := doConvertStrToDeployEnvName(str); err != nil {
		logger.Println(&prefixedError{err})
		os.Exit(1)
	} else {
		result = conversionRes
	}
	return result
}
