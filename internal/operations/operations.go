package operations

import (
	"errors"
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/logger"
	"github.com/entigolabs/entigo-k8s-gitops/internal/operations/copy"
	"github.com/entigolabs/entigo-k8s-gitops/internal/operations/update"
	"os"
)

func Operate() {
	operation := chooseOperation()
	executeOperation(operation)
}

func chooseOperation() func() {
	operationArg := os.Args[1:][0]
	switch operationArg {
	case update.OperationType:
		return update.Update()
	case copy.OperationType:
		return copy.Copy()
	default:
		message := fmt.Sprintf("Unsupported operation: %s", operationArg)
		logger.Logger.Println(&logger.PrefixedError{Reason: errors.New(message)})
		os.Exit(1)
	}
	return nil
}

func executeOperation(operation func()) {
	operation()
}
