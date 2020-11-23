package main

import (
	"github.com/entigolabs/entigo-k8s-gitops/internal/operations"
	"github.com/entigolabs/entigo-k8s-gitops/internal/util"
)

func init() {
	util.Logger = util.ChooseLogger("prod")
}

func main() {
	operations.Operate()
}
