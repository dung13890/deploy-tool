package cmd

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/dung13890/deploy-tool/cmd/task"
	"github.com/dung13890/deploy-tool/config"
	"github.com/dung13890/deploy-tool/remote"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"log"
	"sync"
	"time"
)

type ping struct {
	config     config.Configuration
	privateKey string
	log        bool
}

func NewPing() *cli.Command {
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
	var r remote.Remote

	if p.config.Server.Address == "127.0.0.1" || p.config.Server.Address == "localhost" {
		r = &remote.Localhost{}
	} else {
		r = &remote.Server{}
	}
	defer r.Close()
	r.Load(
		p.config.Server.Address,
		p.config.Server.User,
		p.config.Server.Group,
		p.config.Server.Port,
		p.config.Server.Dir,
		p.config.Server.Project,
	)
	fmt.Println("Testing connection into servers:")
	green := color.New(color.FgHiGreen).SprintFunc()
	sp := spinner.New(spinner.CharSets[50], 100*time.Millisecond)

	sp.Suffix = fmt.Sprintf(" %s: Processing...", r.Prefix())
	sp.Color("fgHiGreen")
	sp.FinalMSG = fmt.Sprintf("%s %s: OK!\n", green("✔"), r.Prefix())
	sp.Start()
	if err := r.Connect(p.privateKey); err != nil {
		log.Fatalf("Error: %s", err)
	}
	t := task.NewTask(r, p.log)
	if err := p.command(t); err != nil {
		log.Fatalf("Error: %s", err)
	}
	sp.Stop()
	// Testing on cluster
	if len(p.config.Cluster.Hosts) > 0 {
		fmt.Println("Testing connection into cluster:")
		if err := p.cluster(t); err != nil {
			log.Fatalf("Error: %s", err)
		}
	}

	return nil
}

func (p *ping) command(t *task.Task) error {
	return t.Run("uname -a")
}

func (p *ping) cluster(t *task.Task) error {
	wg := sync.WaitGroup{}
	rs := make(chan string, len(p.config.Cluster.Hosts))
	er := make(chan string, len(p.config.Cluster.Hosts))
	green := color.New(color.FgHiGreen).SprintFunc()
	red := color.New(color.FgHiRed).SprintFunc()

	for _, item := range p.config.Cluster.Hosts {
		wg.Add(1)
		go func(w *sync.WaitGroup, host string) {
			defer w.Done()
			cmd := fmt.Sprintf("ssh %s 'uname -a'", host)
			if err := t.Run(cmd); err != nil {
				// Push channel when exists error
				er <- fmt.Sprintf("%s [%s]: Failed!", red("✘"), host)
			}
			// Push channel when connection success
			rs <- fmt.Sprintf("%s [%s]: OK!", green("✔"), host)
		}(&wg, item)
	}

	wg.Wait()
	for i := 0; i < len(p.config.Cluster.Hosts); i++ {
		select {
		case results := <-rs:
			fmt.Println(results)
		case errors := <-er:
			fmt.Println(errors)
		default:
			fmt.Println()
		}
	}
	close(rs)
	close(er)

	return nil
}
