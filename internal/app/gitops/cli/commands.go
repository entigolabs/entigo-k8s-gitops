package cli

import (
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
	Action:  action(runCmd),
	Flags:   cliFlags(),
}

var updateCommand = cli.Command{
	Name:    "update",
	Aliases: []string{"up"},
	Usage:   "update corresponding images",
	Action:  action(updateCmd),
	Flags:   cliFlags(),
}

var copyCommand = cli.Command{
	Name:    "copy",
	Aliases: []string{"cp"},
	Usage:   "copy from master branch",
	Action:  action(copyCmd),
	Flags:   cliFlags(),
}
