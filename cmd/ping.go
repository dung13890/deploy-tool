package cmd

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

type ping struct {
}

func PingInit() *cli.Command {
	return &cli.Command{
		Name:    "ping",
		Aliases: []string{"p"},
		Usage:   "Testing connection into servers",
		Action: func(c *cli.Context) error {
			fmt.Println("completed task: ", c.Args().First())
			return nil
		},
	}
}
