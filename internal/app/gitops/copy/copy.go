package copy

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
)

func Run(flags *common.Flags) {
	fmt.Println("copy:", flags)
}
