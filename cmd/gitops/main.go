package main

import (
	"github.com/entigolabs/entigo-k8s-gitops/internal/logger"
	"github.com/entigolabs/entigo-k8s-gitops/internal/operations"
)

func init() {
	logger.Logger = logger.ChooseLogger("dev")
}

func main() {
	operations.Operate()
}
