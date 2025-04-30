package installer

import (
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"os"
	"strings"
)

func (i *Installer) drop(input InstallInput) {
	common.Logger.Printf("started removing %s\n", strings.Join(input.FileNames, ", "))
	for _, file := range input.FileNames {
		if err := os.RemoveAll(file); err != nil {
			common.Logger.Fatal(&common.PrefixedError{Reason: err})
		}
	}
}
