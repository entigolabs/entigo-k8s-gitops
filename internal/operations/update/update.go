package update

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/util"
	"github.com/go-git/go-git/v5"
	"os"
)

const OperationType = "update"
const GitOpsWd = "gitops-workdir"

func Update() func() {
	return func() {
		setupFlags()
		cloneOrPull()
		updateImages()
		applyChanges()
		pushOnDemand()
		printExitMessage()
	}
}

func resetAndUpdate() {
	reset()
	cloneOrPull()
	updateImages()
	applyChanges()
	pushOnDemand()
}

func cloneOrPull() {
	cdToGitOpsWd()
	if !doesRepoExist() {
		cloneAndConfig()
	} else {
		openedRepo := openGitOpsRepo()
		gitPull(openedRepo)
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

func pushOnDemand() {
	if flgs.push {
		openedRepo := openGitOpsRepo()
		gitPush(openedRepo)
	} else {
		util.Logger.Println("commit(s) were chosen not to be pushed")
	}
}

func reset() {
	if err := util.ChangeDir(util.RootPath); err != nil {
		util.Logger.Println(&util.PrefixedError{Reason: err})
		os.Exit(1)
	}
	if err := os.RemoveAll(fmt.Sprintf("%s/%s", util.RootPath, GitOpsWd)); err != nil {
		util.Logger.Println(&util.PrefixedError{Reason: err})
		os.Exit(1)
	}
	util.Logger.Println("gitops-workdir successfully cleaned")
}

func printExitMessage() {
	if flgs.push {
		util.Logger.Println(fmt.Sprintf("repository url: %s", util.GetWebUrl(flgs.repo)))
	}
}
