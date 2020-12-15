package git

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/entigolabs/entigo-k8s-gitops/internal/util"
	"github.com/go-git/go-git/v5"
	goGitSsh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"golang.org/x/crypto/ssh"
	"os"
	"strings"
)

func DoesRepoExist(repoSshUrl string) bool {
	if _, err := git.PlainInit(getRepositoryName(repoSshUrl), false); err != nil {
		switch err {
		case git.ErrRepositoryAlreadyExists:
			return true
		default:
			util.Logger.Println(&util.PrefixedError{Reason: err})
			os.Exit(1)
		}
	}
	removeRepoFolder(repoSshUrl)
	return false
}

func getRepositoryName(repoSshUrl string) string {
	tokens := strings.Split(repoSshUrl, "/")
	lastToken := tokens[len(tokens)-1]
	return lastToken[:len(lastToken)-4]
}

func removeRepoFolder(repoSshUrl string) {
	if err := os.RemoveAll(getRepositoryRootPath(repoSshUrl)); err != nil {
		util.Logger.Println(&util.PrefixedError{Reason: err})
		os.Exit(1)
	}
}

func getRepositoryRootPath(repoSshUrl string) string {
	return fmt.Sprintf("%s/%s/%s", common.RootPath, common.GitOpsWd, getRepositoryName(repoSshUrl))
}

func isRemoteKeyDefined(keyFile string) bool {
	if _, err := os.Stat(keyFile); err != nil {
		common.Logger.Println(fmt.Sprintf("coldn't use SSH key defined via flag. %s", err))
		common.Logger.Println("using SSH key defined in pipeline")
		return false
	}
	common.Logger.Println(fmt.Sprintf("using SSH key defined in %s", keyFile))
	return true
}

func getPublicKeys(gitFlags common.GitFlags) *goGitSsh.PublicKeys {
	if gitFlags.StrictHostKeyChecking {
		return getPublicKeysDefault(gitFlags.KeyFile)
	}
	return getPublicKeysNonStrict(gitFlags.KeyFile)
}

func getPublicKeysDefault(keyFile string) *goGitSsh.PublicKeys {
	publicKeys, err := goGitSsh.NewPublicKeysFromFile("git", keyFile, "")
	if err != nil {
		common.Logger.Println(&common.PrefixedError{Reason: err})
		os.Exit(1)
	}
	return publicKeys
}

func getPublicKeysNonStrict(keyFile string) *goGitSsh.PublicKeys {
	publicKeys := getPublicKeysDefault(keyFile)
	publicKeys.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	return publicKeys
}
