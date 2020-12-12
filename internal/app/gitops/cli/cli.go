package cli

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/entigolabs/entigo-k8s-gitops/internal/util"
	"github.com/urfave/cli/v2"
	"os"
)

var Flags = new(common.Flags)

func Run() {
	app := &cli.App{
		Flags:    cliFlags(),
		Commands: cliCommands(),
		Action:   defaultAction(),
	}
	err := app.Run(os.Args)
	common.Logger = common.ChooseLogger(Flags.LoggingLevel)
	if err != nil {
		common.Logger.Fatal(&util.PrefixedError{Reason: err})
	}
}

func cliFlags() []cli.Flag {
	return []cli.Flag{
		&loggingFlag,
		&gitRepoFlag,
		&gitBranchFlag,
		&gitKeyFileFlag,
		&gitStrictHostKeyCheckingFlag,
		&gitPushFlag,
		&appPrefixFlag,
		&appNamespaceFlag,
		&appNameFlag,
		&appPathFlag,
		&imagesFlag,
		&keepRegistryFlag,
	}
}

func cliCommands() []*cli.Command {
	return []*cli.Command{
		&updateCommand,
		&copyCommand,
	}
}

func defaultAction() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		fmt.Println("defaultAction ->", Flags)
		//cli.ShowAppHelp(c) // manually call help
		return nil
	}
}
