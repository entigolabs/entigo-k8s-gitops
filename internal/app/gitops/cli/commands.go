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
	Usage:   "update corresponding images",
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
