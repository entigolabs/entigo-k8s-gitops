package main

import (
	"fmt"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/copy"
	"github.com/entigolabs/entigo-k8s-gitops/internal/app/gitops/update"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"sort"
)

func runCli() {
	app := &cli.App{
		Flags:    cliFlags(),
		Commands: cliCommands(),
		Action:   defaultAction(),
	}
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err) // TODO unified logging
	}
}

func defaultAction() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		fmt.Println("defaultAction ->", flags)
		//cli.ShowAppHelp(c) // manually call help
		return nil
	}
}

func cliFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "repo",
			Aliases:     []string{"r"},
			DefaultText: "repoDefaultVal",
			Usage:       "Git repository `SSH URL`",
			Destination: &flags.Repo,
		},
		&cli.StringFlag{
			Name:        "branch",
			Aliases:     []string{"b"},
			DefaultText: "branchDefaultVal",
			Usage:       "Branch `name`",
		},
	}
}

func cliCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:    "update",
			Aliases: []string{"up"},
			Usage:   "Update corresponding images",
			Action: func(c *cli.Context) error {
				fmt.Println("update!")
				update.Run(flags)
				return nil
			},
			Flags: cliFlags(),
		},
		{
			Name:    "copy",
			Aliases: []string{"cp"},
			Usage:   "Copy from master branch",
			Action: func(c *cli.Context) error {
				fmt.Println("update!")
				copy.Run(flags)
				return nil
			},
			Flags: cliFlags(),
		},
	}
}
