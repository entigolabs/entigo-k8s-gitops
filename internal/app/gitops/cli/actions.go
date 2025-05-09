package cli

import (
	"context"
	"errors"
	argoCDDelete "github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/argocd/commands/delete"
	argoCDGet "github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/argocd/commands/get"
	argoCDSync "github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/argocd/commands/sync"
	argoCDUpdate "github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/argocd/commands/update"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/commands/copy"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/commands/delete"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/commands/notify"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/commands/update"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/urfave/cli/v3"
)

func action(cmd common.Command) cli.ActionFunc {
	return func(c context.Context, _ *cli.Command) error {
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
	case common.DeploymentNotificationCmd:
		notify.Run(flags)
	case common.ArgoCDGetCmd:
		argoCDGet.Run(flags)
	case common.ArgoCDSyncCmd:
		argoCDSync.Run(flags)
	case common.ArgoCDUpdateCmd:
		argoCDUpdate.Run(flags)
	case common.ArgoCDDeleteCmd:
		argoCDDelete.Run(flags)
	case common.VersionCmd:
		common.PrintVersion()
	default:
		common.Logger.Fatal(&common.PrefixedError{Reason: errors.New("unsupported command")})
	}
}
