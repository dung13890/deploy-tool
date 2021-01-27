package cmd

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

type deploy struct {
}

func DeployInit() *cli.Command {
	return &cli.Command{
		Name:    "deploy",
		Aliases: []string{"d"},
		Usage:   "Deploy into servers",
		Action: func(c *cli.Context) error {
			fmt.Println("completed task: ", c.Args().First())
			return nil
		},
	}
}
