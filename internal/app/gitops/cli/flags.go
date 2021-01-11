package cli

import (
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/urfave/cli/v2"
	"strconv"
)

func cliFlags(cmd common.Command) []cli.Flag {
	baseFlags := []cli.Flag{
		&loggingFlag,
		&gitRepoFlag,
		&gitBranchFlag,
		&gitKeyFileFlag,
		&gitStrictHostKeyCheckingFlag,
		&gitPushFlag,
		&gitAuthorNameFlag,
		&gitAuthorEmailFlag,
		&appPrefixFlag,
		&appNamespaceFlag,
		&appNameFlag,
	}
	baseFlags = appendCmdSpecificFlags(baseFlags, cmd)
	return baseFlags
}

func appendCmdSpecificFlags(baseFlags []cli.Flag, cmd common.Command) []cli.Flag {
	switch cmd {
	case common.RunCmd:
	case common.UpdateCmd:
		baseFlags = append(baseFlags, &imagesFlag)
		baseFlags = append(baseFlags, &keepRegistryFlag)
		baseFlags = append(baseFlags, &appPathFlag)
	case common.CopyCmd:
		baseFlags = append(baseFlags, &appBranchFlag)
		baseFlags = append(baseFlags, &appPrefixArgoFlag)
		baseFlags = append(baseFlags, &appPrefixYamlFlag)
	}
	return baseFlags
}

var loggingFlag = cli.StringFlag{
	Name:        "logging",
	Aliases:     []string{"l"},
	DefaultText: "prod",
	Value:       "prod",
	Usage:       "set `logging level`",
	Destination: &flags.LoggingLevel,
}

var gitRepoFlag = cli.StringFlag{
	Name:        "repo",
	Aliases:     []string{"r"},
	DefaultText: "",
	Usage:       "Git repository `SSH URL`",
	Destination: &flags.Git.Repo,
	Required:    true,
}

var gitBranchFlag = cli.StringFlag{
	Name:        "branch",
	Aliases:     []string{"b"},
	DefaultText: "",
	Usage:       "branch `name`",
	Destination: &flags.Git.Branch,
	Required:    true,
}

var gitKeyFileFlag = cli.StringFlag{
	Name:        "key-file",
	Aliases:     []string{"k"},
	DefaultText: "",
	Usage:       "SSH private key `location`",
	Destination: &flags.Git.KeyFile,
	Required:    true,
}

var gitStrictHostKeyCheckingFlag = cli.BoolFlag{
	Name:        "strict-host-key-checking",
	Aliases:     []string{"s"},
	DefaultText: strconv.FormatBool(false),
	Usage:       "strict host key checking",
	Destination: &flags.Git.StrictHostKeyChecking,
}

var gitPushFlag = cli.BoolFlag{
	Name:        "push",
	Aliases:     []string{"p"},
	DefaultText: strconv.FormatBool(true),
	Value:       true,
	Usage:       "push changes",
	Destination: &flags.Git.Push,
}

var gitAuthorNameFlag = cli.StringFlag{
	Name:        "git-author-name",
	Aliases:     []string{"ga"},
	DefaultText: "jenkins",
	Value:       "jenkins",
	Usage:       "Git author name",
	Destination: &flags.Git.AuthorName,
}

var gitAuthorEmailFlag = cli.StringFlag{
	Name:        "git-author-email",
	Aliases:     []string{"ge"},
	DefaultText: "jenkins@localhost",
	Value:       "jenkins@localhost",
	Usage:       "Git author email",
	Destination: &flags.Git.AuthorEmail,
}

var appPathFlag = cli.StringFlag{
	Name:        "path",
	Aliases:     []string{"t"},
	DefaultText: "",
	Usage:       "path to application folder",
	Destination: &flags.App.Path,
}

var appPrefixFlag = cli.StringFlag{
	Name:        "prefix",
	Aliases:     []string{"x"},
	DefaultText: "",
	Usage:       "`path` prefix to apply",
	Destination: &flags.App.Prefix,
}

var appNamespaceFlag = cli.StringFlag{
	Name:        "namespace",
	Aliases:     []string{"c"},
	DefaultText: "",
	Usage:       "application namespace `name`",
	Destination: &flags.App.Namespace,
}

var appNameFlag = cli.StringFlag{
	Name:        "name",
	Aliases:     []string{"n"},
	DefaultText: "",
	Usage:       "application name",
	Destination: &flags.App.Name,
}

var appBranchFlag = cli.StringFlag{
	Name:        "app-branch",
	Aliases:     []string{"a"},
	DefaultText: "",
	Usage:       "application branch `name`",
	Destination: &flags.App.Branch,
	Required:    true,
}

var appPrefixArgoFlag = cli.StringFlag{
	Name:        "prefix-argo",
	Aliases:     []string{"g"},
	DefaultText: "",
	Usage:       "Argo app `path`",
	Destination: &flags.App.PrefixArgo,
}

var appPrefixYamlFlag = cli.StringFlag{
	Name:        "prefix-yaml",
	Aliases:     []string{"y"},
	DefaultText: "",
	Usage:       "yaml configurations `path`",
	Destination: &flags.App.PrefixYaml,
}

var imagesFlag = cli.StringFlag{
	Name:        "images",
	Aliases:     []string{"i"},
	DefaultText: "",
	Usage:       "images with tags",
	Destination: &flags.Images,
	Required:    true,
}

var keepRegistryFlag = cli.BoolFlag{
	Name:        "keep-registry",
	Aliases:     []string{"g"},
	DefaultText: strconv.FormatBool(false),
	Usage:       "Keeps registry part of the changeable image",
	Destination: &flags.KeepRegistry,
}
