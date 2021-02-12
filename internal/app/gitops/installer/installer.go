package installer

import (
	"errors"
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"strings"
)

const (
	editCmd string = "edit"
	dropCmd string = "drop"
)

type Installer struct {
	Command            common.Command
	KeepRegistry       bool
	DeploymentStrategy common.DeploymentStrategy
}

func (i *Installer) Install(installInput string) {
	cmdLines := strings.Split(installInput, "\n")
	for _, cmdLine := range cmdLines {
		cmdLine = formatCmdLine(cmdLine)
		if !i.isLineValid(cmdLine) {
			continue
		}
		i.runCommand(cmdLine)
	}
}

func (i *Installer) runCommand(line string) {
	lineSplits := strings.Split(line, " ")
	cmdType := lineSplits[0]
	cmdData := lineSplits[1:]

	switch cmdType {
	case editCmd:
		i.edit(cmdData)
	case dropCmd:
		i.drop(cmdData)
	default:
		msg := fmt.Sprintf("unsupported command '%s'", cmdType)
		common.Logger.Fatal(common.PrefixedError{Reason: errors.New(msg)})
	}
	logCommandEnd(cmdType)
}

func formatCmdLine(cmdLine string) string {
	cmdLine = strings.TrimSpace(cmdLine)
	return standardizeSpaces(cmdLine)
}

func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func logCommandEnd(cmdType string) {
	cmdString := ""
	switch cmdType {
	case editCmd:
		cmdString = "edit"
	case dropCmd:
		cmdString = "drop"
	}
	common.Logger.Println(fmt.Sprintf("finised %s command", cmdString))
}
