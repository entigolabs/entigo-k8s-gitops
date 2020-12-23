package common

import (
	"regexp"
	"strings"
)

func sanitizeBranch(branch string) string {
	branch = strings.ReplaceAll(branch, "%2F", "/")
	branch = strings.ReplaceAll(branch, "_", "-")
	branchTokens := strings.Split(branch, "/")
	branch = branchTokens[len(branchTokens)-1]
	branch = cleanIllegalCharacters(branch)
	if len(branch) > 32 {
		Logger.Println("too long branch name, trimming it to 32 characters.")
		branch = branch[0:32]
	}
	return branch
}

func cleanIllegalCharacters(branch string) string {
	cleanedBranch := ""
	for index := range branch {
		character := string([]rune(branch)[index])
		cleanMatch, err := regexp.MatchString("^[A-Za-z0-9-]+$", character)
		if err != nil {
			Logger.Fatal(&PrefixedError{Reason: err})
		}
		if cleanMatch {
			cleanedBranch += character
		}
	}
	return cleanedBranch
}
