package update

import (
	"flag"
	"fmt"
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
	prefix                string
	appNamespace          string
	appName               string
}

func setupFlags() {
	evaluateFlags()
	setupAppPath()
}

func setupAppPath() {
	if flgs.appPath == "" {
		util.Logger.Println("using tokenized path")
		composeAppPath()
	} else {
		util.Logger.Println("--app-path overrides tokenized flags")
	}
}

func composeAppPath() {
	flgs.appPath = fmt.Sprintf("%s/%s/%s", flgs.prefix, flgs.appNamespace, flgs.appName)
}

func evaluateFlags() {
	repo, branch, keyFile, strictHostKeyChecking, push, images, appPath, prefix, appNamespace, appName := parseFlags()
	flgs = Flags{
		repo:                  *repo,
		branch:                *branch,
		keyFile:               *keyFile,
		strictHostKeyChecking: *strictHostKeyChecking,
		push:                  *push,
		images:                *images,
		appPath:               *appPath,
		prefix:                *prefix,
		appNamespace:          *appNamespace,
		appName:               *appName,
	}
}

func parseFlags() (*string, *string, *string, *bool, *bool, *string, *string, *string, *string, *string) {
	flagSet := flag.NewFlagSet("Update Flag Set", flag.ExitOnError)
	repo := flagSet.String("repo", "", "Git repository SSH URL")
	branch := flagSet.String("branch", "", "repository branch")
	keyFile := flagSet.String("key-file", "", "SSH private key location")
	strictHostKeyChecking := flagSet.Bool("strict-host-key-checking", false, "strict host key checking boolean")
	push := flagSet.Bool("push", true, "git push boolean")
	images := flagSet.String("images", "", "image(s) with tag(s)")
	appPath := flagSet.String("app-path", "", "path to application folder")
	prefix := flagSet.String("prefix", "", "path prefix to apply")
	appNamespace := flagSet.String("app-namespace", "", "application namespace")
	appName := flagSet.String("app-name", "", "application name")
	parseErr := flagSet.Parse(os.Args[2:])
	if parseErr != nil {
		util.Logger.Println(&util.PrefixedError{Reason: parseErr})
		os.Exit(1)
	}
	return repo, branch, keyFile, strictHostKeyChecking, push, images, appPath, prefix, appNamespace, appName
}
