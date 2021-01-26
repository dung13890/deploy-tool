package cmd

import (
	"fmt"
	"github.com/dung13890/deploy-tool/config"
	"github.com/dung13890/deploy-tool/remote"
	"github.com/urfave/cli/v2"
	"log"
)

type ping struct {
	config     config.Configuration
	privateKey string
}

func PingInit() *cli.Command {
	return &cli.Command{
		Name:    "ping",
		Aliases: []string{"p"},
		Usage:   "Testing connection into servers",
		Flags: []cli.Flag{
			config.Load,
			config.Identity,
		},
		Action: func(ctx *cli.Context) error {
			p := &ping{}
			err := p.config.ReadFile(ctx.String("config"))
			if err != nil {
				log.Fatal(err)
			}
			p.privateKey = ctx.String("identity")
			p.exec()
			fmt.Printf("%#v", p)
			fmt.Println("completed task: ", ctx.Args().First())
			return nil
		},
	}
}

func (p *ping) exec() error {
	var s remote.Remote = &remote.Server{}
	s.Load(
		p.config.Server.Address,
		p.config.Server.User,
		p.config.Server.Port,
		p.config.Server.Dir,
	)
	defer s.Close()
	if err := s.Connect(p.privateKey); err != nil {
		log.Fatalf("Error: %s", err)
		return nil
	}

	fmt.Printf("%#v", s)
	return nil
}
