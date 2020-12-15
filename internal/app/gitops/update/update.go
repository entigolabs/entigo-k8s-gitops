package update

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/git"
	"os"
)

func Run(flags *common.Flags) {
	fmt.Println("update:", flags)
	printWd()
	common.CdToGitOpsWd()
	git.Clone(flags)
	printWd()

}

func printWd() {
	path, _ := os.Getwd()
	fmt.Println(path)
}
