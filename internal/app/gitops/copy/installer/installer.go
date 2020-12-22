package installer

import (
	"errors"
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"io/ioutil"
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
	input := getFileInput(installFile)

	lines := strings.Split(string(input), "\n")
	for _, line := range lines {
		line = i.specifyLineVars(line)
		runCommand(line)
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
	// TODO impl getFeatureUrl
	return "getFeatureUrl(argoapp, featureBranch)"
}

func saltedVariable(variable string) string {
	return fmt.Sprintf("{{%s}}", variable)
}

func runCommand(line string) {
	lineSplits := strings.Split(line, " ")
	installCmd := lineSplits[0]
	data := lineSplits[1:]

	switch installCmd {
	case editCmd:
		edit(data)
	case dropCmd:
		drop(data)
	default:
		err := fmt.Sprintf("unsupported command '%s'", installCmd)
		common.Logger.Fatal(common.PrefixedError{Reason: errors.New(err)})
	}
}

func edit(data []string) {
	// sh('edityaml.py ' + cmd[1])
	fmt.Println(data)
}

func drop(data []string) {
	// sh('rm -f ' + cmd[1])
	fmt.Println(data)
}

func getFileInput(fileName string) []byte {
	input, err := ioutil.ReadFile(fileName)
	if err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
	return input
}
