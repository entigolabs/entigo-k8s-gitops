package installer

import (
	"errors"
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"strings"
)

func (i *Installer) isLineValid(cmdLine string) bool {
	lineSplits := strings.Split(cmdLine, " ")
	firstSplit := lineSplits[0]
	if isValidCommand(firstSplit) {
		return true
	}
	if isValidForSkip(firstSplit) {
		return false
	}
	msg := errors.New(fmt.Sprintf("skiping invalid install line - %s", cmdLine))
	common.Logger.Println(&common.Warning{Reason: msg})
	return false
}

func isValidForSkip(firstSplit string) bool {
	return firstSplit == "" || strings.Contains(firstSplit, "#")
}

func isValidCommand(firstSplit string) bool {
	return firstSplit == editCmd || firstSplit == dropCmd
}
