package cli

import (
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/urfave/cli/v3"
	"strconv"
)

func cliFlags(cmd common.Command) []cli.Flag {
	var flags []cli.Flag
	flags = appendFlags(flags)
	flags = appendCmdFlags(flags, cmd)
	return flags
}

func appendFlags(flags []cli.Flag) []cli.Flag {
	return append(flags,
		&loggingFlag,
		&gitRepoFlag,
		&gitBranchFlag,
		&gitKeyFileFlag,
		&gitUsernameFlag,
		&gitPasswordFlag,
		&gitStrictHostKeyCheckingFlag,
		&gitPushFlag,
		&gitAuthorNameFlag,
		&gitAuthorEmailFlag,
		&appPrefixFlag,
		&appNamespaceFlag,
		&appNameFlag)
}

func appendCmdFlags(baseFlags []cli.Flag, cmd common.Command) []cli.Flag {
	switch cmd {
	case common.RunCmd:
	case common.UpdateCmd:
		baseFlags = updateSpecificFlags(baseFlags)
	case common.CopyCmd:
		baseFlags = copyAndDeleteSpecificFlags(baseFlags)
	case common.DeleteCmd:
		baseFlags = copyAndDeleteSpecificFlags(baseFlags)
	}
	return baseFlags
}

func updateSpecificFlags(baseFlags []cli.Flag) []cli.Flag {
	baseFlags = append(baseFlags, &imagesFlag)
	baseFlags = append(baseFlags, &keepRegistryFlag)
	baseFlags = append(baseFlags, &appPathFlag)
	baseFlags = append(baseFlags, &deploymentStrategyFlag)
	baseFlags = append(baseFlags, &recursiveFlag)
	baseFlags = append(baseFlags, &notifyApiUrlFlag)
	baseFlags = append(baseFlags, &notifyRegistryUriFlag)
	baseFlags = append(baseFlags, &notifyDevelopmentEnvFlag)
	baseFlags = append(baseFlags, &notifyAuthTokenFlag)
	return baseFlags
}

func copyAndDeleteSpecificFlags(baseFlags []cli.Flag) []cli.Flag {
	baseFlags = append(baseFlags, &appDestBranchFlag)
	baseFlags = append(baseFlags, &appPrefixArgoFlag)
	baseFlags = append(baseFlags, &appPrefixYamlFlag)
	baseFlags = append(baseFlags, &appDomainFlag)
	baseFlags = append(baseFlags, &appSourceBranchFlag)

	return baseFlags
}

var loggingFlag = cli.StringFlag{
	Name:        "logging",
	Aliases:     []string{"l"},
	Sources:     cli.EnvVars("LOGGING"),
	DefaultText: "prod",
	Value:       "prod",
	Usage:       "set `logging level` (prod | dev)",
	Destination: &flags.LoggingLevel,
}

var gitRepoFlag = cli.StringFlag{
	Name:        "git-repo",
	Sources:     cli.EnvVars("GIT_REPO"),
	DefaultText: "",
	Usage:       "Git repository `SSH URL`",
	Destination: &flags.Git.Repo,
	Required:    true,
}

var gitBranchFlag = cli.StringFlag{
	Name:        "git-branch",
	Sources:     cli.EnvVars("GIT_BRANCH"),
	DefaultText: "",
	Usage:       "branch `name`",
	Destination: &flags.Git.Branch,
	Required:    true,
}

var gitKeyFileFlag = cli.StringFlag{
	Name:        "git-key-file",
	Sources:     cli.EnvVars("GIT_KEY_FILE"),
	DefaultText: "",
	Usage:       "SSH private key `location`",
	Destination: &flags.Git.KeyFile,
	Required:    false,
}

var gitUsernameFlag = cli.StringFlag{
	Name:        "git-username",
	Sources:     cli.EnvVars("GIT_USERNAME"),
	DefaultText: "",
	Usage:       "git username",
	Destination: &flags.Git.Username,
	Required:    false,
}

var gitPasswordFlag = cli.StringFlag{
	Name:        "git-password",
	Sources:     cli.EnvVars("GIT_PASSWORD"),
	DefaultText: "",
	Usage:       "git password",
	Destination: &flags.Git.KeyFile,
	Required:    false,
}

var gitStrictHostKeyCheckingFlag = cli.BoolFlag{
	Name:        "git-strict-host-key-checking",
	Sources:     cli.EnvVars("GIT_STRICT_HOST_KEY_CHECKING"),
	DefaultText: strconv.FormatBool(false),
	Usage:       "strict host key checking",
	Destination: &flags.Git.StrictHostKeyChecking,
}

var gitPushFlag = cli.BoolFlag{
	Name:        "git-push",
	Sources:     cli.EnvVars("GIT_PUSH"),
	DefaultText: strconv.FormatBool(true),
	Value:       true,
	Usage:       "push changes",
	Destination: &flags.Git.Push,
}

var gitAuthorNameFlag = cli.StringFlag{
	Name:        "git-author-name",
	Sources:     cli.EnvVars("GIT_AUTHOR_NAME"),
	DefaultText: "jenkins",
	Value:       "jenkins",
	Usage:       "Git author name",
	Destination: &flags.Git.AuthorName,
}

