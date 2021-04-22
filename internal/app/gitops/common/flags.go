package common

import (
	"errors"
	"fmt"
	"strings"
)

type Flags struct {
	LoggingLevel       string
	Git                GitFlags
	App                AppFlags
	ArgoCD             ArgoCDFlags
	Images             string
	KeepRegistry       bool
	DeploymentStrategy string
	Recursive          bool
}

type GitFlags struct {
	Repo                  string
	Branch                string
	KeyFile               string
	AuthorName            string
	AuthorEmail           string
	StrictHostKeyChecking bool
	Push                  bool
}

type AppFlags struct {
	Path       string
	Prefix     string
	Namespace  string
	Name       string
	Branch     string
	PrefixArgo string
	PrefixYaml string
}

type ArgoCDFlags struct {
	ServerAddr string
	Insecure   bool
	AuthToken  string
	Timeout    int
	Async      bool
	Cascade    bool
}

func (f *Flags) Setup(cmd Command) error {
	if err := f.validate(cmd); err != nil {
		return err
	}
	f.setup()
	f.cmdSpecificSetup(cmd)
	return nil
}

func (f *Flags) ComposeYamlPath() string {
	yamlPath := ""
	if f.App.Prefix != "" {
		yamlPath = appendSlash(f.App.Prefix)
	}
	yamlPath = yamlPath + appendSlash(f.App.PrefixYaml)
	return yamlPath + appendSlash(f.App.Namespace) + f.App.Name
}

func (f *Flags) ComposeArgoPath() string {
	yamlPath := ""
	if f.App.Prefix != "" {
		yamlPath = appendSlash(f.App.Prefix)
	}
	yamlPath = yamlPath + appendSlash(f.App.PrefixArgo)
	return yamlPath + appendSlash(f.App.Namespace) + f.App.Name
}

func (f *Flags) composeAppPath() {
	f.App.Path = fmt.Sprintf("%s/%s/%s", f.App.Prefix, f.App.Namespace, f.App.Name)
}

func (f *Flags) setup() {
	f.App.Namespace = strings.ToLower(f.App.Namespace)
	f.App.Name = strings.ToLower(f.App.Name)
	f.App.Branch = sanitizeBranch(f.App.Branch)
}

func (f *Flags) cmdSpecificSetup(cmd Command) {
	switch cmd {
	case UpdateCmd:
		if f.isTokenizedPath() {
			f.composeAppPath()
		}
	case CopyCmd:
		f.App.Path = fmt.Sprintf("%s/%s/%s", f.App.Prefix, f.App.Namespace, f.App.Name)
	case DeleteCmd:
		f.App.Path = fmt.Sprintf("%s/%s/%s", f.App.Prefix, f.App.Namespace, f.App.Name)
	case ArgoCDSyncCmd:
	case ArgoCDDeleteCmd:
	default:
		Logger.Fatal(&PrefixedError{Reason: errors.New("unsupported command")})

	}
}

func appendSlash(str string) string {
	return str + "/"
}
