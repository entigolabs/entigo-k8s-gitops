package sync

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/argocd/api"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
)

func Run(flags *common.Flags) {
	var appName = flags.App.Name
	common.Logger.Println(fmt.Sprintf("Syncing ArgoCD app: %s", appName))
	client := api.NewClientOrDie(flags)
	defer client.Close()
	client.SyncRequest(appName)
	if !flags.ArgoCD.Async {
		// TODO Wait for healthy
	}
}
