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
	AppBranch string // featureBranch  // todo rm comment
	AppName   string // argoapp // todo rm comment
}

func (i *Installer) Install(installInput string) {
	cmdLines := strings.Split(installInput, "\n")
	for _, cmdLine := range cmdLines {
		if cmdLine == "" {
			return
		}
		cmdLine = i.specifyLineVars(cmdLine)
		runCommand(cmdLine)
	}
}

func (i *Installer) specifyLineVars(line string) string {
	line = strings.ReplaceAll(line, saltedVariable("featureBranch"), i.AppBranch)
	line = strings.ReplaceAll(line, saltedVariable("workname"), fmt.Sprintf("%s-%s", i.AppName, i.AppBranch))
	line = strings.ReplaceAll(line, saltedVariable("url"), i.getFeatureUrl())
	return line
}

func (i *Installer) getFeatureUrl() string {
	if i.AppBranch == "master" {
		return i.AppName
	}
	return fmt.Sprintf("%s-%s.fleetcomplete.dev", i.AppName, i.AppBranch)
}

func saltedVariable(variable string) string {
	return fmt.Sprintf("{{%s}}", variable)
}

func runCommand(line string) {
	lineSplits := strings.Split(line, " ")
	cmdType := lineSplits[0]
	cmdData := lineSplits[1:]

	switch cmdType {
	case editCmd:
		edit(cmdData)
	case dropCmd:
		drop(cmdData)
	default:
		msg := fmt.Sprintf("unsupported command '%s'", cmdType)
		common.Logger.Fatal(common.PrefixedError{Reason: errors.New(msg)})
	}
	logCommandEnd(cmdType)
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
