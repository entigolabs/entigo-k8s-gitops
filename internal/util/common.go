package util

import (
	"fmt"
	"os"
	"strings"
)

var RootPath = GetWd()

func GetWd() string {
	wd, _ := os.Getwd()
	return wd
}

func ChangeDir(path string) error {
	if err := os.Chdir(path); err != nil {
		return err
	}
	Logger.Println(fmt.Sprintf("changed directory to: %s", path))
	return nil
}

func GetTag(img string) string {
	tag := strings.Split(img, ":")[0]
	return tag
}

func GetAppName(appPath string) string {
	pathTokens := strings.Split(appPath, "/")
	return pathTokens[len(pathTokens)-1]
}

func GetWebUrl(repositoryUrl string) string {
	webUrl := repositoryUrl
	webUrl = strings.TrimPrefix(webUrl, "git@")
	webUrl = strings.TrimSuffix(webUrl, ".git")
	webUrl = strings.ReplaceAll(webUrl, ":", "/")
	return fmt.Sprintf("https://www.%s", webUrl)
}
