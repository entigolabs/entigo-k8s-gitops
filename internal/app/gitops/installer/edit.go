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

type editInformation struct {
	workingFile    string
	workingKey     string
	keyDontExist   bool
	isImageUpdated bool
}

var editInfo = new(editInformation)

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
		i.editStrategy(&yamlNode, cmdData[0])
		if err := encoder.Encode(&yamlNode); err != nil {
			common.Logger.Fatal(&common.PrefixedError{Reason: err})
		}
	}
	if err := encoder.Close(); err != nil {
		logEncoderClosing(yamlFileName, err)
	}
	return &buffer
}

func (i *Installer) editYaml(yamlNode *yaml.Node, cmdData []string) {
	replaceLocations := strings.Split(cmdData[1], ",")
	newValue := cmdData[2]
	for _, replaceLocation := range replaceLocations {
		editInfo.workingKey = replaceLocation
		editInfo.keyDontExist = false
		keys := strings.Split(replaceLocation, ".")
		i.replace(yamlNode, keys, newValue)
	}
}

func (i *Installer) editStrategy(yamlNode *yaml.Node, fileNames string) {
	if i.DeploymentStrategy != common.UnspecifiedStrategy && editInfo.isImageUpdated {
		strategyLocation := "spec.strategy.type"
		newStrategy := common.ConvDeploymentStrategyToStr(i.DeploymentStrategy)
		editData := []string{fileNames, strategyLocation, newStrategy}
		i.editYaml(yamlNode, editData)
	}
	editInfo.isImageUpdated = false
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
		for i := range node.Content {
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
			inst.logIfKeyDontExist(node, keys, i)
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
	if editInfo.isImageUpdated {
		return getStrategyChangeSpecificNewValue(newValue)
	}
	return i.getImageChangeSpecificNewValue(oldValue, newValue)
}

func (i *Installer) getImageChangeSpecificNewValue(oldValue string, newValue string) string {
	editInfo.keyDontExist = true
	oldImage := strings.Split(oldValue, ":")[0]
	newImage := strings.Split(newValue, ":")[0]
	if isOldImageContainingNewImage(oldImage, newImage) {
		editInfo.isImageUpdated = true
		if i.KeepRegistry {
			newTag := strings.Split(newValue, ":")[1]
			return fmt.Sprintf("%s:%s", oldImage, newTag)
		}
		return newValue
	}
	logImageCouldNotBeFound(newImage)
	return oldValue
}

func isOldImageContainingNewImage(oldImage string, newImage string) bool {
	return strings.Contains(oldImage, newImage) && areImageEndingsMatching(oldImage, newImage)
}

func areImageEndingsMatching(oldImage string, newImage string) bool {
	newImgIndex := strings.Index(oldImage, newImage)
	if len(oldImage) != newImgIndex+len(newImage) {
		return false
	}
	return true
}

func getStrategyChangeSpecificNewValue(newValue string) string {
	return newValue
}
