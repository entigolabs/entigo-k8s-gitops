package installer

import (
	"errors"
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
)

type Installer struct {
	Command            common.Command
	KeepRegistry       bool
	DeploymentStrategy common.DeploymentStrategy
	InstallHistory     []InstallHistory
}

type InstallHistory struct {
	NewValue      string
	OldValue      string
	workingFile   string
	documentIndex int
	workingKey    string
}

type InstallCommand int

const (
	EditCmd InstallCommand = iota
	DropCmd
)

func (ic InstallCommand) String() string {
	return [...]string{"edit", "drop"}[ic]
}

type InstallInput struct {
	Command      InstallCommand
	FileNames    []string
	KeyLocations []string
	NewValue     string
}

func (i *Installer) Install(installInputs []InstallInput) []InstallHistory {
	for _, installInput := range installInputs {
		editInfo.documentIndex = 0
		i.installSingleInput(installInput)
	}
	return i.InstallHistory
}

func (i *Installer) installSingleInput(input InstallInput) {
	switch input.Command {
	case EditCmd:
		i.edit(input)
	case DropCmd:
		i.drop(input)
	default:
		msg := fmt.Sprintf("unsupported command '%s'", input.Command.String())
		common.Logger.Fatal(common.PrefixedError{Reason: errors.New(msg)})
	}
	common.Logger.Printf("finised %s command\n", input.Command.String())
}

func ConvStrToInstallCommand(str string) InstallCommand {
	switch str {
	case EditCmd.String():
		return EditCmd
	case DropCmd.String():
		return DropCmd
	default:
		msg := fmt.Sprintf("unsupported command '%s'", str)
		common.Logger.Fatal(common.PrefixedError{Reason: errors.New(msg)})
	}
	return EditCmd
}
