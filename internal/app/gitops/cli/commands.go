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
		&argoCDSyncCommand,
		&argoCDDeleteCommand,
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
	Usage:   "copy from master branch",
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

var argoCDSyncCommand = cli.Command{
	Name:    "argocd-sync",
	Aliases: []string{"a-sync"},
	Usage:   "synchronizes ArgoCD application",
	Action:  action(common.ArgoCDSyncCmd),
	Flags:   ArgoCDFlags(common.ArgoCDSyncCmd),
}

var argoCDDeleteCommand = cli.Command{
	Name:    "argocd-delete",
	Aliases: []string{"a-delete"},
	Usage:   "deletes ArgoCD application",
	Action:  action(common.ArgoCDDeleteCmd),
	Flags:   ArgoCDFlags(common.ArgoCDDeleteCmd),
}
