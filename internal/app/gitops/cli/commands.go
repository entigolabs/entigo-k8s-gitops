package cli

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/copy"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/update"
	"github.com/urfave/cli/v2"
)

var updateCommand = cli.Command{
	Name:    "update",
	Aliases: []string{"up"},
	Usage:   "Update corresponding images",
	Action: func(c *cli.Context) error {
		fmt.Println("update!")
		update.Run(Flags)
		return nil
	},
	Flags: cliFlags(),
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
	Flags: cliFlags(),
}
