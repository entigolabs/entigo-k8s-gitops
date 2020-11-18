package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"
)

var logger = log.New(os.Stderr, "gitops: ", log.LstdFlags|log.Lshortfile)
var rootPath = getwd()
var deploymentDir = "deployment"

func changeDir(path string) error {
	if err := os.Chdir(path); err != nil {
		return err
	}
	fmt.Printf("Changed directory to: %s\n", path)
	return nil
}

func getwd() string {
	wd, _ := os.Getwd()
	return wd
}

func genMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func getTag(img string) string {
	tag := strings.Split(img, ":")[0]
	return tag
}

func getWebUrl(repositoryUrl string) string {
	webUrl := repositoryUrl
	webUrl = strings.TrimPrefix(webUrl, "git@")
	webUrl = strings.TrimSuffix(webUrl, ".git")
	webUrl = strings.ReplaceAll(webUrl, ":", "/")
	return fmt.Sprintf("https://www.%s", webUrl)
}
