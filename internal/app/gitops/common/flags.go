package common

import "fmt"

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
	Path      string
	Prefix    string
	Namespace string
	Name      string
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
	return fmt.Sprintf("%s/%s/%s/%s", f.App.Prefix, "yaml", f.App.Namespace, f.App.Name)
}

func (f *Flags) composeAppPath() {
	f.App.Path = fmt.Sprintf("%s/%s/%s", f.App.Prefix, f.App.Namespace, f.App.Name)
}

func (f *Flags) setup() {
	f.sanitizeBranch()
}

func (f *Flags) sanitizeBranch() {
	// TODO getFeatureBranch(branch)
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
