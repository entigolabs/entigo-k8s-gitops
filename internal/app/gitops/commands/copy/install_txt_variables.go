package copy

import (
	"fmt"
	"strings"
)

type installTxtVariables struct {
	appBranch  string
	appName    string
	featureUrl string
}

func initInstallTxtVariables(appBranch string, appName string, appDomain string) *installTxtVariables {
	featureUrl := ""
	if appDomain == "" {
	    appDomain = "localhost"
	}
	if appBranch == "master" {
		featureUrl = appName
	} else {
		featureUrl = fmt.Sprintf("%s-%s.%s", appName, appBranch, appDomain)
	}
	return &installTxtVariables{appBranch: appBranch, appName: appName, featureUrl: featureUrl}
}

func (i *installTxtVariables) specifyInstallVariables(installInput string) string {
	unspecifiedLines := strings.Split(installInput, "\n")
	specifiedLines := ""
	for _, unspecifiedLine := range unspecifiedLines {
		specifiedLines += fmt.Sprintf("%s\n", i.specifyLineVariables(unspecifiedLine))
	}
	fmt.Println("specifyInstallVariables: %s",specifiedLines)
	return specifiedLines
}

func (i *installTxtVariables) specifyLineVariables(line string) string {
	line = strings.ReplaceAll(line, i.saltedVariable("featureBranch"), i.appBranch)
	line = strings.ReplaceAll(line, i.saltedVariable("workname"), fmt.Sprintf("%s-%s", i.appName, i.appBranch))
	line = strings.ReplaceAll(line, i.saltedVariable("url"), i.featureUrl)
	return line
}

func (i *installTxtVariables) saltedVariable(variable string) string {
	return fmt.Sprintf("{{%s}}", variable)
}
