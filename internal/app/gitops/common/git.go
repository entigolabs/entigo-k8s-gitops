package common

import (
	"fmt"
	"strings"
)

func GetRepositoryRootPath(repoSshUrl string) string {
	return fmt.Sprintf("%s/%s/%s", RootPath, GitOpsWd, GetRepositoryName(repoSshUrl))
}

func GetRepositoryName(repoSshUrl string) string {
	tokens := strings.Split(repoSshUrl, "/")
	lastToken := tokens[len(tokens)-1]
	return lastToken[:len(lastToken)-4]
}
