package copy

import (
	"errors"
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
)

func Run(flags *common.Flags) {
	if flags.LoggingLevel == "dev" {
		common.Logger.Println(&common.Warning{Reason: errors.New(fmt.Sprintf("copy:", flags))})
	}
	common.Logger.Fatal("copy not implemented")
}
