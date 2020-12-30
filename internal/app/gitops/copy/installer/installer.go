package installer

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/otiai10/copy"
	"gopkg.in/yaml.v2"
	"io/ioutil"
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

func edit(data []string) {
	yamlFileNames := getFileNames(data[0])
	//keys := getKeys(data[1])
	//newValue := data[2]
	for _, yamlFileName := range yamlFileNames {
		//tmpFileName := copyTempFile(yamlFile)
		editedBuffer := getEditedBuffer(yamlFileName)

		fmt.Println("xxxxxxxx>")
		fmt.Println(editedBuffer)
		fmt.Println("<xxxxxxxx")

		if err := ioutil.WriteFile(yamlFileName, editedBuffer.Bytes(), 0644); err != nil {
			common.Logger.Fatal(&common.PrefixedError{Reason: err})
		}
	}
}

func getEditedBuffer(tmpFileName string) *bytes.Buffer {
	inputYaml := common.GetFileInput(tmpFileName)
	reader := bytes.NewReader(inputYaml)
	yamlMap := yaml.MapSlice{}
	decoder := yaml.NewDecoder(reader)
	buffer := *new(bytes.Buffer)
	encoder := yaml.NewEncoder(&buffer)

	for decoder.Decode(&yamlMap) == nil {
		// todo do edit logic here
		if err := encoder.Encode(yamlMap); err != nil {
			common.Logger.Fatal(&common.PrefixedError{Reason: err})
		}
	}
	if err := encoder.Close(); err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
	return &buffer
}

func copyTempFile(yamlFile string) string {
	tmpFileName := fmt.Sprintf("%s.yaml.tmp", stripYamlSuffix(yamlFile))
	if err := copy.Copy(yamlFile, tmpFileName); err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
	return tmpFileName
}

func stripYamlSuffix(yamlFileName string) string {
	if strings.HasSuffix(yamlFileName, ".yml") {
		return strings.TrimSuffix(yamlFileName, ".yml")
	} else if strings.HasSuffix(yamlFileName, ".yaml") {
		return strings.TrimSuffix(yamlFileName, ".yaml")
	}
	common.Logger.Fatal(&common.PrefixedError{Reason: errors.New("unrecognised yaml file")})
	return ""
}

func drop(data []string) {
	filesToRemove := getFileNames(data[0])
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
