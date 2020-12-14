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
	app := &cli.App{Commands: cliCommands()}
	addAppInfo(app)
	err := app.Run(os.Args)
	common.Logger = common.ChooseLogger(Flags.LoggingLevel)
	if err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
}

func addAppInfo(app *cli.App) {
	const gitOps = "gitops"
	app.Name = gitOps
	app.HelpName = gitOps
	app.Usage = "helper utility"
}
