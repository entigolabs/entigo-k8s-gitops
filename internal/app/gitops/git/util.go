package git

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/go-git/go-git/v5"
	goGitSsh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"golang.org/x/crypto/ssh"
	"os"
	"path/filepath"
	"strings"
)

func (r *Repository) DoesRepoExist() bool {
	if _, err := git.PlainInit(common.GetRepositoryName(r.Repo), false); err != nil {
		switch err {
		case git.ErrRepositoryAlreadyExists:
			return true
		default:
			common.Logger.Fatal(&common.PrefixedError{Reason: err})
		}
	}
	removeRepoFolder(r.Repo)
	return false
}

func removeRepoFolder(repoSshUrl string) {
	if err := os.RemoveAll(common.GetRepositoryRootPath(repoSshUrl)); err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
}

func isRemoteKeyDefined(keyFile string) bool {
	keyAbsPath := getKeyFileAbsPath(keyFile)
	if _, err := os.Stat(keyAbsPath); err != nil {
		common.Logger.Println(fmt.Sprintf("coldn't use SSH key defined via flag. %s", err))
		common.Logger.Println("using SSH key defined in pipeline")
		return false
	}
	common.Logger.Println(fmt.Sprintf("using SSH key defined in %s", keyFile))
	return true
}

func getKeyFileAbsPath(keyPath string) string {
	if filepath.IsAbs(keyPath) {
		return keyPath
	}
	if err := common.ChangeDir(common.RootPath); err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
	absKeyPath, err := filepath.Abs(keyPath)
	if err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
	return absKeyPath
}

func (r *Repository) getPublicKeys() *goGitSsh.PublicKeys {
	if r.StrictHostKeyChecking {
		return getPublicKeysDefault(r.KeyFile)
	}
	return getPublicKeysNonStrict(r.KeyFile)
}

func getPublicKeysDefault(keyFile string) *goGitSsh.PublicKeys {
	publicKeys, err := goGitSsh.NewPublicKeysFromFile("git", keyFile, "")
	if err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
	return publicKeys
}

func getPublicKeysNonStrict(keyFile string) *goGitSsh.PublicKeys {
	publicKeys := getPublicKeysDefault(keyFile)
	publicKeys.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	return publicKeys
}

func getAppName(appPath string) string {
	pathTokens := strings.Split(appPath, "/")
	return pathTokens[len(pathTokens)-1]
}
