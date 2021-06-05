package main

import (
	"errors"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/cli"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/common"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	terminated := make(chan os.Signal, 1)
	signal.Notify(terminated, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		cli.Run()
		close(terminated)
	}()

	sig := <-terminated
	if sig != nil {
		common.Logger.Println(&common.Warning{Reason: errors.New("utility was terminated, exiting")})
	}
}
