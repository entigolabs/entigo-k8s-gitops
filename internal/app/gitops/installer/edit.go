package installer

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"gopkg.in/yaml.v3"
	"strconv"
	"strings"
)

type logInfo struct {
	workingFile string
	workingKey  string
}

var editInfo = new(logInfo)

func (i *Installer) edit(cmdData []string) {
	logEditStart(cmdData)
	yamlFileNames := strings.Split(cmdData[0], ",")
	for _, yamlFileName := range yamlFileNames {
		editInfo.workingFile = yamlFileName
		editedBuffer := i.getEditedBuffer(yamlFileName, cmdData)
		common.OverwriteFile(yamlFileName, editedBuffer.Bytes())
	}
}

func (i *Installer) getEditedBuffer(yamlFileName string, cmdData []string) *bytes.Buffer {
	inputYaml := common.GetFileInput(yamlFileName)
	reader := bytes.NewReader(inputYaml)
	yamlNode := yaml.Node{}
	decoder := yaml.NewDecoder(reader)
	buffer := *new(bytes.Buffer)
	encoder := yaml.NewEncoder(&buffer)
	encoder.SetIndent(2)
	for decoder.Decode(&yamlNode) == nil {
		i.editYaml(&yamlNode, cmdData)
		if err := encoder.Encode(&yamlNode); err != nil {
			common.Logger.Fatal(&common.PrefixedError{Reason: err})
		}
	}
	if err := encoder.Close(); err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
	return &buffer
}

func (i *Installer) editYaml(yamlNode *yaml.Node, cmdData []string) {
	replaceLocations := strings.Split(cmdData[1], ",")
	newValue := cmdData[2]
	for _, replaceLocation := range replaceLocations {
		editInfo.workingKey = replaceLocation
		keys := strings.Split(replaceLocation, ".")
		i.replace(yamlNode, keys, newValue)
	}
}

// TODO Should it be refactored?
func (inst *Installer) replace(node *yaml.Node, keys []string, newValue string) {
	identifier := keys[0]
	if node.Kind == yaml.DocumentNode {
		inst.replace(node.Content[0], keys, newValue)
	}
	if seqPos, err := strconv.Atoi(identifier); err == nil {
		if len(node.Content)-1 < seqPos {
			msg := errors.New(fmt.Sprintf("skiping '%s' copy in %s - key doesn't exist", editInfo.workingKey, editInfo.workingFile))
			common.Logger.Println(&common.Warning{Reason: msg})
			return
		}
		seqPosNode := node.Content[seqPos]
		if seqPosNode.Kind == yaml.ScalarNode {
			seqPosNode.Value = inst.getNewValue(seqPosNode.Value, newValue)
		} else {
			inst.replace(node.Content[seqPos], keys[1:], newValue)
		}
	}
	if identifier == "*" {
		for i, _ := range node.Content {
			inst.replace(node.Content[i], keys[1:], newValue)
		}
	}
	for i, n := range node.Content {
		if n.Value == identifier {
			if len(keys) <= 1 && node.Content[i+1].Content == nil {
				node.Content[i+1].Value = inst.getNewValue(node.Content[i+1].Value, newValue)
			} else {
				inst.replace(node.Content[i+1], keys[1:], newValue)
			}
		} else {
			continue
		}
	}
}

func (i *Installer) getNewValue(oldValue string, newValue string) string {
	switch i.Command {
	case common.UpdateCmd:
		return i.getUpdateSpecificNewValue(oldValue, newValue)
	case common.CopyCmd:
		return newValue
	default:
		common.Logger.Fatal(&common.PrefixedError{Reason: errors.New("unsupported command")})
	}
	return newValue
}

func (i *Installer) getUpdateSpecificNewValue(oldValue string, newValue string) string {
	image := strings.Split(newValue, ":")[0]
	if strings.Contains(oldValue, image) {
		if i.KeepRegistry {
			registry := strings.Split(oldValue, ":")[0]
			newTag := strings.Split(newValue, ":")[1]
			return fmt.Sprintf("%s:%s", registry, newTag)
		}
		return newValue
	}
	return oldValue
}

func logEditStart(cmdData []string) {
	fileNames := formatCommaSeparatedString(cmdData[0])
	formattedKeys := formatCommaSeparatedString(cmdData[1])
	newValue := cmdData[2]
	common.Logger.Println(fmt.Sprintf("started editing %s", fileNames))
	common.Logger.Println(fmt.Sprintf("changing keys %s to %s", formattedKeys, newValue))
}