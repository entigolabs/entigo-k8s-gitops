package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/argoproj/argo-cd/v3/pkg/apiclient"
	applicationpkg "github.com/argoproj/argo-cd/v3/pkg/apiclient/application"
	"github.com/argoproj/argo-cd/v3/pkg/apis/application/v1alpha1"
	argoappv1 "github.com/argoproj/argo-cd/v3/pkg/apis/application/v1alpha1"
	"github.com/argoproj/argo-cd/v3/util/io"
	"github.com/argoproj/gitops-engine/pkg/health"
	argocdtypes "github.com/argoproj/gitops-engine/pkg/sync/common"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
	"k8s.io/apimachinery/pkg/watch"
	"strconv"
	"strings"
	"time"
)

const (
	initialRetryDelay int = 1
	maxRetryDelay     int = 30
)

type Client interface {
	DeleteRequest(applicationName string, cascade bool, timeout int)
	GetRequest(applicationName string, timeout int, refresh bool) *v1alpha1.Application
	SyncRequest(applicationName string, timeout int) *v1alpha1.Application
	WaitApplicationSync(applicationName string, timeout int, resourceVersion string, waitFailure bool)
}

type client struct {
	ArgoCDClient apiclient.Client
}

func NewClientOrDie(flags *common.Flags) Client {
	common.Logger.Printf("ArgoCD server: %s\n", flags.ArgoCD.ServerAddr)
	options := apiclient.ClientOptions{
		ServerAddr: flags.ArgoCD.ServerAddr,
		Insecure:   flags.ArgoCD.Insecure,
		AuthToken:  flags.ArgoCD.AuthToken,
		GRPCWeb:    flags.ArgoCD.GRPCWeb,
	}
	var (
		client client
		err    error
	)
	client.ArgoCDClient, err = apiclient.NewClient(&options)
	if err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
	return &client
}

func (c *client) GetRequest(applicationName string, timeout int, refresh bool) *v1alpha1.Application {
	common.Logger.Printf("Getting ArgoCD app: %s, timeout: %d seconds\n", applicationName, timeout)
	ctx, cancel := context.WithCancel(context.Background())
	if timer := addFatalTimeout(timeout, cancel, "Getting application timed out"); timer != nil {
		defer timer.Stop()
	}
	defer cancel()
	retryDelay := initialRetryDelay
	for ctx.Err() == nil {
		connection, applicationServiceClient := c.newApplicationClient()
		if applicationServiceClient != nil {
			refreshString := strconv.FormatBool(refresh)
			appQuery := &applicationpkg.ApplicationQuery{
				Name:    &applicationName,
				Refresh: &refreshString,
			}
			app, err := applicationServiceClient.Get(ctx, appQuery)
			closeConnection(connection)
			if err == nil {
				return app
			} else {
				status, _ := grpcstatus.FromError(err)
				switch status.Code() {
				case codes.Unavailable:
					common.Logger.Println("Connection failed")
				case codes.NotFound:
					return nil
				default:
					common.Logger.Fatal(&common.PrefixedError{Reason: err})
				}
			}
		}
		sleep(&retryDelay)
	}
	common.Logger.Fatal(&common.PrefixedError{Reason: ctx.Err()})
	return nil
}

func (c *client) SyncRequest(applicationName string, timeout int) *v1alpha1.Application {
	common.Logger.Printf("Syncing ArgoCD app: %s, timeout: %d seconds\n", applicationName, timeout)
	ctx, cancel := context.WithCancel(context.Background())
	if timer := addFatalTimeout(timeout, cancel, "Syncing application timed out"); timer != nil {
		defer timer.Stop()
	}
	defer cancel()
	prune := true
	syncReq := applicationpkg.ApplicationSyncRequest{
		Name:  &applicationName,
		Prune: &prune,
	}
	syncReq.Strategy = &argoappv1.SyncStrategy{Hook: &argoappv1.SyncStrategyHook{}}
	syncReq.Strategy.Hook.Force = true
	return c.sendSyncReq(ctx, syncReq)
}

