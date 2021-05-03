package cli

import (
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/urfave/cli/v2"
)

func ArgoCDFlags(cmd common.Command) []cli.Flag {
	var flags []cli.Flag
	flags = appendArgoCDFlags(flags)
	flags = appendArgoCDCmdFlags(flags, cmd)
	return flags
}

func appendArgoCDFlags(flags []cli.Flag) []cli.Flag {
	return append(flags,
		&loggingFlag,
		&appNameFlag,
		&argoCDServerAddrFlag,
		&argoCDTokenFlag,
		&argoCDInsecureFlag,
		&argoCDTimeoutFlag,
		&argoCDGRPCWebFlag)
}

func appendArgoCDCmdFlags(baseFlags []cli.Flag, cmd common.Command) []cli.Flag {
	switch cmd {
	case common.ArgoCDSyncCmd:
		baseFlags = append(baseFlags,
			&argoCDAsyncFlag,
			&argoCDWaitFailureFlag)
	case common.ArgoCDDeleteCmd:
		baseFlags = append(baseFlags, &argoCDCascadeFlag)
	case common.ArgoCDUpdateCmd:
		baseFlags = argoCDUpdateSpecificFlags(baseFlags)
	}
	return baseFlags
}

func argoCDUpdateSpecificFlags(baseFlags []cli.Flag) []cli.Flag {
	return append(baseFlags,
		&gitKeyFileFlag,
		&imagesFlag,
		&gitStrictHostKeyCheckingFlag,
		&gitPushFlag,
		&gitAuthorNameFlag,
		&gitAuthorEmailFlag,
		&keepRegistryFlag,
		&deploymentStrategyFlag,
		&recursiveFlag,
		&argoCDAsyncFlag,
		&argoCDWaitFailureFlag)
}

var argoCDServerAddrFlag = cli.StringFlag{
	Name:        "server",
	EnvVars:     []string{"ARGO_CD_SERVER"},
	DefaultText: "",
	Usage:       "Server tcp address with port",
	Destination: &flags.ArgoCD.ServerAddr,
	Required:    true,
}

var argoCDInsecureFlag = cli.BoolFlag{
	Name:        "insecure",
	EnvVars:     []string{"ARGO_CD_INSECURE"},
	Value:       false,
	DefaultText: "false",
	Usage:       "Insecure connection",
	Destination: &flags.ArgoCD.Insecure,
	Required:    false,
}

var argoCDTokenFlag = cli.StringFlag{
	Name:        "auth-token",
	Aliases:     []string{"token"},
	EnvVars:     []string{"ARGO_CD_TOKEN"},
	DefaultText: "",
	Usage:       "Authentication token",
	Destination: &flags.ArgoCD.AuthToken,
	Required:    true,
}

var argoCDTimeoutFlag = cli.IntFlag{
	Name:        "timeout",
	EnvVars:     []string{"ARGO_CD_TIMEOUT"},
	Value:       300,
	DefaultText: "300",
	Usage:       "Timeout for single ArgoCD operation",
	Destination: &flags.ArgoCD.Timeout,
	Required:    false,
}

var argoCDGRPCWebFlag = cli.BoolFlag{
	Name:        "grpc-web",
	EnvVars:     []string{"ARGO_CD_GRPC_WEB"},
	Value:       false,
	DefaultText: "false",
	Usage:       "Enables gRPC-web protocol. Useful if Argo CD server is behind proxy which does not support HTTP2",
	Destination: &flags.ArgoCD.GRPCWeb,
	Required:    false,
}

var argoCDAsyncFlag = cli.BoolFlag{
	Name:        "async",
	EnvVars:     []string{"ARGO_CD_ASYNC"},
	Value:       false,
	DefaultText: "false",
	Usage:       "Don't wait for sync to complete",
	Destination: &flags.ArgoCD.Async,
	Required:    false,
}

var argoCDWaitFailureFlag = cli.BoolFlag{
	Name:        "wait-failure",
	EnvVars:     []string{"ARGO_CD_WAIT_FAILURE"},
	Value:       true,
	DefaultText: "true",
	Usage:       "Fail the command when waiting for the sync to complete exceeds the timeout",
	Destination: &flags.ArgoCD.WaitFailure,
	Required:    false,
}

var argoCDCascadeFlag = cli.BoolFlag{
	Name:        "cascade",
	EnvVars:     []string{"ARGO_CD_CASCADE"},
	Value:       true,
	DefaultText: "true",
	Usage:       "Perform a cascaded deletion of all application resources",
	Destination: &flags.ArgoCD.Async,
	Required:    false,
}
