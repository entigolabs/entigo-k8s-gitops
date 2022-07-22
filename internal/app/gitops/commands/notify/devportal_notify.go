package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type notificationPayload struct {
	EnvName     string `json:"deploymentEnvironmentName"`
	NewImg      string `json:"newImageWithTag"`
	OldImg      string `json:"oldImageWithTag"`
	RegistryURI string `json:"registryUri"`
}

type authToken struct {
	key, value string
}

func Run(flags *common.Flags) {
	RunNotificationRequest(flags.Notification)
}

func RunNotificationRequest(notificationFlags common.DeploymentNotificationFlags) {
	req := createNotificationRequest(notificationFlags)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			common.Logger.Println(&common.Warning{Reason: err})
		}
	}(resp.Body)
	logResponseInfo(resp)
}

func logResponseInfo(resp *http.Response) {
	common.Logger.Println(fmt.Sprintf("response Status: %s", resp.Status))
	common.Logger.Println(fmt.Sprintf("response Headers: %s", resp.Header))
	body, _ := ioutil.ReadAll(resp.Body)
	common.Logger.Println(fmt.Sprintf("response Body: %s", string(body)))
}

func createNotificationRequest(notificationFlags common.DeploymentNotificationFlags) *http.Request {
	req, _ := http.NewRequest("POST", notificationFlags.URL, getPayloadBuffer(notificationFlags))
	req.Header.Set("Content-Type", "application/json")
	token := getAuthToken(notificationFlags)
	req.Header.Set(token.key, token.value)
	return req
}

func getAuthToken(notificationFlags common.DeploymentNotificationFlags) authToken {
	tokenSplits := strings.Split(notificationFlags.AuthToken, "=")
	return authToken{
		key:   tokenSplits[0],
		value: tokenSplits[1],
	}
}

func getPayloadBuffer(notificationFlags common.DeploymentNotificationFlags) *bytes.Buffer {
	body := &notificationPayload{
		EnvName:     notificationFlags.Environment,
		NewImg:      notificationFlags.NewImage,
		OldImg:      notificationFlags.OldImage,
		RegistryURI: notificationFlags.RegistryUri,
	}
	payloadBuf := new(bytes.Buffer)
	err := json.NewEncoder(payloadBuf).Encode(body)
	if err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
	return payloadBuf
}
