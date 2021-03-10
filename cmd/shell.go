package cmd

import (
	"fmt"
	"github.com/dung13890/deploy-tool/config"
	"github.com/dung13890/deploy-tool/remote"
	"github.com/urfave/cli/v2"
	"log"
)

type shell struct {
	config     config.Configuration
	privateKey string
}

func NewShell() *cli.Command {
	return &cli.Command{
		Name:    "shell",
		Aliases: []string{"s"},
		Usage:   "Running shell into multiple servers",
		Flags: []cli.Flag{
			config.Load,
			config.Identity,
		},
		Action: func(ctx *cli.Context) error {
			s := &shell{}
			err := s.config.ReadFile(ctx.String("config"))
			if err != nil {
				log.Fatal(err)
			}
			s.privateKey = ctx.String("identity")
			s.exec()
			return nil
		},
	}
}

func (s *shell) exec() error {
	var r remote.Remote

	if s.config.Server.Address == "127.0.0.1" || s.config.Server.Address == "localhost" {
		r = &remote.Localhost{}
	} else {
		r = &remote.Server{}
	}
	defer r.Close()

	fmt.Println("Running shell into multiple servers:")

	return nil
}
