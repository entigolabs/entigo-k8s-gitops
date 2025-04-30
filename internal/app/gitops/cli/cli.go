package cli

import (
	"context"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/urfave/cli/v3"
	"os"
	"strings"
)

var flags *common.Flags = new(common.Flags)

func Run(ctx context.Context) {
	app := &cli.Command{Commands: cliCommands()}
	addAppInfo(app)
	loggingLvl := getLoggingLvl()
	common.ChooseLogger(loggingLvl)
	err := app.Run(ctx, os.Args)
	if err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
}

func addAppInfo(app *cli.Command) {
	const gitOps = "gitops"
	app.Name = gitOps
	app.Usage = "helper utility"
}

func getLoggingLvl() common.LoggingLevel {
	for i, arg := range os.Args {
		if isLoggerFlag(arg) {
			return getLoggerFlagValue(i, arg)
		}
	}
	return common.ProdLoggingLvl
}

func isLoggerFlag(arg string) bool {
	isLongLoggingFlag := strings.Contains(arg, "--logging=") || arg == "--logging"
	isShortLoggingFlag := strings.Contains(arg, "-l=") || arg == "-l"
	return isLongLoggingFlag || isShortLoggingFlag
}

func getLoggerFlagValue(index int, loggerArg string) common.LoggingLevel {
	if strings.Contains(loggerArg, "=") {
		splits := strings.Split(loggerArg, "=")
		loggingLvlAsString := strings.TrimSpace(splits[len(splits)-1])
		return common.ConvStrToLoggingLvl(loggingLvlAsString)
	} else {
		return common.ConvStrToLoggingLvl(os.Args[index+1])
	}
}
