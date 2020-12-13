package cli

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/copy"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/update"
	"github.com/urfave/cli/v2"
)

var runCommand = cli.Command{
	Name:    "run",
	Aliases: []string{"rn"},
	Usage:   "Copy and update",
	Action: func(c *cli.Context) error {
		fmt.Println("run!")
		return nil
	},
	Flags: cliFlags(runCmd),
}

var updateCommand = cli.Command{
	Name:    "update",
	Aliases: []string{"up"},
	Usage:   "Update corresponding images",
	Action: func(c *cli.Context) error {
		fmt.Println("update!")
		update.Run(Flags)
		return nil
	},
	Flags: cliFlags(updateCmd),
}

var copyCommand = cli.Command{
	Name:    "copy",
	Aliases: []string{"cp"},
	Usage:   "Copy from master branch",
	Action: func(c *cli.Context) error {
		fmt.Println("copy!")
		copy.Run(Flags)
		return nil
	},
	Flags: cliFlags(copyCmd),
}
