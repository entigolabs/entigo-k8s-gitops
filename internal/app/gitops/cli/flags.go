package cli

import (
	"github.com/urfave/cli/v2"
)

var loggingFlag = cli.StringFlag{
	Name:        "logging",
	Aliases:     []string{"l"},
	DefaultText: "prod",
	Value:       "prod",
	Usage:       "Set `logging level`",
	Destination: &Flags.LoggingLevel,
}

var repoFlag = cli.StringFlag{
	Name:        "repo",
	Aliases:     []string{"r"},
	DefaultText: "repoDefaultVal",
	Usage:       "Git repository `SSH URL`",
	Destination: &Flags.Repo,
}

var branchFlag = cli.StringFlag{
	Name:        "branch",
	Aliases:     []string{"b"},
	DefaultText: "branchDefaultVal",
	Usage:       "Branch `name`",
	Destination: &Flags.Branch,
}
