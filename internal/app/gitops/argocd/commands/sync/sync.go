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
	// TODO Is waitFailure boolean needed to not fail the build when waiting times out?
	if !flags.ArgoCD.Async {
		common.Logger.Println(fmt.Sprintf("Waiting for application to sync, timeout: %d seconds", timeout))
		client.WaitApplicationSync(appName, timeout, application.ResourceVersion)
	} else {
		common.Logger.Println("Waiting disabled, won't wait for sync to complete")
	}
}
