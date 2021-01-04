package installer

import (
	"bytes"
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"gopkg.in/yaml.v3"
)

func edit(cmdData []string) {
	yamlFileNames := getFileNames(cmdData[0])
	for _, yamlFileName := range yamlFileNames {
		editedBuffer := getEditedBuffer(yamlFileName, cmdData)
		common.OverwriteFile(yamlFileName, editedBuffer.Bytes())
	}
}

func getEditedBuffer(yamlFileName string, cmdData []string) *bytes.Buffer {
	inputYaml := common.GetFileInput(yamlFileName)
	reader := bytes.NewReader(inputYaml)
	yamlNode := yaml.Node{}
	decoder := yaml.NewDecoder(reader)
	buffer := *new(bytes.Buffer)
	encoder := yaml.NewEncoder(&buffer)
	encoder.SetIndent(2)
	for decoder.Decode(&yamlNode) == nil {
		editYaml(yamlNode, cmdData)
		if err := encoder.Encode(&yamlNode); err != nil {
			common.Logger.Fatal(&common.PrefixedError{Reason: err})
		}
	}
	if err := encoder.Close(); err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
	return &buffer
}

func editYaml(yamlNode yaml.Node, cmdData []string) {
	//keys := getKeys(cmdData[1])
	//newValue := cmdData[2]
	fmt.Println(yamlNode) // todo implement edit logic
}
