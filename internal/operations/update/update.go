package update

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/util"
)

const OperationType = "update"
const GitOpsWd = "gitops-workdir"

func Update() func() {
	return func() {
		evaluateFlags()
		cloneIfNecessary()
		updateImages()
		applyChanges()
		util.Logger.Println(fmt.Sprintf("repository url: %s", util.GetWebUrl(flgs.repo)))
	}
}

func cloneIfNecessary() {
	cdToGitOpsWd()
	if !doesRepoExist() {
		cloneAndConfig()
	} else {
		util.Logger.Println("repo already exist, skip cloning")
	}
}

func cloneAndConfig() {
	repo := gitClone()
	repo = configRepo(repo)
}

func updateImages() {
	cdToAppDir()
	changeImages()
}

func applyChanges() {
	openedRepo := openGitOpsRepo()
	gitAdd(openedRepo)
	exitIfUnmodified(openedRepo)
	gitCommit(openedRepo)
	gitPush(openedRepo)
}
