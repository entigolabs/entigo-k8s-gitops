package update

import (
	"fmt"
	"github.com/argoproj/argo-cd/v3/pkg/apis/application/v1alpha1"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/argocd/api"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/argocd/commands/sync"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/commands/update"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
)

func Run(flags *common.Flags) {
	appName := flags.App.Name
	client := api.NewClientOrDie(flags)
	app := client.GetRequest(appName, flags.ArgoCD.Timeout, flags.ArgoCD.Refresh)
	if app == nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: fmt.Errorf("app '%s' not found", appName)})
	}
	addFlags(flags, app)
	update.Run(flags)
	if flags.ArgoCD.Sync {
		sync.Run(flags)
	}
}

func addFlags(flags *common.Flags, app *v1alpha1.Application) {
	flags.Git.Repo = app.Spec.Source.RepoURL
	flags.Git.Branch = app.Spec.Source.TargetRevision
	flags.App.Path = app.Spec.Source.Path
}
