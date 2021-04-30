package sync

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/argocd/api"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
)

func Run(flags *common.Flags) {
	appName := flags.App.Name
	timeout := flags.ArgoCD.Timeout
	client := api.NewClientOrDie(flags)
	application := client.SyncRequest(appName, timeout)
	if !flags.ArgoCD.Async {
		common.Logger.Println(fmt.Sprintf("Waiting for application to sync, timeout: %d seconds", timeout))
		client.WaitApplicationSync(appName, timeout, application.ResourceVersion, flags.ArgoCD.WaitFailure)
	} else {
		common.Logger.Println("Waiting disabled, won't wait for sync to complete")
	}
}
