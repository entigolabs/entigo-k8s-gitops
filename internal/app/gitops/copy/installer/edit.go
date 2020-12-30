package installer

import (
	"bytes"
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
	yamlMap := yaml.Node{}
	decoder := yaml.NewDecoder(reader)
	buffer := *new(bytes.Buffer)
	encoder := yaml.NewEncoder(&buffer)
	encoder.SetIndent(2)
	for decoder.Decode(&yamlMap) == nil {
		editYaml(yamlMap, cmdData)
		if err := encoder.Encode(&yamlMap); err != nil {
			common.Logger.Fatal(&common.PrefixedError{Reason: err})
		}
	}
	if err := encoder.Close(); err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
	return &buffer
}

func editYaml(yamlMap yaml.Node, cmdData []string) {
	//keys := getKeys(cmdData[1])
	//newValue := cmdData[2]
	//fmt.Println(yamlMap) // todo implement edit logic
}
