package api

import (
	"context"
	"fmt"
	"github.com/argoproj/argo-cd/pkg/apiclient"
	applicationpkg "github.com/argoproj/argo-cd/pkg/apiclient/application"
	"github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	argoappv1 "github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	"github.com/argoproj/argo-cd/util/io"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
)

type Client interface {
	DeleteRequest(applicationName string, cascade bool)
	SyncRequest(applicationName string) *v1alpha1.Application
	Close()
}

type client struct {
	ApplicationServiceClient applicationpkg.ApplicationServiceClient
	Connection               io.Closer
	Context                  context.Context
}

func NewClientOrDie(flags *common.Flags) Client {
	common.Logger.Println(fmt.Sprintf("ArgoCD server: %s", flags.ArgoCD.ServerAddr))
	options := apiclient.ClientOptions{
		ServerAddr: flags.ArgoCD.ServerAddr,
		Insecure:   flags.ArgoCD.Insecure,
		AuthToken:  flags.ArgoCD.AuthToken,
	}
	argoCDClient, err := apiclient.NewClient(&options)
	if err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
	var client client
	client.Connection, client.ApplicationServiceClient, err = argoCDClient.NewApplicationClient()
	if err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
	client.Context = context.Background()
	return &client
}

func (c *client) SyncRequest(applicationName string) *v1alpha1.Application {
	// TODO Retry strategy
	// TODO Dry run?
	syncReq := applicationpkg.ApplicationSyncRequest{
		Name:  &applicationName,
		Prune: true,
	}
	syncReq.Strategy = &argoappv1.SyncStrategy{Apply: &argoappv1.SyncStrategyApply{}}
	syncReq.Strategy.Apply.Force = true
	application, err := c.ApplicationServiceClient.Sync(c.Context, &syncReq)
	if err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
	return application
}

func (c *client) DeleteRequest(applicationName string, cascade bool) {
	// TODO Retry strategy
	// TODO Dry run?
	deleteReq := applicationpkg.ApplicationDeleteRequest{
		Name:    &applicationName,
		Cascade: &cascade,
	}
	// TODO Is returned ApplicationContext useful?
	_, err := c.ApplicationServiceClient.Delete(c.Context, &deleteReq)
	if err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
}

func (c *client) Close() {
	io.Close(c.Connection)
}
