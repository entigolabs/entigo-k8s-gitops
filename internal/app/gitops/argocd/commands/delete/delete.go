package sync

import (
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/argocd/api"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
)

func Run(flags *common.Flags) {
	appName := flags.App.Name
	client := api.NewClientOrDie(flags)
	client.DeleteRequest(appName, flags.ArgoCD.Cascade, flags.ArgoCD.Timeout)
}
