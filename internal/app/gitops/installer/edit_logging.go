package installer

import (
	"errors"
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"gopkg.in/yaml.v3"
	"strings"
)

func (i *Installer) logIfKeyDontExist(node *yaml.Node, keys []string, index int) {
	if len(keys) > 1 && index == len(node.Content)-1 && !editInfo.keyDontExist && i.Command == common.UpdateCmd {
		msg := errors.New(fmt.Sprintf("skiping '%s' update in %s - key doesn't exist", editInfo.workingKey, editInfo.workingFile))
		common.Logger.Println(&common.Warning{Reason: msg})
		editInfo.keyDontExist = true
	}
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

func logEditStart(cmdData []string) {
	fileNames := formatCommaSeparatedString(cmdData[0])
	formattedKeys := formatCommaSeparatedString(cmdData[1])
	newValue := cmdData[2]
	common.Logger.Println(fmt.Sprintf("started editing %s", fileNames))
	common.Logger.Println(fmt.Sprintf("changing keys %s to %s", formattedKeys, newValue))
}
