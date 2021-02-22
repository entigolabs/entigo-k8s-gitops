package common

import (
	"errors"
	"fmt"
	"strings"
)

func ConvStrToLoggingLvl(str string) LoggingLevel {
	switch str {
	case "dev":
		return DevLoggingLvl
	case "prod":
		return ProdLoggingLvl
	default:
		msg := fmt.Sprintf("unsupported logging level: %s", str)
		Logger.Fatal(&PrefixedError{Reason: errors.New(msg)})
	}
	return ProdLoggingLvl
}

func ConvStrToDeploymentStrategy(str string) DeploymentStrategy {
	switch strings.ToLower(str) {
	case "rollingupdate":
		return RollingUpdateStrategy
	case "recreate":
		return RecreateStrategy
	case "":
		return UnspecifiedStrategy
	default:
		msg := fmt.Sprintf("unsupported deployment strategy: %s", str)
		Logger.Fatal(&PrefixedError{Reason: errors.New(msg)})
	}
	return UnspecifiedStrategy
}
