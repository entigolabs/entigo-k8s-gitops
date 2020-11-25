package update

import (
	"flag"
	"github.com/entigolabs/entigo-k8s-gitops/internal/util"
	"os"
)

var flgs = Flags{}

type Flags = struct {
	repo                  string
	branch                string
	keyFile               string
	strictHostKeyChecking bool
	push                  bool
	images                string
	appPath               string
}

func evaluateFlags() {
	repo, branch, keyFile, strictHostKeyChecking, push, images, appPath := parseFlags()
	flgs = Flags{
		repo:                  *repo,
		branch:                *branch,
		keyFile:               *keyFile,
		strictHostKeyChecking: *strictHostKeyChecking,
		push:                  *push,
		images:                *images,
		appPath:               *appPath,
	}
}

func parseFlags() (*string, *string, *string, *bool, *bool, *string, *string) {
	flagSet := flag.NewFlagSet("Update Flag Set", flag.ExitOnError)
	repo := flagSet.String("repo", "", "Git repository SSH URL")
	branch := flagSet.String("branch", "", "repository branch")
	keyFile := flagSet.String("key-file", "", "SSH private key location")
	strictHostKeyChecking := flagSet.Bool("strict-host-key-checking", false, "strict host key checking boolean (default is false)")
	push := flagSet.Bool("push", true, "git push boolean (default is true)")
	images := flagSet.String("images", "", "image(s) with tag(s)")
	appPath := flagSet.String("app-path", "", "path to application folder")
	parseErr := flagSet.Parse(os.Args[2:])
	if parseErr != nil {
		util.Logger.Println(&util.PrefixedError{Reason: parseErr})
		os.Exit(1)
	}
	return repo, branch, keyFile, strictHostKeyChecking, push, images, appPath
}
