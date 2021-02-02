package cmd

import (
	"fmt"
	"github.com/briandowns/spinner"
	cmdDep "github.com/dung13890/deploy-tool/cmd/deploy"
	"github.com/dung13890/deploy-tool/config"
	"github.com/dung13890/deploy-tool/remote"
	"github.com/dung13890/deploy-tool/utils"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"io"
	"log"
	"os"
	"sync"
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
	fmt.Printf("[%s] Executing task deploy:\n", d.config.Server.Address)
	if err := r.Connect(d.privateKey); err != nil {
		log.Fatalf("Error: %s", err)
		return nil
	}
	d.commands(r, "deploy:prepare")
	d.commands(r, "deploy:install")
	d.commands(r, "deploy:shared")
	d.commands(r, "deploy:vendors")
	d.commands(r, "deploy:migrate")
	d.commands(r, "deploy:release")
	success := color.New(color.FgHiGreen, color.Bold).PrintlnFunc()
	success("Successfully deployed!")

	return nil
}

func (d *deploy) commands(r remote.Remote, cmd string) error {
	green := color.New(color.FgHiGreen, color.Bold).SprintFunc()
	sp := spinner.New(spinner.CharSets[50], 100*time.Millisecond)
	sp.Suffix = fmt.Sprintf(" [%s]:	Processing...", cmd)
	sp.Color("fgHiGreen")
	sp.FinalMSG = fmt.Sprintf("%s [%s]:	Completed!\n", green("✔"), cmd)
	sp.Start()
	utils.Call(mfuncs, cmd, r)
	if d.log {
		d.printLog(r)
	}
	sp.Stop()

	return nil
}

func (d *deploy) printLog(r remote.Remote) error {
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

	return nil
}
