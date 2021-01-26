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
	for i := range str {
		char := string(str[i])
		cleanMatch, err := regexp.MatchString("[A-Za-z0-9-]", char)
		if err != nil {
			Logger.Fatal(&PrefixedError{Reason: err})
		}
		if cleanMatch {
			cleanedStr += char
		}
	}
	return cleanedStr
}
