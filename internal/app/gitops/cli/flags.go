package cli

import (
	"github.com/urfave/cli/v2"
	"strconv"
)

var loggingFlag = cli.StringFlag{
	Name:        "logging",
	Aliases:     []string{"l"},
	DefaultText: "prod",
	Value:       "prod",
	Usage:       "set `logging level`",
	Destination: &Flags.LoggingLevel,
}

var gitRepoFlag = cli.StringFlag{
	Name:        "repo",
	Aliases:     []string{"r"},
	DefaultText: "",
	Usage:       "Git repository `SSH URL`",
	Destination: &Flags.Git.Repo,
	Required:    true,
}

var gitBranchFlag = cli.StringFlag{
	Name:        "branch",
	Aliases:     []string{"b"},
	DefaultText: "",
	Usage:       "branch `name`",
	Destination: &Flags.Git.Branch,
	Required:    true,
}

var gitKeyFileFlag = cli.StringFlag{
	Name:        "key-file",
	Aliases:     []string{"k"},
	DefaultText: "",
	Usage:       "SSH private key `location`",
	Destination: &Flags.Git.KeyFile,
	Required:    true,
}

var gitStrictHostKeyCheckingFlag = cli.BoolFlag{
	Name:        "strict-host-key-checking",
	Aliases:     []string{"s"},
	DefaultText: strconv.FormatBool(false),
	Usage:       "strict host key checking",
	Destination: &Flags.Git.StrictHostKeyChecking,
}

var gitPushFlag = cli.BoolFlag{
	Name:        "push",
	Aliases:     []string{"p"},
	DefaultText: strconv.FormatBool(true),
	Usage:       "push changes",
	Destination: &Flags.Git.Push,
}

var appPathFlag = cli.StringFlag{
	Name:        "path",
	Aliases:     []string{"t"},
	DefaultText: "",
	Usage:       "path to application folder",
	Destination: &Flags.App.Path,
}

var appPrefixFlag = cli.StringFlag{
	Name:        "prefix",
	Aliases:     []string{"x"},
	DefaultText: "",
	Usage:       "`path` prefix to apply",
	Destination: &Flags.App.Prefix,
}

var appNamespaceFlag = cli.StringFlag{
	Name:        "namespace",
	Aliases:     []string{"c"},
	DefaultText: "",
	Usage:       "application namespace `name`",
	Destination: &Flags.App.Namespace,
}

var appNameFlag = cli.StringFlag{
	Name:        "name",
	Aliases:     []string{"n"},
	DefaultText: "",
	Usage:       "application name",
	Destination: &Flags.App.Name,
}

var imagesFlag = cli.StringFlag{
	Name:        "images",
	Aliases:     []string{"i"},
	DefaultText: "",
	Usage:       "images with tags",
	Destination: &Flags.Images,
	Required:    true,
}

var keepRegistryFlag = cli.BoolFlag{
	Name:        "keep-registry",
	Aliases:     []string{"g"},
	DefaultText: strconv.FormatBool(false),
	Usage:       "Keeps registry part of the changeable image",
	Destination: &Flags.KeepRegistry,
}
