package update

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/util"
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
