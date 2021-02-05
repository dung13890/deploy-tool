package cmd

import (
	"fmt"
	"github.com/briandowns/spinner"
	cmdDep "github.com/dung13890/deploy-tool/cmd/deploy"
	"github.com/dung13890/deploy-tool/cmd/task"
	"github.com/dung13890/deploy-tool/config"
	"github.com/dung13890/deploy-tool/remote"
	"github.com/dung13890/deploy-tool/utils"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"log"
	"time"
)

var mfuncs = map[string]interface{}{
	"deploy:prepare": cmdDep.Prepare,
	"deploy:install": cmdDep.Prepare,
	"deploy:shared":  cmdDep.Prepare,
	"deploy:vendors": cmdDep.Prepare,
	"deploy:migrate": cmdDep.Prepare,
	"deploy:release": cmdDep.Prepare,
}

type deploy struct {
	config     config.Configuration
	privateKey string
	log        bool
}

func DeployInit() *cli.Command {
	return &cli.Command{
		Name:    "deploy",
		Aliases: []string{"d"},
		Usage:   "Deploy into servers",
		Flags: []cli.Flag{
			config.Load,
			config.Identity,
			config.EnableLog,
			cmdDep.Tag,
			cmdDep.Branch,
		},
		Action: func(ctx *cli.Context) error {
			d := &deploy{}
			err := d.config.ReadFile(ctx.String("config"))
			if err != nil {
				log.Fatal(err)
			}
			d.log = ctx.Bool("log")
			d.privateKey = ctx.String("identity")
			d.exec()
			return nil
		},
	}
}

func (d *deploy) exec() error {
	var r remote.Remote = &remote.Server{}
	defer r.Close()
	r.Load(
		d.config.Server.Address,
		d.config.Server.User,
		d.config.Server.Port,
		d.config.Server.Dir,
	)
	fmt.Printf("%s Executing task deploy:\n", r.Prefix())
	if err := r.Connect(d.privateKey); err != nil {
		log.Fatalf("Error: %s", err)
	}
	t := task.New(r, d.log)
	d.commands(t, "deploy:prepare")
	d.commands(t, "deploy:install")
	d.commands(t, "deploy:shared")
	d.commands(t, "deploy:vendors")
	d.commands(t, "deploy:migrate")
	d.commands(t, "deploy:release")
	success := color.New(color.FgHiGreen, color.Bold).PrintlnFunc()
	success("Successfully deployed!")

	return nil
}

func (d *deploy) commands(t *task.Task, cmds string) error {
	green := color.New(color.FgHiGreen, color.Bold).SprintFunc()
	sp := spinner.New(spinner.CharSets[50], 100*time.Millisecond)
	sp.Suffix = fmt.Sprintf(" [%s]:	Processing...", cmds)
	sp.Color("fgHiGreen")
	sp.FinalMSG = fmt.Sprintf("%s [%s]:	Completed!\n", green("✔"), cmds)
	sp.Start()
	_, err := utils.Call(mfuncs, cmds, t)
	if err != nil {
		log.Fatal(err)
	}
	sp.Stop()

	return nil
}
