package get

import (
	"encoding/json"
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/argocd/api"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
)

func Run(flags *common.Flags) {
	appName := flags.App.Name
	client := api.NewClientOrDie(flags)
	app := client.GetRequest(appName, flags.ArgoCD.Timeout, flags.ArgoCD.Refresh)
	appJson, err := json.Marshal(app)
	if err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
	fmt.Println(string(appJson))
}
