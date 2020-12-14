package main

import (
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/cli"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
)

func init() {
	common.Logger = common.ChooseLogger("prod")
}

func main() {
	//operations.Operate()
	cli.Run()
	//f := new(cli.Flags)
	//f.Repo = "repoVal"
	//f.ComposeAppPath()
	//copy.Run(f)
	//update.Run(f)
}
