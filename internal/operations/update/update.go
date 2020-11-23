package update

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/util"
	"runtime"
	"strings"
)

const OperationType = "update"
const GitOpsWd = "gitops"

func Update() func() {
	return func() {
		evaluateFlags()
		cloneRepo()
		updateImages()
		applyChanges()
		util.Logger.Println(fmt.Sprintf("repository url: %s", util.GetWebUrl(flgs.repo)))
	}
}

func cloneRepo() {
	cdToGitOps()
	repo := gitClone()
	repo = configRepo(repo)
}

func updateImages() {
	cdToWorkPath()
	changeImages()
}

func applyChanges() {
	openedRepo := openGitOpsRepo()

	gitAdd(openedRepo)
	exitIfUnmodified(openedRepo)

	gitCommit(openedRepo)
	gitPush(openedRepo)
}

func changeImages() {
	imgNameTokens := strings.Split(flgs.images, ",")
	for _, img := range imgNameTokens {
		changeSpecificImage(img)
	}
}

func changeSpecificImage(image string) {
	switch opSys := runtime.GOOS; opSys {
	case "darwin":
		changeImageOsx(image)
	default:
		changeImageDefault(image)
	}
}
