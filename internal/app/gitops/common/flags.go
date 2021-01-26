package common

import (
	"errors"
	"fmt"
)

type Flags struct {
	LoggingLevel string
	Git          GitFlags
	App          AppFlags
	Images       string
	KeepRegistry bool
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
	if f.App.PrefixYaml != "" {
		yamlPath = yamlPath + appendSlash(f.App.PrefixYaml)
	} else {
		defaultPrefixYaml := "yaml"
		yamlPath = yamlPath + appendSlash(defaultPrefixYaml)
	}
	return yamlPath + appendSlash(f.App.Namespace) + f.App.Name
}

func (f *Flags) ComposeArgoPath() string {
	yamlPath := ""
	if f.App.Prefix != "" {
		yamlPath = appendSlash(f.App.Prefix)
	}
	if f.App.PrefixArgo != "" {
		yamlPath = yamlPath + appendSlash(f.App.PrefixArgo)
	} else {
		defaultPrefixArgo := "argoapps"
		yamlPath = yamlPath + appendSlash(defaultPrefixArgo)
	}
	return yamlPath + appendSlash(f.App.Namespace) + f.App.Name
}

func (f *Flags) composeAppPath() {
	f.App.Path = fmt.Sprintf("%s/%s/%s", f.App.Prefix, f.App.Namespace, f.App.Name)
}

func (f *Flags) setup() {
	f.Git.Branch = sanitizeBranch(f.Git.Branch)
}

func (f *Flags) cmdSpecificSetup(cmd Command) {
	switch cmd {
	case UpdateCmd:
		if f.isTokenizedPath() {
			f.composeAppPath()
		}
	case CopyCmd:
		f.App.Path = fmt.Sprintf("%s/%s/%s", f.App.Prefix, f.App.Namespace, f.App.Name)
	default:
		Logger.Fatal(&PrefixedError{Reason: errors.New("unsupported command")})

	}
}

func appendSlash(str string) string {
	return str + "/"
}
