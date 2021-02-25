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
		installInputs = append(installInputs, installLineToInstallInput(line))
	}
	return installInputs
}

func installLineToInstallInput(installLine string) installer.InstallInput {
	lineSplits := strings.Split(installLine, " ")
	return installer.InstallInput{
		Command:      installer.ConvStrToInstallCommand(lineSplits[0]),
		FileNames:    strings.Split(lineSplits[1], ","),
		KeyLocations: strings.Split(lineSplits[2], ","),
		NewValue:     lineSplits[3],
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

func isValidCommand(firstSplit string) bool {
	return firstSplit == installer.EditCmd.String() || firstSplit == installer.DropCmd.String()
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
		FileNames:    []string{fmt.Sprintf("%s.yaml", flags.App.Branch)},
		KeyLocations: []string{"spec.destination.namespace"},
		NewValue:     flags.App.Namespace,
	}
}

func getEditPathInput(flags *common.Flags) installer.InstallInput {
	destinationDir := fmt.Sprintf("%s/%s", flags.ComposeYamlPath(), flags.App.Branch)
	return installer.InstallInput{
		Command:      installer.EditCmd,
		FileNames:    []string{fmt.Sprintf("%s.yaml", flags.App.Branch)},
		KeyLocations: []string{"spec.source.path"},
		NewValue:     destinationDir,
	}
}

func getEditNameInput(flags *common.Flags) installer.InstallInput {
	return installer.InstallInput{
		Command:      installer.EditCmd,
		FileNames:    []string{fmt.Sprintf("%s.yaml", flags.App.Branch)},
		KeyLocations: []string{"metadata.name"},
		NewValue:     fmt.Sprintf("%s-%s", flags.App.Name, flags.App.Branch),
	}
}
