package cli

import (
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/urfave/cli/v2"
)

func cliCommands() []*cli.Command {
	return []*cli.Command{
		&runCommand,
		&updateCommand,
		&copyCommand,
		&deleteCommand,
		&notificationCommand,
		&argoCDGetCommand,
		&argoCDSyncCommand,
		&argoCDUpdateCommand,
		&argoCDDeleteCommand,
		&versionCommand,
	}
}

var runCommand = cli.Command{
	Name:    "run",
	Aliases: []string{"rn"},
	Usage:   "copy and update",
	Action:  action(common.RunCmd),
	Flags:   cliFlags(common.RunCmd),
}

var updateCommand = cli.Command{
	Name:    "update",
	Aliases: []string{"up"},
	Usage:   "update specified images",
	Action:  action(common.UpdateCmd),
	Flags:   cliFlags(common.UpdateCmd),
}

var copyCommand = cli.Command{
	Name:    "copy",
	Aliases: []string{"cp"},
	Usage:   "copy branch (default: master) and install changes from install.txt",
	Action:  action(common.CopyCmd),
	Flags:   cliFlags(common.CopyCmd),
}

var deleteCommand = cli.Command{
	Name:    "delete",
	Aliases: []string{"del"},
	Usage:   "deletes specified branch/folder and ArgoCD configuration",
	Action:  action(common.DeleteCmd),
	Flags:   cliFlags(common.DeleteCmd),
}

var notificationCommand = cli.Command{
	Name:    "notify",
	Aliases: []string{"ntf"},
	Usage:   "notifies API about new deployment",
	Action:  action(common.DeploymentNotificationCmd),
	Flags:   deploymentNotificationFlags(common.DeploymentNotificationCmd),
}

var argoCDGetCommand = cli.Command{
	Name:    "argocd-get",
	Aliases: []string{"a-get"},
	Usage:   "gets ArgoCD application",
	Action:  action(common.ArgoCDGetCmd),
	Flags:   ArgoCDFlags(common.ArgoCDGetCmd),
}

var argoCDSyncCommand = cli.Command{
	Name:    "argocd-sync",
	Aliases: []string{"a-sync"},
	Usage:   "synchronizes ArgoCD application",
	Action:  action(common.ArgoCDSyncCmd),
	Flags:   ArgoCDFlags(common.ArgoCDSyncCmd),
}

var argoCDUpdateCommand = cli.Command{
	Name:    "argocd-update",
	Aliases: []string{"a-update"},
	Usage:   "updates ArgoCD application",
	Action:  action(common.ArgoCDUpdateCmd),
	Flags:   ArgoCDFlags(common.ArgoCDUpdateCmd),
}

var argoCDDeleteCommand = cli.Command{
	Name:    "argocd-delete",
	Aliases: []string{"a-delete"},
	Usage:   "deletes ArgoCD application",
	Action:  action(common.ArgoCDDeleteCmd),
	Flags:   ArgoCDFlags(common.ArgoCDDeleteCmd),
}

var versionCommand = cli.Command{
	Name:    "version",
	Aliases: []string{"ver"},
	Usage:   "utility version",
	Action:  action(common.VersionCmd),
	Flags:   nil,
}