func (c *client) sendSyncReq(ctx context.Context, syncReq applicationpkg.ApplicationSyncRequest) *v1alpha1.Application {
	retryDelay := initialRetryDelay
	for ctx.Err() == nil {
		connection, applicationServiceClient := c.newApplicationClient()
		if applicationServiceClient != nil {
			application, err := applicationServiceClient.Sync(ctx, &syncReq)
			closeConnection(connection)
			if err == nil {
				return application
			} else {
				status, _ := grpcstatus.FromError(err)
				switch status.Code() {
				case codes.FailedPrecondition:
					common.Logger.Println("Another operation in progress")
				case codes.Unavailable:
					common.Logger.Println("Connection failed")
				default:
					common.Logger.Fatal(&common.PrefixedError{Reason: err})
				}
			}
		}
		sleep(&retryDelay)
	}
	common.Logger.Fatal(&common.PrefixedError{Reason: ctx.Err()})
	return nil
}

func (c *client) DeleteRequest(applicationName string, cascade bool, timeout int) {
	common.Logger.Printf("Deleting ArgoCD app: %s, timeout: %d seconds\n", applicationName, timeout)
	ctx, cancel := context.WithCancel(context.Background())
	if timer := addFatalTimeout(timeout, cancel, "Deleting application timed out"); timer != nil {
		defer timer.Stop()
	}
	defer cancel()
	deleteReq := applicationpkg.ApplicationDeleteRequest{
		Name:    &applicationName,
		Cascade: &cascade,
	}
	retryDelay := initialRetryDelay
	for ctx.Err() == nil {
		connection, applicationServiceClient := c.newApplicationClient()
		if applicationServiceClient != nil {
			_, err := applicationServiceClient.Delete(ctx, &deleteReq)
			closeConnection(connection)
			if err == nil {
				return
			} else {
				status, _ := grpcstatus.FromError(err)
				switch status.Code() {
				case codes.Unavailable:
					common.Logger.Println("Connection failed")
				default:
					common.Logger.Fatal(&common.PrefixedError{Reason: err})
				}
			}
		}
		sleep(&retryDelay)
	}
	common.Logger.Fatal(&common.PrefixedError{Reason: ctx.Err()})
}

func (c *client) WaitApplicationSync(applicationName string, timeout int, resourceVersion string, waitFailure bool) {
	ctx, cancel := context.WithCancel(context.Background())
	timer := addTimeout(timeout, cancel, "Waiting for application to sync timed out", waitFailure)
	if timer != nil {
		defer timer.Stop()
	}
	defer cancel()

	appEventCh := c.ArgoCDClient.WatchApplicationWithRetry(ctx, applicationName, resourceVersion)
	var statusMessage string // Holds last statusMessage to avoid logging duplicates
	for appEvent := range appEventCh {
		if isApplicationReady(appEvent, &statusMessage) {
			common.Logger.Println("Application is synced and healthy")
			return
		}
	}
	if waitFailure && ctx.Err() != context.Canceled {
		common.Logger.Fatal(&common.PrefixedError{Reason: errors.New("waiting was interrupted before sync completed")})
	}
}

func addFatalTimeout(timeout int, cancel context.CancelFunc, message string) *time.Timer {
	return addTimeout(timeout, cancel, message, true)
}

func addTimeout(timeout int, cancel context.CancelFunc, message string, fatal bool) *time.Timer {
	if timeout > 0 {
		return time.AfterFunc(time.Duration(timeout)*time.Second, func() {
			if fatal {
				common.Logger.Fatal(&common.PrefixedError{Reason: errors.New(message)})
			} else {
				cancel()
				common.Logger.Println(&common.Warning{Reason: errors.New(message)})
			}
		})
	} else {
		return nil
	}
}

func (c *client) newApplicationClient() (io.Closer, applicationpkg.ApplicationServiceClient) {
	connection, applicationServiceClient, err := c.ArgoCDClient.NewApplicationClient()
	if err == nil {
		return connection, applicationServiceClient
	} else {
		if err.Error() == "EOF" {
			common.Logger.Println(&common.Warning{Reason: errors.New("failed to connect to the server")})
		} else {
			common.Logger.Fatal(&common.PrefixedError{Reason: err})
		}
	}
	return nil, nil
}

