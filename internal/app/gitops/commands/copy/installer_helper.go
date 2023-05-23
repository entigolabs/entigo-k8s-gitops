package copy

import (
	"errors"
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/installer"
	"strings"
)

func composeInstallInput(installTxt string) []installer.InstallInput {
	var installInputs []installer.InstallInput
	for _, line := range strings.Split(installTxt, "\n") {
		line = formatCmdLine(line)
		if !isLineValid(line) {
			continue
		}
		installInput, err := convInstallLineToInstallInput(line)
		if err != nil {
			common.Logger.Fatal(&common.PrefixedError{Reason: err})
		}
		installInputs = append(installInputs, installInput)
	}
	return installInputs
}

func convInstallLineToInstallInput(installLine string) (installer.InstallInput, error) {
	lineSplits := strings.Split(installLine, " ")
	installCommand := installer.ConvStrToInstallCommand(lineSplits[0])
	switch installCommand {
	case installer.EditCmd:
		return convInstallLineToEditInput(lineSplits), nil
	case installer.DropCmd:
		return convInstallLineToDropInput(lineSplits), nil
	}
	return installer.InstallInput{}, errors.New(fmt.Sprintf("unsupported installer command: %s", lineSplits[0]))
}

func convInstallLineToEditInput(lineSplits []string) installer.InstallInput {
	return installer.InstallInput{
		Command:      installer.EditCmd,
		FileNames:    strings.Split(lineSplits[1], ","),
		KeyLocations: strings.Split(lineSplits[2], ","),
		NewValue:     lineSplits[3],
	}
}

func convInstallLineToDropInput(lineSplits []string) installer.InstallInput {
	return installer.InstallInput{
		Command:   installer.DropCmd,
		FileNames: strings.Split(lineSplits[1], ","),
	}
}

func formatCmdLine(cmdLine string) string {
	cmdLine = strings.TrimSpace(cmdLine)
	return standardizeSpaces(cmdLine)
}

func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func isLineValid(cmdLine string) bool {
	lineSplits := strings.Split(cmdLine, " ")
	firstSplit := lineSplits[0]
	if isValidCommand(firstSplit) {
		return true
	}
	if isValidForSkip(firstSplit) {
		return false
	}
	msg := errors.New(fmt.Sprintf("skiping invalid install line - %s", cmdLine))
	common.Logger.Println(&common.Warning{Reason: msg})
	return false
}

func isValidForSkip(firstSplit string) bool {
	return firstSplit == "" || strings.Contains(firstSplit, "#")
}

func isValidCommand(str string) bool {
	return str == installer.EditCmd.String() || str == installer.DropCmd.String()
}

func composeArgoAppInstallInput(flags *common.Flags) []installer.InstallInput {
	var installInputs []installer.InstallInput
	installInputs = append(installInputs, getEditNameInput(flags))
	installInputs = append(installInputs, getEditPathInput(flags))
	installInputs = append(installInputs, getEditNamespaceInput(flags))
	return installInputs
}

func getEditNamespaceInput(flags *common.Flags) installer.InstallInput {
	return installer.InstallInput{
		Command:      installer.EditCmd,
		FileNames:    []string{fmt.Sprintf("%s.yaml", flags.App.DestBranch)},
		KeyLocations: []string{"spec.destination.namespace"},
		NewValue:     flags.App.Namespace,
	}
}

func getEditPathInput(flags *common.Flags) installer.InstallInput {
	destinationDir := fmt.Sprintf("%s/%s", flags.ComposeYamlPath(), flags.App.DestBranch)
	return installer.InstallInput{
		Command:      installer.EditCmd,
		FileNames:    []string{fmt.Sprintf("%s.yaml", flags.App.DestBranch)},
		KeyLocations: []string{"spec.source.path"},
		NewValue:     destinationDir,
	}
}

func getEditNameInput(flags *common.Flags) installer.InstallInput {
	return installer.InstallInput{
		Command:      installer.EditCmd,
		FileNames:    []string{fmt.Sprintf("%s.yaml", flags.App.DestBranch)},
		KeyLocations: []string{"metadata.name"},
		NewValue:     fmt.Sprintf("%s-%s", flags.App.Name, flags.App.DestBranch),
	}
}
