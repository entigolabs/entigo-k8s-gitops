package update

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/entigolabs/entigo-k8s-gitops/internal/util"
	"io/ioutil"
	"strings"
)

type updaterFlags struct {
	images       string
	keepRegistry bool
}

func changeImages(flags updaterFlags) {
	imgNameTokens := strings.Split(flags.images, ",")
	for _, img := range imgNameTokens {
		changeSpecificImage(img, flags.keepRegistry)
	}
}

func changeSpecificImage(image string, keepRegistry bool) {
	yamlFiles, _ := readDirFiltered(".", ".yaml")
	for _, fileName := range yamlFiles {
		imageChanger(image, fileName, keepRegistry)
	}
}

func imageChanger(image string, fileName string, keepRegistry bool) {
	input := getFileInput(fileName)
	output := getChangedOutput(input, image, keepRegistry)
	overwriteFile(fileName, output)
}

func getFileInput(fileName string) []byte {
	input, err := ioutil.ReadFile(fileName)
	if err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
	return input
}

func getChangedOutput(input []byte, image string, keepRegistry bool) string {
	lines := strings.Split(string(input), "\n")
	changeImageLines(lines, image, keepRegistry)
	return strings.Join(lines, "\n")
}

func changeImageLines(lines []string, image string, keepRegistry bool) {
	for lineIndex, line := range lines {
		if isImageLine(line, getImageName(image)) {
			if keepRegistry {
				changeKeepingRegistry(lines, image, line, lineIndex)
			} else {
				changeCompletely(lines, image, line, lineIndex)
			}
		}
	}
}

func overwriteFile(fileName string, output string) {
	if err := ioutil.WriteFile(fileName, []byte(output), 0644); err != nil {
		common.Logger.Fatal(&util.PrefixedError{Reason: err})
	}
}

func changeCompletely(lines []string, image string, line string, lineIndex int) {
	imageFieldTitle := strings.Split(line, ":")[0]
	lines[lineIndex] = fmt.Sprintf("%s: %s", imageFieldTitle, image)
}

func changeKeepingRegistry(lines []string, image string, line string, lineIndex int) {
	imageLineSplits := strings.Split(line, ":")
	imageFieldTitle := imageLineSplits[0]
	registryName := imageLineSplits[1]
	lines[lineIndex] = fmt.Sprintf("%s: %s:%s", imageFieldTitle, strings.TrimSpace(registryName), getImageTag(image))
}

func isImageLine(line string, imageName string) bool {
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "image:") && strings.Contains(line, imageName) {
		return true
	}
	return false
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

func getImageName(image string) string {
	return strings.Split(image, ":")[0]
}

func getImageTag(image string) string {
	return strings.Split(image, ":")[1]
}
