package installer

import (
	"errors"
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"os"
	"strings"
)

const installFile = "install.txt"

const (
	editCmd string = "edit"
	dropCmd string = "drop"
)

type Installer struct {
	GitBranch string // featureBranch
	AppName   string // argoapp
}

func (i *Installer) Install() {
	installInput := common.GetFileInput(installFile)
	cmdLines := strings.Split(string(installInput), "\n")
	for _, cmdLine := range cmdLines {
		if cmdLine == "" {
			return
		}
		cmdLine = i.specifyLineVars(cmdLine)
		runCommand(cmdLine)
	}
}

func (i *Installer) specifyLineVars(line string) string {
	// TODO check replacers are correct
	line = strings.ReplaceAll(line, saltedVariable("featureBranch"), i.GitBranch)
	line = strings.ReplaceAll(line, saltedVariable("workname"), fmt.Sprintf("%s-%s", i.AppName, i.GitBranch))
	line = strings.ReplaceAll(line, saltedVariable("url"), i.getFeatureUrl())
	return line
}

func (i *Installer) getFeatureUrl() string {
	if i.GitBranch == "master" {
		return i.AppName
	}
	return fmt.Sprintf("%s-%s.fleetcomplete.dev", i.AppName, i.GitBranch)
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
		err := fmt.Sprintf("unsupported command '%s'", cmdType)
		common.Logger.Fatal(common.PrefixedError{Reason: errors.New(err)})
	}
}

func drop(cmdData []string) {
	filesToRemove := getFileNames(cmdData[0])
	for _, file := range filesToRemove {
		if err := os.RemoveAll(file); err != nil {
			common.Logger.Fatal(&common.PrefixedError{Reason: err})
		}
	}
}

func getKeys(keysInString string) []string {
	return strings.Split(keysInString, ",")
}

func getFileNames(fileNamesInString string) []string {
	return strings.Split(fileNamesInString, ",")
}
