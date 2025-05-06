package common

import (
	"regexp"
	"strings"
)

func sanitizeBranch(branch string) string {
	branch = strings.ReplaceAll(branch, "%2F", "/")
	branch = strings.ReplaceAll(branch, "_", "-")
	branch = strings.ReplaceAll(branch, " ", "-")
	branchTokens := strings.Split(branch, "/")
	branch = branchTokens[len(branchTokens)-1]
	branch = cleanIllegalCharacters(branch)
	branch = strings.ToLower(branch)
	if len(branch) > 32 {
		Logger.Println("too long branch name, trimming it to 32 characters.")
		branch = branch[0:32]
	}
	return branch
}

func cleanIllegalCharacters(str string) string {
	cleanedStr := ""
	regex := regexp.MustCompile("[^A-Za-z0-9-]")
	for i := range str {
		char := string(str[i])
		if regex.MatchString(char) {
			cleanedStr += char
		}
	}
	return cleanedStr
}
