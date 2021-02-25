package update

import (
	"errors"
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/installer"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func composeInstallInput(images []string, recursiveFileSearching bool) []installer.InstallInput {
	yamlNames := getFileNames(recursiveFileSearching)
	updateLocations := getUpdateInstallLocations()
	var installInputs []installer.InstallInput
	for _, image := range images {
		installInput := installer.InstallInput{
			Command:      installer.EditCmd,
			FileNames:    yamlNames,
			KeyLocations: updateLocations,
			NewValue:     image,
		}
		installInputs = append(installInputs, installInput)
	}
	return installInputs
}

func getFileNames(recursively bool) []string {
	var fileNames []string
	if recursively {
		fileNames = getFilesRecursively()
	} else {
		fileNames = getFiles()
	}
	validateFileNames(fileNames)
	return fileNames
}

func validateFileNames(fileNames []string) {
	if fileNames == nil {
		msg := ".yaml files could not be found"
		common.Logger.Fatal(&common.PrefixedError{Reason: errors.New(msg)})
	}
}

func getFilesRecursively() []string {
	var yamlNames []string
	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.HasSuffix(info.Name(), ".yaml") {
				yamlNames = append(yamlNames, path)
			}
			return nil
		})
	if err != nil {
		common.Logger.Println(common.PrefixedError{Reason: err})
	}
	return yamlNames
}

func getFiles() []string {
	yamlNames, err := readDirFiltered(".", ".yaml")
	if err != nil {
		common.Logger.Println(common.PrefixedError{Reason: err})
	}
	return yamlNames
}

func getUpdateInstallLocations() []string {
	locations := append(getTypeSpecificContainerLocations(containers), getTypeSpecificContainerLocations(initContainers)...)
	return locations
}

func getTypeSpecificContainerLocations(ct containersType) []string {
	cTypeAsString := ct.string()
	deploymentStatefulSetDaemonSetJobLocationTemplate := fmt.Sprintf("spec.template.spec.%s.*.image", cTypeAsString)
	podLocationTemplate := fmt.Sprintf("spec.%s.*.image", cTypeAsString)
	cronJobLocationTemplate := fmt.Sprintf("spec.jobTemplate.spec.template.spec.%s.*.image", cTypeAsString)
	return []string{deploymentStatefulSetDaemonSetJobLocationTemplate, podLocationTemplate, cronJobLocationTemplate}
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
