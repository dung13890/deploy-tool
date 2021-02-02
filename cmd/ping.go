package cmd

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/dung13890/deploy-tool/config"
	"github.com/dung13890/deploy-tool/remote"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

type ping struct {
	config     config.Configuration
	privateKey string
	log        bool
}

func PingInit() *cli.Command {
	return &cli.Command{
		Name:    "ping",
		Aliases: []string{"p"},
		Usage:   "Testing connection into servers",
		Flags: []cli.Flag{
			config.Load,
			config.Identity,
			config.EnableLog,
		},
		Action: func(ctx *cli.Context) error {
			p := &ping{}
			err := p.config.ReadFile(ctx.String("config"))
			if err != nil {
				log.Fatal(err)
			}
			p.log = ctx.Bool("log")
			p.privateKey = ctx.String("identity")
			p.exec()
			return nil
		},
	}
}

func (p *ping) exec() error {
	var r remote.Remote = &remote.Server{}
	defer r.Close()
	r.Load(
		p.config.Server.Address,
		p.config.Server.User,
		p.config.Server.Port,
		p.config.Server.Dir,
	)
	fmt.Println("Testing connection into servers:")
	green := color.New(color.FgHiGreen).SprintFunc()
	sp := spinner.New(spinner.CharSets[50], 100*time.Millisecond)

	sp.Suffix = fmt.Sprintf(" [%s]: Processing...", p.config.Server.Address)
	sp.Color("fgHiGreen")
	sp.FinalMSG = fmt.Sprintf("%s [%s]: OK!\n", green("✔"), p.config.Server.Address)
	sp.Start()
	if err := r.Connect(p.privateKey); err != nil {
		log.Fatalf("Error: %s", err)
		return nil
	}
	p.command(r)
	sp.Stop()

	return nil
}

func (p *ping) command(r remote.Remote) error {
	r.Run("uname -a")
	if p.log {
		wg := sync.WaitGroup{}
		// Copy over tasks's STDOUT.
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := io.Copy(os.Stdout, r.Stdout())
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}
		}()
		// Copy over tasks's STDERR.
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := io.Copy(os.Stderr, r.StdErr())
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}
		}()
		wg.Wait()
	}
	return nil
}
