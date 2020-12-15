package git

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/go-git/go-git/v5"
)

var Repository *git.Repository = new(git.Repository)

func Clone(gitFlags common.GitFlags) {
	common.Logger.Println(fmt.Sprintf("git clone %s", gitFlags.Repo))
	repo, err := git.PlainClone(getRepositoryName(gitFlags.Repo), false, getCloneOptions(gitFlags))
	if err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
	common.Logger.Println("git clone was successful")
	Repository = repo
}
