package cli

import (
	"errors"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/copy"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/update"
	"github.com/urfave/cli/v2"
)

func action(cmd Command) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		if err := flags.Setup(); err != nil {
			common.Logger.Fatal(&common.PrefixedError{Reason: err})
		}
		run(cmd)
		return nil
	}
}

func run(cmd Command) {
	switch cmd {
	case runCmd:
		update.Run(flags)
		copy.Run(flags)
	case updateCmd:
		update.Run(flags)
	case copyCmd:
		copy.Run(flags)
	default:
		common.Logger.Fatal(&common.PrefixedError{Reason: errors.New("unsupported command")})
	}

}
