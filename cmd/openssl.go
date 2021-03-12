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
	"time"
)

type openssl struct {
	config     config.Configuration
	privateKey string
	log        bool
}

func NewOpenssl() *cli.Command {
	return &cli.Command{
		Name:    "openssl",
		Aliases: []string{"o"},
		Usage:   "Check the SSL/TLS Cipher Suites in Servers",
		Flags: []cli.Flag{
			config.Load,
			config.Identity,
			config.EnableLog,
		},
		Action: func(ctx *cli.Context) error {
			o := &openssl{}
			err := o.config.ReadFile(ctx.String("config"))
			if err != nil {
				log.Fatal(err)
			}
			o.log = ctx.Bool("log")
			o.privateKey = ctx.String("identity")
			o.exec()
			return nil
		},
	}
}

func (o *openssl) exec() error {
	var r remote.Remote

	if o.config.Server.Address == "127.0.0.1" || o.config.Server.Address == "localhost" {
		r = &remote.Localhost{}
	} else {
		r = &remote.Server{}
	}
	defer r.Close()
	r.Load(
		o.config.Server.Address,
		o.config.Server.User,
		o.config.Server.Group,
		o.config.Server.Port,
		o.config.Server.Dir,
		o.config.Server.Project,
	)
	fmt.Println("Testing connection into servers:")
	green := color.New(color.FgHiGreen).SprintFunc()
	sp := spinner.New(spinner.CharSets[50], 100*time.Millisecond)

	sp.Suffix = fmt.Sprintf(" %s: Processing...", r.Prefix())
	sp.Color("fgHiGreen")
	sp.FinalMSG = fmt.Sprintf("%s %s: OK!\n", green("✔"), r.Prefix())
	sp.Start()
	if err := r.Connect(o.privateKey); err != nil {
		log.Fatalf("Error: %s", err)
	}
	t := task.NewTask(r, o.log)
	if err := o.command(t); err != nil {
		log.Fatalf("Error: %s", err)
	}
	sp.Stop()

	return nil
}

func (o *openssl) command(t *task.Task) error {
	return t.Run("uname -a")
}
