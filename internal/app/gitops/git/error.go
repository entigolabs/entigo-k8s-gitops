package git

import (
	"errors"
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/go-git/go-git/v5"
	"os"
	"strings"
)

var errUpdateReference = errors.New("failed to update ref")

func handlePullErr(err error) error {
	pullOp := "pull"
	if err == git.NoErrAlreadyUpToDate {
		alreadyUpToDateLogging(pullOp, err)
	} else if isConflictErr(err) {
		common.Logger.Println(fmt.Sprintf("couldn't git %s, %s", pullOp, err))
		return err
	} else {
		defaultGitOpErrLogging(pullOp, err)
		os.Exit(1)
	}
	return nil
}

func alreadyUpToDateLogging(gitOpName string, err error) {
	common.Logger.Println(fmt.Sprintf("skipping git %s, %s", gitOpName, err))
}

func isConflictErr(err error) bool {
	isErrNonFastForwardUpdate := strings.Contains(err.Error(), git.ErrNonFastForwardUpdate.Error()) ||
		err == git.ErrNonFastForwardUpdate
	isErrUpdateReference := strings.Contains(err.Error(), errUpdateReference.Error())
	return isErrNonFastForwardUpdate || isErrUpdateReference
}

func defaultGitOpErrLogging(gitOpName string, err error) {
	errorMessage := fmt.Sprintf("git %s failed, %s", gitOpName, err)
	common.Logger.Println(&common.PrefixedError{Reason: errors.New(errorMessage)})
}
