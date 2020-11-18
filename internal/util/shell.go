package util

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

func changeImageOsx(image string) {
	yamlFiles, _ := readDirFiltered(".", ".yaml")
	tag := getTag(image)

	for _, fileName := range yamlFiles {
		substitution := fmt.Sprintf("s#image:(.*)%s:(.*)+#image: %s#g", tag, image)
		output, err := exec.Command("sed", "-Ei", "", substitution, fileName).CombinedOutput()

		if err != nil {
			if output != nil {
				logger.Println(&warning{errors.New(fmt.Sprintf("%s\n", output))})
			}
			logger.Println(&prefixedError{err})
		}
	}
}

func changeImageDefault(image string) {
	yamlFiles, _ := readDirFiltered(".", ".yaml")
	tag := getTag(image)

	for _, fileName := range yamlFiles {
		substitution := fmt.Sprintf("s#image:(.*)/%s:(.*)+#image:\\1/%s#g", tag, image)
		output, err := exec.Command("sed", "-iE", substitution, fileName).CombinedOutput()

		if err != nil {
			if output != nil {
				logger.Println(&warning{errors.New(fmt.Sprintf("%s\n", output))})
			}
			logger.Println(&prefixedError{err})
		}
	}
}

func readDirFiltered(path string, suffix string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, f := range files {
		if f.IsDir() == false && strings.HasSuffix(f.Name(), suffix) {
			result = append(result, f.Name())
		}
	}
	return result, err
}
