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
		logAlreadyUpToDate(pullOp, err)
	} else if isConflictErr(err) {
		common.Logger.Printf("couldn't git %s, %s", pullOp, err)
		return err
	} else {
		logDefaultGitOpErr(pullOp, err)
		os.Exit(1)
	}
	return nil
}

func handlePushErr(err error) error {
	pushOp := "push"
	if err == git.NoErrAlreadyUpToDate {
		logAlreadyUpToDate(pushOp, err)
	} else if isConflictErr(err) {
		common.Logger.Printf("couldn't git %s, %s\n", pushOp, err)
		return err
	} else {
		logDefaultGitOpErr(pushOp, err)
		os.Exit(1)
	}
	return nil
}

func logAlreadyUpToDate(gitOpName string, err error) {
	common.Logger.Printf("skipping git %s, %s\n", gitOpName, err)
}

func isConflictErr(err error) bool {
	isErrNonFastForwardUpdate := strings.Contains(err.Error(), git.ErrNonFastForwardUpdate.Error()) ||
		err == git.ErrNonFastForwardUpdate
	isErrUpdateReference := strings.Contains(err.Error(), errUpdateReference.Error())
	return isErrNonFastForwardUpdate || isErrUpdateReference
}

func logDefaultGitOpErr(gitOpName string, err error) {
	errorMessage := fmt.Sprintf("git %s failed, %s", gitOpName, err)
	common.Logger.Println(&common.PrefixedError{Reason: errors.New(errorMessage)})
}
