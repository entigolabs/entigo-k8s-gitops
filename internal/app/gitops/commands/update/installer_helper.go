package update

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"io/ioutil"
	"strings"
)

type installInput struct {
	changeData []string
	fileNames  []string
}

func getInstallInput(changeData []string) string {
	input := installInput{changeData: changeData}
	yamlNames, err := readDirFiltered(".", ".yaml")
	if err != nil {
		common.Logger.Println(common.PrefixedError{Reason: err})
	}
	input.fileNames = yamlNames
	return composeInstallInput(input)
}

func composeInstallInput(input installInput) string {
	composedInput := ""
	for _, d := range input.changeData {
		yamlNamesStr := strings.Join(input.fileNames, ",")
		editLine := fmt.Sprintf("edit %s %s %s", yamlNamesStr, getUpdateInstallLocations(), d)
		composedInput += fmt.Sprintf("%s\n", editLine)
	}
	return composedInput
}

func getUpdateInstallLocations() string {
	deploymentStatefulSetDaemonSetJobLocation := "spec.template.spec.containers.*.image"
	podLocation := "spec.containers.*.image"
	cronJobLocation := "spec.jobTemplate.spec.template.spec.containers.*.image"
	return strings.Join([]string{deploymentStatefulSetDaemonSetJobLocation, podLocation, cronJobLocation}, ",")
}

func readDirFiltered(path string, suffix string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, file := range files {
		if file.IsDir() == false && strings.HasSuffix(file.Name(), suffix) {
			result = append(result, file.Name())
		}
	}
	return result, err
}
