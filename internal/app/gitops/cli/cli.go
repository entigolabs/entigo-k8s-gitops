package cli

import (
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/urfave/cli/v2"
	"os"
)

var Flags = new(common.Flags)

type Command int

const (
	runCmd Command = iota
	updateCmd
	copyCmd
)

func Run() {
	app := &cli.App{
		Commands: cliCommands(),
	}
	err := app.Run(os.Args)
	common.Logger = common.ChooseLogger(Flags.LoggingLevel)
	if err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
}

func cliFlags(cmd Command) []cli.Flag {
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
		&runCommand,
		&updateCommand,
		&copyCommand,
	}
}
