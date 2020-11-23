package update

import (
	"errors"
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/util"
	"io/ioutil"
	"os/exec"
	"strings"
)

func changeImageOsx(image string) {
	yamlFiles, _ := readDirFiltered(".", ".yaml")
	tag := util.GetTag(image)

	for _, fileName := range yamlFiles {
		substitution := fmt.Sprintf("s#image:(.*)%s:(.*)+#image: %s#g", tag, image)
		output, err := exec.Command("sed", "-Ei", "", substitution, fileName).CombinedOutput()

		if err != nil {
			if output != nil {
				util.Logger.Println(&util.Warning{Reason: errors.New(fmt.Sprintf("%s", output))})
			}
			util.Logger.Println(&util.PrefixedError{Reason: err})
		}
	}
}

func changeImageDefault(image string) {
	yamlFiles, _ := readDirFiltered(".", ".yaml")
	tag := util.GetTag(image)

	for _, fileName := range yamlFiles {
		substitution := fmt.Sprintf("s#image:(.*)/%s:(.*)+#image:\\1/%s#g", tag, image)
		output, err := exec.Command("sed", "-iE", substitution, fileName).CombinedOutput()

		if err != nil {
			if output != nil {
				util.Logger.Println(&util.Warning{Reason: errors.New(fmt.Sprintf("%s\n", output))})
			}
			util.Logger.Println(&util.PrefixedError{Reason: err})
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
