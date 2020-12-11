package common

import "fmt"

type Flags struct {
	Repo    string
	Branch  string
	AppPath string
}

func (f *Flags) ComposeAppPath() {
	f.AppPath = fmt.Sprintf("composed AppPath with %s", f.Repo)
}
