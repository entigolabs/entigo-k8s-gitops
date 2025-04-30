package git

import (
	"errors"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	goGitSsh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"golang.org/x/crypto/ssh"
	"os"
	"path/filepath"
)

func (r *Repository) DoesRepoExist() bool {
	if _, err := git.PlainInit(common.GetRepositoryName(r.Repo), false); err != nil {
		if errors.Is(err, git.ErrRepositoryAlreadyExists) {
			return true
		}
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
	removeRepoFolder(r.Repo)
	return false
}

func removeRepoFolder(repoSshUrl string) {
	if err := os.RemoveAll(common.GetRepositoryRootPath(repoSshUrl)); err != nil {
		common.Logger.Fatal(&common.PrefixedError{Reason: err})
	}
}

func (r *Repository) isRemoteKeyDefined(keyFile string) bool {
	if r.KeyFile == "" {
		return false
	}
	keyAbsPath := r.getKeyFileAbsPath(keyFile)
	if _, err := os.Stat(keyAbsPath); err != nil {
		common.Logger.Printf("coldn't use SSH key defined via flag. %s\n", err)
		common.Logger.Println("using SSH key defined in pipeline")
		return false
	}
	common.Logger.Printf("using SSH key defined in %s\n", keyFile)
	return true
}

func (r *Repository) getKeyFileAbsPath(keyPath string) string {
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
	r.KeyFile = absKeyPath
	common.CdToGitOpsWd()
	return absKeyPath
}

func (r *Repository) getPublicKeys() *goGitSsh.PublicKeys {
	if r.StrictHostKeyChecking {
		return getPublicKeysDefault(r.KeyFile)
	}
	return getPublicKeysNonStrict(r.KeyFile)
}

func (r *Repository) getBasicAuth() transport.AuthMethod {
	return &http.BasicAuth{
		Username: r.Username,
		Password: r.Password,
	}
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
