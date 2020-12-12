package common

import "fmt"

type Flags struct {
	LoggingLevel string
	Repo         string
	Branch       string
	AppPath      string
}

func (f *Flags) ComposeAppPath() {
	f.AppPath = fmt.Sprintf("composed AppPath with %s", f.Repo)
}
