package common

import (
	"fmt"
)

// todo refactor flags
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
	if f.App.PrefixYaml == "" {
		return fmt.Sprintf("%s/%s/%s", f.App.Prefix, f.App.Namespace, f.App.Name)
	}
	return fmt.Sprintf("%s/%s/%s/%s", f.App.Prefix, f.App.PrefixYaml, f.App.Namespace, f.App.Name)
}

func (f *Flags) ComposeArgoPath() string {
	if f.App.PrefixArgo == "" {
		return fmt.Sprintf("%s/%s/%s", f.App.Prefix, f.App.Namespace, f.App.Name)
	}
	return fmt.Sprintf("%s/%s/%s/%s", f.App.Prefix, f.App.PrefixArgo, f.App.Namespace, f.App.Name)
}

func (f *Flags) composeAppPath() {
	f.App.Path = fmt.Sprintf("%s/%s/%s", f.App.Prefix, f.App.Namespace, f.App.Name)
}

func (f *Flags) setup() {
	f.Git.Branch = sanitizeBranch(f.Git.Branch)
}

func (f *Flags) cmdSpecificSetup(cmd Command) {
	if cmd == UpdateCmd {
		if f.isTokenizedPath() {
			f.composeAppPath()
		}
	} else {
		f.App.Path = f.ComposeYamlPath()
	}
}