var gitAuthorEmailFlag = cli.StringFlag{
	Name:        "git-author-email",
	Sources:     cli.EnvVars("GIT_AUTHOR_EMAIL"),
	DefaultText: "jenkins@localhost",
	Value:       "jenkins@localhost",
	Usage:       "Git author email",
	Destination: &flags.Git.AuthorEmail,
}

var appPathFlag = cli.StringFlag{
	Name:        "app-path",
	Sources:     cli.EnvVars("APP_PATH"),
	DefaultText: "",
	Usage:       "path to application folder",
	Destination: &flags.App.Path,
}

var appPrefixFlag = cli.StringFlag{
	Name:        "app-prefix",
	Sources:     cli.EnvVars("APP_PREFIX"),
	DefaultText: "",
	Usage:       "`path` prefix to apply",
	Destination: &flags.App.Prefix,
}

var appNamespaceFlag = cli.StringFlag{
	Name:        "app-namespace",
	Sources:     cli.EnvVars("APP_NAMESPACE"),
	DefaultText: "",
	Usage:       "application namespace `name`",
	Destination: &flags.App.Namespace,
}

var appNameFlag = cli.StringFlag{
	Name:        "app-name",
	Sources:     cli.EnvVars("APP_NAME"),
	DefaultText: "",
	Usage:       "application name",
	Destination: &flags.App.Name,
}

var appDestBranchFlag = cli.StringFlag{
	Name:        "app-dest-branch",
	Sources:     cli.EnvVars("APP_DEST_BRANCH"),
	DefaultText: "",
	Usage:       "application destination branch `name`",
	Destination: &flags.App.DestBranch,
	Required:    true,
}

var appDomainFlag = cli.StringFlag{
	Name:        "app-domain",
	Sources:     cli.EnvVars("APP_DOMAIN"),
	DefaultText: "localhost",
	Usage:       "application domain",
	Destination: &flags.App.Domain,
}

var appSourceBranchFlag = cli.StringFlag{
	Name:        "app-source-branch",
	Sources:     cli.EnvVars("APP_SOURCE_BRANCH"),
	DefaultText: "master",
	Value:       "master",
	Usage:       "application source branch `name`",
	Destination: &flags.App.SourceBranch,
}

var appPrefixArgoFlag = cli.StringFlag{
	Name:        "app-prefix-argo",
	Sources:     cli.EnvVars("APP_PREFIX_ARGO"),
	DefaultText: "argoapps",
	Value:       "argoapps",
	Usage:       "Argo app `path`",
	Destination: &flags.App.PrefixArgo,
}

var appPrefixYamlFlag = cli.StringFlag{
	Name:        "app-prefix-yaml",
	Sources:     cli.EnvVars("APP_PREFIX_YAML"),
	DefaultText: "yaml",
	Value:       "yaml",
	Usage:       "yaml configurations `path`",
	Destination: &flags.App.PrefixYaml,
}

var imagesFlag = cli.StringFlag{
	Name:        "images",
	Aliases:     []string{"i"},
	Sources:     cli.EnvVars("IMAGES"),
	DefaultText: "",
	Usage:       "images with tags",
	Destination: &flags.Images,
	Required:    true,
}

var keepRegistryFlag = cli.BoolFlag{
	Name:        "keep-registry",
	Aliases:     []string{"k"},
	Sources:     cli.EnvVars("KEEP_REGISTRY"),
	DefaultText: strconv.FormatBool(false),
	Usage:       "keeps registry part of the changeable image",
	Destination: &flags.KeepRegistry,
}

var deploymentStrategyFlag = cli.StringFlag{
	Name:        "deployment-strategy",
	Aliases:     []string{"s"},
	Sources:     cli.EnvVars("DEPLOYMENT-STRATEGY"),
	DefaultText: "if not defined then strategy will remain unchanged",
	Usage:       "change deployment strategy (RollingUpdate | Recreate)",
	Destination: &flags.DeploymentStrategy,
}

var recursiveFlag = cli.BoolFlag{
	Name:        "recursive",
	Sources:     cli.EnvVars("RECURSIVE"),
	DefaultText: strconv.FormatBool(false),
	Usage:       "updates directories and their contents recursively",
	Destination: &flags.Recursive,
}

var notifyApiUrlFlag = cli.StringFlag{
	Name:        "notify-api-url",
	Sources:     cli.EnvVars("NOTIFY_API_URL"),
	DefaultText: "",
	Usage:       "URL where to post notification",
	Destination: &flags.Notification.URL,
}

var notifyDevelopmentEnvFlag = cli.StringFlag{
	Name:        "notify-env",
	Sources:     cli.EnvVars("NOTIFY_ENV"),
	DefaultText: "",
	Usage:       "environment name where the image has been deployed",
	Destination: &flags.Notification.Environment,
}

var notifyRegistryUriFlag = cli.StringFlag{
	Name:        "notify-registry-uri",
	Sources:     cli.EnvVars("NOTIFY_REGISTRY_URI"),
	DefaultText: "",
	Usage:       "docker registry URI",
	Destination: &flags.Notification.RegistryUri,
}

var notifyAuthTokenFlag = cli.StringFlag{
	Name:        "notify-auth-token",
	Sources:     cli.EnvVars("NOTIFY_AUTH_TOKEN"),
	DefaultText: "",
	Usage:       "authentication token as key and value pair ('key=value')",
	Destination: &flags.Notification.AuthToken,
}
