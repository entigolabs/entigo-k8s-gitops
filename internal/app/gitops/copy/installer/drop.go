package installer

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"os"
	"strings"
)

func drop(cmdData []string) {
	logDropStart(cmdData)
	filesToRemove := strings.Split(cmdData[0], ",")
	for _, file := range filesToRemove {
		if err := os.RemoveAll(file); err != nil {
			common.Logger.Fatal(&common.PrefixedError{Reason: err})
		}
	}
}

func logDropStart(cmdData []string) {
	fileNames := formatCommaSeparatedString(cmdData[0])
	common.Logger.Println(fmt.Sprintf("started removeing %s", fileNames))
}
