package update

import (
	"flag"
	"github.com/entigolabs/entigo-k8s-gitops/internal/util"
	"os"
)

var flgs = Flags{}

type Flags = struct {
	repo    string
	branch  string
	keyPath string
	images  string
	appPath string
}

func evaluateFlags() {
	repo, branch, keyPath, images, appPath := parseFlags()
	flgs = Flags{
		repo:    *repo,
		branch:  *branch,
		keyPath: *keyPath,
		images:  *images,
		appPath: *appPath,
	}
}

func parseFlags() (*string, *string, *string, *string, *string) {
	flagSet := flag.NewFlagSet("Update Flag Set", flag.ExitOnError)
	repo := flagSet.String("repo", "", "Repository clone URL")
	branch := flagSet.String("branch", "", "Repository branch")
	keyPath := flagSet.String("key-path", "", "SSH key path")
	images := flagSet.String("images", "", "New image(s)")
	appPath := flagSet.String("app-path", "", "Application path")
	parseErr := flagSet.Parse(os.Args[2:])
	if parseErr != nil {
		util.Logger.Println(&util.PrefixedError{Reason: parseErr})
		os.Exit(1)
	}
	return repo, branch, keyPath, images, appPath
}
