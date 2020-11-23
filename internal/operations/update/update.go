package update

import (
	"errors"
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/logger"
)

const OperationType = "update"

func Update() func() {
	return func() {
		flgs = getFlags()
		logger.Logger.Println(&logger.Warning{Reason: errors.New(fmt.Sprint("Flags: ", flgs))})

		// mkdir gitops
		// cd gitops
		// git clone repo
		// change files
		// git add
		// git commit
		// git push
	}
}