func closeConnection(conn io.Closer) {
	if conn != nil {
		if err := conn.Close(); err != nil {
			common.Logger.Println(&common.Warning{Reason: err})
		}
	}
}

func sleep(retryDelay *int) {
	common.Logger.Printf("Retrying in %d seconds\n", *retryDelay)
	time.Sleep(time.Duration(*retryDelay) * time.Second)
	if *retryDelay < maxRetryDelay {
		*retryDelay = common.MinInt(maxRetryDelay, *retryDelay*2)
	}
}

func isApplicationReady(applicationEvent *v1alpha1.ApplicationWatchEvent, lastStatusMessage *string) bool {
	checkEventType(applicationEvent)
	application := applicationEvent.Application
	if application.Operation != nil {
		return false
	} else if application.Status.OperationState != nil {
		status := application.Status
		operationState := status.OperationState
		if operationState.Phase == argocdtypes.OperationFailed {
			syncFailure(operationState)
		} else if operationState.FinishedAt == nil || (status.ReconciledAt == nil ||
			status.ReconciledAt.Before(operationState.FinishedAt)) {
			logStatus(getStatus(application, true), lastStatusMessage)
			return false
		}
	}

	healthStatus := string(application.Status.Health.Status)
	syncStatus := string(application.Status.Sync.Status)
	logStatus(getStatus(application, false), lastStatusMessage)
	return healthStatus == string(health.HealthStatusHealthy) && syncStatus == string(argoappv1.SyncStatusCodeSynced)
}

func checkEventType(applicationEvent *v1alpha1.ApplicationWatchEvent) {
	switch applicationEvent.Type {
	case watch.Error:
		common.Logger.Fatal(applicationEvent)
	case watch.Deleted:
		common.Logger.Fatal(&common.PrefixedError{Reason: errors.New("application was deleted while waiting")})
	}
}

func syncFailure(operationState *argoappv1.OperationState) {
	common.Logger.Println(&common.PrefixedError{Reason: errors.New("argoCD operation failed with message: " + operationState.Message)})
	result := operationState.SyncResult
	if result != nil && result.Resources != nil {
		for i := range result.Resources {
			resource := result.Resources[i]
			if resource.Status == argocdtypes.ResultCodeSyncFailed {
				common.Logger.Printf("%s (%s) - %s, %s", resource.Name,
					resource.Kind, resource.Status, resource.Message)
			}
		}
	}
	common.Logger.Fatal(&common.PrefixedError{Reason: errors.New("application failed to reach synced and ready status")})
}

func logStatus(message string, lastStatusMessage *string) {
	if message != *lastStatusMessage {
		common.Logger.Println(message)
		*lastStatusMessage = message
	}
}

func getStatus(application argoappv1.Application, operationInProgress bool) string {
	var sb strings.Builder
	if operationInProgress {
		sb.WriteString("Operation in progress")
	} else {
		sb.WriteString("Operation finished")
	}
	sb.WriteString(", resource statuses: ")
	sb.WriteString(getResourceStatuses(application.Status.Resources))
	return sb.String()
}

func getResourceStatuses(resources []v1alpha1.ResourceStatus) string {
	if len(resources) == 0 {
		return "no resources found"
	} else {
		var statuses []string
		for i := range resources {
			resource := resources[i]
			if string(resource.Status) != string(argoappv1.SyncStatusCodeSynced) {
				statuses = append(statuses, getResourceStatus(resource, string(resource.Status)))
			} else if resource.Health != nil && string(resource.Health.Status) != string(health.HealthStatusHealthy) {
				statuses = append(statuses, getResourceStatus(resource, string(resource.Health.Status)))
			}
		}
		if len(statuses) == 0 {
			return "all ready"
		} else {
			return strings.Join(statuses, "; ")
		}
	}
}

func getResourceStatus(resource v1alpha1.ResourceStatus, status string) string {
	return fmt.Sprintf("%s (%s) - %s", resource.Name, resource.Kind, status)
}
