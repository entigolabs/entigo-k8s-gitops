package installer

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"gopkg.in/yaml.v3"
	"io"
	"strconv"
	"strings"
)

type editInformation struct {
	workingFile         string
	objectKind          string
	documentIndex       int
	workingKey          string
	keyExist            bool
	allowStrategyChange bool
}

var editInfo = new(editInformation)

func (i *Installer) edit(input InstallInput) {
	logEditStart(input)
	for _, yamlFileName := range input.FileNames {
		editInfo.workingFile = yamlFileName
		editedBuffer := i.getEditedBuffer(yamlFileName, input)
		if editInfo.objectKind == "ConfigMap" {
			continue
		}
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
	if isFileEmpty(reader.Len(), yamlFileName) {
		return &buffer
	}
	for {
		err := decoder.Decode(&yamlNode)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			common.Logger.Fatal(&common.PrefixedError{Reason: err})
		}
		editInfo.documentIndex = editInfo.documentIndex + 1
		i.skipOrEdit(yamlNode, input)
		if err := encoder.Encode(&yamlNode); err != nil {
			common.Logger.Fatal(&common.PrefixedError{Reason: err})
		}
	}
	if err := encoder.Close(); err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
	return &buffer
}

func (i *Installer) skipOrEdit(yamlNode yaml.Node, input InstallInput) {
	if kind, _ := i.getObjectKind(&yamlNode); kind == "ConfigMap" {
		msg := fmt.Errorf("skipping update in %s in document nr %v, object is ConfigMap", editInfo.workingFile, editInfo.documentIndex)
		common.Logger.Println(&common.Warning{Reason: msg})
		return
	} else if i.isObjectProtected(&yamlNode) {
		msg := fmt.Errorf("skipping update in %s in document nr %v, object is protected", editInfo.workingFile, editInfo.documentIndex)
		common.Logger.Println(&common.Warning{Reason: msg})
		return
	}
	i.editDocumentNode(yamlNode, input)
}

func (i *Installer) editDocumentNode(yamlNode yaml.Node, input InstallInput) {
	i.editYaml(&yamlNode, input)
	i.editStrategy(&yamlNode, input.FileNames)
}

func (i *Installer) isObjectProtected(yamlNode *yaml.Node) bool {
	keyValue, keyNotFoundErr := i.getValue(yamlNode, "metadata.annotations.entigo-k8s-gitops/protected")
	if keyNotFoundErr != nil {
		return false
	}
	isProtected, parseErr := strconv.ParseBool(keyValue)
	if parseErr != nil {
		msg := fmt.Errorf("unsupported key value, could not parse to boolean: %s", keyValue)
		common.Logger.Fatal(&common.PrefixedError{Reason: msg})
	}
	return isProtected
}

func (i *Installer) getValue(yamlNode *yaml.Node, keyLocation string) (string, error) {
	editInfo.workingKey = keyLocation
	editInfo.keyExist = false
	keyValue, keyNotFoundErr := i.search(yamlNode, strings.Split(keyLocation, "."))
	if keyNotFoundErr != nil {
		return "", keyNotFoundErr
	}
	return keyValue, nil
}

func (i *Installer) getObjectKind(yamlNode *yaml.Node) (string, error) {
	kind, keyNotFoundErr := i.getValue(yamlNode, "kind")
	if keyNotFoundErr != nil {
		return "", keyNotFoundErr
	}
	editInfo.objectKind = kind
	return kind, nil
}

func (i *Installer) editYaml(yamlNode *yaml.Node, input InstallInput) {
	for _, location := range input.KeyLocations {
		editInfo.workingKey = location
		editInfo.keyExist = false
		keys := strings.Split(location, ".")
		i.replace(yamlNode, keys, input.NewValue)
	}
}

func (i *Installer) editStrategy(yamlNode *yaml.Node, fileNames []string) {
	editInfo.allowStrategyChange = true
	if i.DeploymentStrategy != common.UnspecifiedStrategy {
		newStrategyInput := InstallInput{
			Command:      EditCmd,
			FileNames:    fileNames,
			KeyLocations: []string{"spec.strategy.type"},
			NewValue:     i.DeploymentStrategy.String(),
		}
		i.editYaml(yamlNode, newStrategyInput)
	}
	editInfo.allowStrategyChange = false
}

