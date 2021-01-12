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
	EnvVars:     []string{"LOGGING"},
	DefaultText: "prod",
	Value:       "prod",
	Usage:       "set `logging level`",
	Destination: &flags.LoggingLevel,
}

var gitRepoFlag = cli.StringFlag{
	Name:        "git-repo",
	EnvVars:     []string{"GIT_REPO"},
	DefaultText: "",
	Usage:       "Git repository `SSH URL`",
	Destination: &flags.Git.Repo,
	Required:    true,
}

var gitBranchFlag = cli.StringFlag{
	Name:        "git-branch",
	EnvVars:     []string{"GIT_BRANCH"},
	DefaultText: "",
	Usage:       "branch `name`",
	Destination: &flags.Git.Branch,
	Required:    true,
}

var gitKeyFileFlag = cli.StringFlag{
	Name:        "git-key-file",
	EnvVars:     []string{"GIT_KEY_FILE"},
	DefaultText: "",
	Usage:       "SSH private key `location`",
	Destination: &flags.Git.KeyFile,
	Required:    true,
}

var gitStrictHostKeyCheckingFlag = cli.BoolFlag{
	Name:        "git-strict-host-key-checking",
	EnvVars:     []string{"GIT_STRICT_HOST_KEY_CHECKING"},
	DefaultText: strconv.FormatBool(false),
	Usage:       "strict host key checking",
	Destination: &flags.Git.StrictHostKeyChecking,
}

var gitPushFlag = cli.BoolFlag{
	Name:        "git-push",
	EnvVars:     []string{"GIT_PUSH"},
	DefaultText: strconv.FormatBool(true),
	Value:       true,
	Usage:       "push changes",
	Destination: &flags.Git.Push,
}

var gitAuthorNameFlag = cli.StringFlag{
	Name:        "git-author-name",
	EnvVars:     []string{"GIT_AUTHOR_NAME"},
	DefaultText: "jenkins",
	Usage:       "Git author name",
	Destination: &flags.Git.AuthorName,
}

var gitAuthorEmailFlag = cli.StringFlag{
	Name:        "git-author-email",
	EnvVars:     []string{"GIT_AUTHOR_EMAIL"},
	DefaultText: "jenkins@localhost",
	Value:       "jenkins@localhost",
	Usage:       "Git author email",
	Destination: &flags.Git.AuthorEmail,
}

var appPathFlag = cli.StringFlag{
	Name:        "app-path",
	EnvVars:     []string{"APP_PATH"},
	DefaultText: "",
	Usage:       "path to application folder",
	Destination: &flags.App.Path,
}

var appPrefixFlag = cli.StringFlag{
	Name:        "app-prefix",
	EnvVars:     []string{"APP_PREFIX"},
	DefaultText: "",
	Usage:       "`path` prefix to apply",
	Destination: &flags.App.Prefix,
}

var appNamespaceFlag = cli.StringFlag{
	Name:        "app-namespace",
	EnvVars:     []string{"APP_NAMESPACE"},
	DefaultText: "",
	Usage:       "application namespace `name`",
	Destination: &flags.App.Namespace,
}

var appNameFlag = cli.StringFlag{
	Name:        "app-name",
	EnvVars:     []string{"APP_NAME"},
	DefaultText: "",
	Usage:       "application name",
	Destination: &flags.App.Name,
}

var appBranchFlag = cli.StringFlag{
	Name:        "app-branch",
	EnvVars:     []string{"APP_BRANCH"},
	DefaultText: "",
	Usage:       "application branch `name`",
	Destination: &flags.App.Branch,
	Required:    true,
}

var appPrefixArgoFlag = cli.StringFlag{
	Name:        "app-prefix-argo",
	EnvVars:     []string{"APP_PREFIX_ARGO"},
	DefaultText: "",
	Usage:       "Argo app `path`",
	Destination: &flags.App.PrefixArgo,
}

var appPrefixYamlFlag = cli.StringFlag{
	Name:        "app-prefix-yaml",
	EnvVars:     []string{"APP_PREFIX_YAML"},
	DefaultText: "",
	Usage:       "yaml configurations `path`",
	Destination: &flags.App.PrefixYaml,
}

var imagesFlag = cli.StringFlag{
	Name:        "images",
	Aliases:     []string{"i"},
	EnvVars:     []string{"IMAGES"},
	DefaultText: "",
	Usage:       "images with tags",
	Destination: &flags.Images,
	Required:    true,
}

var keepRegistryFlag = cli.BoolFlag{
	Name:        "keep-registry",
	Aliases:     []string{"k"},
	EnvVars:     []string{"KEEP_REGISTRY"},
	DefaultText: strconv.FormatBool(false),
	Usage:       "Keeps registry part of the changeable image",
	Destination: &flags.KeepRegistry,
}
