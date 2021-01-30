package cmd

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/dung13890/deploy-tool/config"
	"github.com/dung13890/deploy-tool/remote"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"log"
	"time"
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
			return nil
		},
	}
}

func (p *ping) exec() error {
	var s remote.Remote = &remote.Server{}
	defer s.Close()
	s.Load(
		p.config.Server.Address,
		p.config.Server.User,
		p.config.Server.Port,
		p.config.Server.Dir,
	)
	fmt.Println("Testing connection into servers:")
	green := color.New(color.FgHiGreen).SprintFunc()
	sp := spinner.New(spinner.CharSets[24], 100*time.Millisecond)

	sp.Suffix = fmt.Sprintf(" [%s]: Processing...", p.config.Server.Address)
	sp.Color("fgHiGreen")
	sp.Start()
	if err := s.Connect(p.privateKey); err != nil {
		log.Fatalf("Error: %s", err)
		return nil
	}
	sp.FinalMSG = fmt.Sprintf("%s [%s]: OK!\n", green("✔"), p.config.Server.Address)
	s.Run("uname -a")
	sp.Stop()

	return nil
}