func (i *Installer) replace(node *yaml.Node, keys []string, newValue string) {
	identifier := keys[0]
	if node.Kind == yaml.DocumentNode {
		i.replace(node.Content[0], keys, newValue)
	}
	if seqPos, err := strconv.Atoi(identifier); err == nil {
		if len(node.Content)-1 < seqPos {
			msg := fmt.Errorf("skiping '%s' copy in %s - key doesn't exist", editInfo.workingKey, editInfo.workingFile)
			common.Logger.Println(&common.Warning{Reason: msg})
			return
		}
		seqPosNode := node.Content[seqPos]
		if seqPosNode.Kind == yaml.ScalarNode {
			seqPosNode.Value = i.getNewValue(seqPosNode.Value, newValue)
		} else {
			i.replace(node.Content[seqPos], keys[1:], newValue)
		}
	}
	if identifier == "*" {
		for j := range node.Content {
			i.replace(node.Content[j], keys[1:], newValue)
		}
	}
	for j, n := range node.Content {
		if n.Value == identifier {
			if len(keys) <= 1 && node.Content[j+1].Content == nil {
				node.Content[j+1].Value = i.getNewValue(node.Content[j+1].Value, newValue)
			} else {
				i.replace(node.Content[j+1], keys[1:], newValue)
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
	if editInfo.allowStrategyChange {
		return getStrategyChangeSpecificNewValue(oldValue, newValue)
	}
	return i.getImageChangeSpecificNewValue(oldValue, newValue)
}

func (i *Installer) appendHistoryEntry(oldValue string, newValue string) {
	historyEntry := InstallHistory{
		NewValue:      newValue,
		OldValue:      oldValue,
		workingFile:   editInfo.workingFile,
		documentIndex: editInfo.documentIndex,
		workingKey:    editInfo.workingKey,
	}
	i.InstallHistory = append(i.InstallHistory, historyEntry)
}

func (i *Installer) getImageChangeSpecificNewValue(oldValue string, newValue string) string {
	editInfo.keyExist = true
	oldValueSplits := strings.Split(oldValue, ":")
	oldImage := oldValueSplits[0]
	newValueSplits := strings.Split(newValue, ":")
	newImage := newValueSplits[0]
	if isOldImageContainingNewImage(oldImage, newImage) {
		oldValueImgAndTag := fmt.Sprintf("%s:%s", newImage, oldValueSplits[1])
		i.appendHistoryEntry(oldValueImgAndTag, newValue)
		if i.KeepRegistry {
			logImageChangeWithRegistry(oldValue, newValue)
			newTag := newValueSplits[1]
			return fmt.Sprintf("%s:%s", oldImage, newTag)
		}
		logImageChange(oldValue, newValue)
		return newValue
	}
	logImageCouldNotBeFound(newImage)
	return oldValue
}

func logImageChange(oldValue string, newValue string) {
	msg := fmt.Errorf("updating key '%s' in %s in document nr %v from '%s' to '%s'",
		editInfo.workingKey, editInfo.workingFile, editInfo.documentIndex, oldValue, newValue)
	common.Logger.Println(msg)
}

func logImageChangeWithRegistry(oldValue string, newValue string) {
	newValueWithRegistry := fmt.Sprintf("%s:%s", strings.Split(oldValue, ":")[0], newValue)
	msg := fmt.Errorf("updating key '%s' in %s in document nr %v from '%s' to '%s'",
		editInfo.workingKey, editInfo.workingFile, editInfo.documentIndex, oldValue, newValueWithRegistry)
	common.Logger.Println(msg)
}

func isOldImageContainingNewImage(oldImage string, newImage string) bool {
	return strings.Contains(oldImage, newImage) && areImageEndingsMatching(oldImage, newImage)
}

func areImageEndingsMatching(oldImage string, newImage string) bool {
	newImgIndex := strings.Index(oldImage, newImage)
	return len(oldImage) == newImgIndex+len(newImage)
}

func getStrategyChangeSpecificNewValue(oldValue string, newValue string) string {
	logImageChange(oldValue, newValue)
	return newValue
}

func isFileEmpty(readerLen int, yamlFileName string) bool {
	if readerLen == 0 {
		msg := fmt.Sprintf("%s is empty, nothing is changed", yamlFileName)
		common.Logger.Println(&common.Warning{Reason: errors.New(msg)})
		return true
	}
	return false
}

func logImageCouldNotBeFound(image string) {
	msg := fmt.Errorf("skiping '%s' update in %s - '%s' couldn't be found", editInfo.workingKey, editInfo.workingFile, image)
	common.Logger.Println(&common.Warning{Reason: msg})
}

func logEditStart(input InstallInput) {
	common.Logger.Printf("started editing %s\n", strings.Join(input.FileNames, ", "))
	common.Logger.Printf("changing keys %s to %s\n", strings.Join(input.KeyLocations, ", "), input.NewValue)
}
