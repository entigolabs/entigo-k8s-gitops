package update

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/util"
	"github.com/go-git/go-git/v5"
)

const OperationType = "update"
const GitOpsWd = "gitops-workdir"

func Update() func() {
	return func() {
		setupFlags()
		cloneIfNecessary()
		updateImages()
		applyChanges()
		pushIfWanted()
		printExitMessage()
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
	commitIfModified(openedRepo)
}

func commitIfModified(openedRepo *git.Repository) {
	if isRepoModified(openedRepo) {
		gitCommit(openedRepo)
	} else {
		util.Logger.Println("nothing to commit, working tree clean")
	}
}

func pushIfWanted() {
	if flgs.push {
		openedRepo := openGitOpsRepo()
		gitPush(openedRepo)
	} else {
		util.Logger.Println("commit(s) were chosen not to be pushed")
	}
}

func printExitMessage() {
	if flgs.push {
		util.Logger.Println(fmt.Sprintf("repository url: %s", util.GetWebUrl(flgs.repo)))
	}
}
