package installer

import (
	"bytes"
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"gopkg.in/yaml.v2"
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
	yamlMap := yaml.MapSlice{}
	decoder := yaml.NewDecoder(reader)
	buffer := *new(bytes.Buffer)
	encoder := yaml.NewEncoder(&buffer)
	for decoder.Decode(&yamlMap) == nil {
		editYaml(yamlMap, cmdData)
		if err := encoder.Encode(yamlMap); err != nil {
			common.Logger.Fatal(&common.PrefixedError{Reason: err})
		}
	}
	if err := encoder.Close(); err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
	return &buffer
}

func editYaml(yamlMap yaml.MapSlice, cmdData []string) {
	//keys := getKeys(cmdData[1])
	//newValue := cmdData[2]
	fmt.Println(len(yamlMap))
	if yamlMap[0].Key == "apiVersion" {
		fmt.Println(yamlMap[0].Value)
	}
}
