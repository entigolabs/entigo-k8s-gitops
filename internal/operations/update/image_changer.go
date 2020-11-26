package update

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/util"
	"io/ioutil"
	"os"
	"strings"
)

func changeImages() {
	imgNameTokens := strings.Split(flgs.images, ",")
	for _, img := range imgNameTokens {
		changeSpecificImage(img)
	}
}

func changeSpecificImage(image string) {
	yamlFiles, _ := readDirFiltered(".", ".yaml")
	for _, fileName := range yamlFiles {
		imageChanger(image, fileName)
	}
}

func imageChanger(image string, fileName string) {
	input := getFileInput(fileName)
	output := getChangedOutput(input, image)
	overwriteFile(fileName, output)
}

func getFileInput(fileName string) []byte {
	input, err := ioutil.ReadFile(fileName)
	if err != nil {
		util.Logger.Println(&util.PrefixedError{Reason: err})
		os.Exit(1)
	}
	return input
}

func getChangedOutput(input []byte, image string) string {
	lines := strings.Split(string(input), "\n")
	changeImageLines(lines, image)
	return strings.Join(lines, "\n")
}

func changeImageLines(lines []string, image string) {
	for lineIndex, line := range lines {
		if isImageLine(line, getImageName(image)) {
			if flgs.keepRegistry {
				changeKeepingRegistry(lines, image, line, lineIndex)
			} else {
				changeCompletely(lines, image, line, lineIndex)
			}
		}
	}
}

func overwriteFile(fileName string, output string) {
	if err := ioutil.WriteFile(fileName, []byte(output), 0644); err != nil {
		util.Logger.Println(&util.PrefixedError{Reason: err})
		os.Exit(1)
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
