package cli

import (
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/urfave/cli/v2"
)

func deploymentNotificationFlags(cmd common.Command) []cli.Flag {
	var flags []cli.Flag
	flags = appendNotificationFlags(flags)
	flags = appendNotificationCmdFlags(flags, cmd)
	return flags
}

func appendNotificationFlags(flags []cli.Flag) []cli.Flag {
	return append(flags,
		&loggingFlag)
}

func appendNotificationCmdFlags(baseFlags []cli.Flag, cmd common.Command) []cli.Flag {
	switch cmd {
	case common.DeploymentNotificationCmd:
		baseFlags = notificationSpecificFlags(baseFlags)
	}
	return baseFlags
}

func notificationSpecificFlags(baseFlags []cli.Flag) []cli.Flag {
	baseFlags = append(baseFlags, &notificationUrlFlag)
	baseFlags = append(baseFlags, &notificationDeploymentEnvNameFlag)
	baseFlags = append(baseFlags, &oldImgWithTagFlag)
	baseFlags = append(baseFlags, &newImgWithTagFlag)
	baseFlags = append(baseFlags, &registryUriFlag)
	baseFlags = append(baseFlags, &authenticationTokenFlag)
	return baseFlags
}

var notificationUrlFlag = cli.StringFlag{
	Name:        "url",
	EnvVars:     []string{"NOTIFICATION_URL"},
	DefaultText: "",
	Usage:       "URL where to post notification",
	Destination: &flags.Notification.URL,
	Required:    true,
}

var notificationDeploymentEnvNameFlag = cli.StringFlag{
	Name:        "deployment-env",
	EnvVars:     []string{"NOTIFICATION_DEPLOYMENT_ENV"},
	DefaultText: "",
	Usage:       "environment name where the image has been deployed",
	Destination: &flags.Notification.Environment,
	Required:    true,
}

var oldImgWithTagFlag = cli.StringFlag{
	Name:        "old-img",
	EnvVars:     []string{"NOTIFICATION_IMG_OLD"},
	DefaultText: "",
	Usage:       "old image value",
	Destination: &flags.Notification.OldImage,
	Required:    true,
}

var newImgWithTagFlag = cli.StringFlag{
	Name:        "new-img",
	EnvVars:     []string{"NOTIFICATION_IMG_NEW"},
	DefaultText: "",
	Usage:       "new image value",
	Destination: &flags.Notification.NewImage,
	Required:    true,
}

var registryUriFlag = cli.StringFlag{
	Name:        "registry-uri",
	EnvVars:     []string{"NOTIFICATION_REGISTRY_URI"},
	DefaultText: "",
	Usage:       "docker registry URI",
	Destination: &flags.Notification.RegistryUri,
	Required:    true,
}

var authenticationTokenFlag = cli.StringFlag{
	Name:        "auth-token",
	EnvVars:     []string{"NOTIFICATION_AUTH_TOKEN"},
	DefaultText: "",
	Usage:       "authentication token as key and value pair ('key=value')",
	Destination: &flags.Notification.AuthToken,
	Required:    true,
}
