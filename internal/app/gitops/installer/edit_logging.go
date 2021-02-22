package installer

import (
	"errors"
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"strings"
)

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
