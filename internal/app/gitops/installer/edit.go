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
	documentIndex  int
	workingKey     string
	keyExist       bool
	isImageUpdated bool
}

var editInfo = new(editInformation)

func (i *Installer) edit(input InstallInput) {
	logEditStart(input)
	for _, yamlFileName := range input.FileNames {
		editInfo.workingFile = yamlFileName
		editedBuffer := i.getEditedBuffer(yamlFileName, input)
		common.OverwriteFile(yamlFileName, editedBuffer.Bytes())
	}
}

func (i *Installer) getEditedBuffer(yamlFileName string, input InstallInput) *bytes.Buffer {
	inputYaml := common.GetFileInput(yamlFileName)
	reader := bytes.NewReader(inputYaml)
	yamlNode := yaml.Node{}
	decoder := yaml.NewDecoder(reader)
	buffer := *new(bytes.Buffer)
	encoder := yaml.NewEncoder(&buffer)
	encoder.SetIndent(2)
	for documentIndex := 1; decoder.Decode(&yamlNode) == nil; documentIndex++ {
		editInfo.documentIndex = documentIndex
		i.protectOrEdit(yamlNode, input)
		if err := encoder.Encode(&yamlNode); err != nil {
			common.Logger.Fatal(&common.PrefixedError{Reason: err})
		}
	}
	if err := encoder.Close(); err != nil {
		logEncoderClosing(yamlFileName, err)
	}
	return &buffer
}

func (i *Installer) protectOrEdit(yamlNode yaml.Node, input InstallInput) {
	if !i.isObjectProtected(&yamlNode) {
		i.editNode(yamlNode, input)
	} else {
		msg := errors.New(fmt.Sprintf("skiping update in %s in document nr %v, object is protected", editInfo.workingFile, editInfo.documentIndex))
		common.Logger.Println(&common.Warning{Reason: msg})
	}
}

func (i *Installer) editNode(yamlNode yaml.Node, input InstallInput) {
	i.editYaml(&yamlNode, input)
	i.editStrategy(&yamlNode, input)
}

func (i *Installer) isObjectProtected(yamlNode *yaml.Node) bool {
	protectionLocation := "metadata.annotations.entigo-k8s-gitops/protected"
	editInfo.workingKey = protectionLocation
	editInfo.keyExist = false
	keyValue, keyNotFoundErr := i.search(yamlNode, strings.Split(protectionLocation, "."))
	if keyNotFoundErr != nil {
		return false
	}
	isProtected, parseErr := strconv.ParseBool(keyValue)
	if parseErr != nil {
		msg := errors.New(fmt.Sprintf("unsupported key value, could not parse to boolean: %s", keyValue))
		common.Logger.Fatal(&common.PrefixedError{Reason: msg})
	}
	return isProtected
}

func (i *Installer) editYaml(yamlNode *yaml.Node, input InstallInput) {
	for _, location := range input.KeyLocations {
		editInfo.workingKey = location
		editInfo.keyExist = false
		keys := strings.Split(location, ".")
		i.replace(yamlNode, keys, input.NewValue)
	}
}

func (i *Installer) editStrategy(yamlNode *yaml.Node, input InstallInput) {
	if i.DeploymentStrategy != common.UnspecifiedStrategy && editInfo.isImageUpdated {
		newStrategyInput := InstallInput{
			Command:      input.Command,
			FileNames:    input.FileNames,
			KeyLocations: []string{"spec.strategy.type"},
			NewValue:     i.DeploymentStrategy.String(),
		}
		i.editYaml(yamlNode, newStrategyInput)
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
	editInfo.keyExist = true
	oldImage := strings.Split(oldValue, ":")[0]
	newImage := strings.Split(newValue, ":")[0]
	if isOldImageContainingNewImage(oldImage, newImage) {
		editInfo.isImageUpdated = true
		if i.KeepRegistry {
			logImageChangeWithRegistry(oldValue, newValue)
			newTag := strings.Split(newValue, ":")[1]
			return fmt.Sprintf("%s:%s", oldImage, newTag)
		}
		logImageChange(oldValue, newValue)
		return newValue
	}
	logImageCouldNotBeFound(newImage)
	return oldValue
}

func logImageChange(oldValue string, newValue string) {
	msg := errors.New(fmt.Sprintf("updating key '%s' in %s in document nr %v from '%s' to '%s'",
		editInfo.workingKey, editInfo.workingFile, editInfo.documentIndex, oldValue, newValue))
	common.Logger.Println(msg)
}

func logImageChangeWithRegistry(oldValue string, newValue string) {
	newValueWithRegistry := fmt.Sprintf("%s:%s", strings.Split(oldValue, ":")[0], newValue)
	msg := errors.New(fmt.Sprintf("updating key '%s' in %s in document nr %v from '%s' to '%s'",
		editInfo.workingKey, editInfo.workingFile, editInfo.documentIndex, oldValue, newValueWithRegistry))
	common.Logger.Println(msg)
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

func logImageCouldNotBeFound(image string) {
	msg := errors.New(fmt.Sprintf("skiping '%s' update in %s - '%s' couldn't be found", editInfo.workingKey, editInfo.workingFile, image))
	common.Logger.Println(&common.Warning{Reason: msg})
}

func logEncoderClosing(yamlFileName string, err error) {
	if strings.Contains(err.Error(), "yaml: expected STREAM-START") {
		msg := fmt.Sprintf("%s in %s", err, yamlFileName)
		common.Logger.Println(&common.Warning{Reason: errors.New(msg)})
	} else {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
}

func logEditStart(input InstallInput) {
	common.Logger.Println(fmt.Sprintf("started editing %s", strings.Join(input.FileNames, ", ")))
	common.Logger.Println(fmt.Sprintf("changing keys %s to %s", strings.Join(input.KeyLocations, ", "), input.NewValue))
}
