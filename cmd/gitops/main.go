package main

import (
	flagsPkg "github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/entigolabs/entigo-k8s-gitops/internal/util"
)

func init() {
	util.Logger = util.ChooseLogger("prod")
}

var flags = new(flagsPkg.Flags)

func main() {
	//operations.Operate()
	runCli()

	//f := new(cli.Flags)
	//f.Repo = "repoVal"
	//f.ComposeAppPath()
	//copy.Run(f)
	//update.Run(f)
}
