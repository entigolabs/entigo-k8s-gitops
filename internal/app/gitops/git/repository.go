package git

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/go-git/go-git/v5"
	"strings"
)

var Repository *git.Repository = new(git.Repository)

func Clone(flags *common.Flags) {
	common.Logger.Println(fmt.Sprintf("git clone %s", flags.Git.Repo))
	repo, err := git.PlainClone(getRepositoryName(flags.Git.Repo), false, getCloneOptions(flags.Git))
	if err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
	common.Logger.Println("git clone was successful")
	Repository = repo
}

func getRepositoryName(sshUrl string) string {
	tokens := strings.Split(sshUrl, "/")
	lastToken := tokens[len(tokens)-1]
	return lastToken[:len(lastToken)-4]
}
