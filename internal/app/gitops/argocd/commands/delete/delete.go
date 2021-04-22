package sync

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/argocd/api"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
)

func Run(flags *common.Flags) {
	var appName = flags.App.Name
	common.Logger.Println(fmt.Sprintf("Deleting ArgoCD app: %s", appName))
	client := api.NewClientOrDie(flags)
	defer client.Close()
	client.DeleteRequest(appName, flags.ArgoCD.Cascade)
}
