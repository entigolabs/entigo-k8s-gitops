package cli

import (
	"errors"
	argoCDDelete "github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/argocd/commands/delete"
	argoCDSync "github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/argocd/commands/sync"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/commands/copy"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/commands/delete"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/commands/update"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/urfave/cli/v2"
)

func action(cmd common.Command) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		if err := flags.Setup(cmd); err != nil {
			common.Logger.Fatal(&common.PrefixedError{Reason: err})
		}
		run(cmd)
		return nil
	}
}

func run(cmd common.Command) {
	switch cmd {
	case common.RunCmd:
		update.Run(flags)
		copy.Run(flags)
	case common.UpdateCmd:
		update.Run(flags)
	case common.CopyCmd:
		copy.Run(flags)
	case common.DeleteCmd:
		delete.Run(flags)
	case common.ArgoCDSyncCmd:
		argoCDSync.Run(flags)
	case common.ArgoCDDeleteCmd:
		argoCDDelete.Run(flags)
	default:
		common.Logger.Fatal(&common.PrefixedError{Reason: errors.New("unsupported command")})
	}

}
