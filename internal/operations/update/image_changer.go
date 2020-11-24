package update

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/util"
	"github.com/rwtodd/Go.Sed/sed"
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
	tag := util.GetTag(image)
	for _, fileName := range yamlFiles {
		err := execSed(image, tag, fileName)
		if err != nil {
			break
		}
	}
}

func execSed(image string, tag string, fileName string) error {
	substitution := fmt.Sprintf("s#image:(.*)%s:(.*)+#image: %s#g", tag, image)
	sedEngine, err := sed.New(strings.NewReader(substitution))
	file, err := os.Open(fileName)
	if err == nil {
		readData, _ := ioutil.ReadAll(sedEngine.Wrap(file))
		file.Write(readData)
		err = ioutil.WriteFile(fileName, readData, 0644)

	}
	file.Close()
	return err
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
