package cli

import (
	"github.com/urfave/cli/v2"
)

var runCommand = cli.Command{
	Name:    "run",
	Aliases: []string{"rn"},
	Usage:   "Copy and update",
	Action:  action(runCmd),
	Flags:   cliFlags(runCmd),
}

var updateCommand = cli.Command{
	Name:    "update",
	Aliases: []string{"up"},
	Usage:   "Update corresponding images",
	Action:  action(updateCmd),
	Flags:   cliFlags(updateCmd),
}

var copyCommand = cli.Command{
	Name:    "copy",
	Aliases: []string{"cp"},
	Usage:   "Copy from master branch",
	Action:  action(copyCmd),
	Flags:   cliFlags(copyCmd),
}
