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

func (f *Flags) Setup() error {
	if err := f.validatePaths(); err != nil {
		return err
	}
	f.composeAppPath()
	return nil
}

func (f *Flags) composeAppPath() {
	f.App.Path = fmt.Sprintf("%s/%s/%s", f.App.Prefix, f.App.Namespace, f.App.Name)
}
