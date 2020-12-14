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

func (f *Flags) ComposeAppPath() {
	f.App.Path = fmt.Sprintf("composed AppPath with %s", f.Git.Repo)
}
