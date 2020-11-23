package copy

import (
	"errors"
	"github.com/entigolabs/entigo-k8s-gitops/internal/util"
	"os"
)

const OperationType = "copy"

func Copy() func() {
	return func() {
		// TODO implement multibranch logic
		util.Logger.Println(&util.PrefixedError{Reason: errors.New("operation 'copy' not yet implemented")})
		os.Exit(1)
	}
}
