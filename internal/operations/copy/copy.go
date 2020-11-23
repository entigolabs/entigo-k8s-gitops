package copy

import (
	"errors"
	"github.com/entigolabs/entigo-k8s-gitops/internal/logger"
	"os"
)

const OperationType = "copy"

func Copy() func() {
	return func() {
		// TODO implement multibranch logic
		logger.Logger.Println(&logger.PrefixedError{Reason: errors.New("operation 'copy' not yet implemented")})
		os.Exit(1)
	}
}
