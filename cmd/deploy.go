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

type deploy struct {
	mfuncs     map[string]interface{}
	config     config.Configuration
	privateKey string
	log        bool
	repo       *cmdDep.Repo
	shared     *cmdDep.Shared
	tasks      *cmdDep.Tasks
	cluster    *cmdDep.Cluster
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
			d.loadRepo(ctx.String("tag"), ctx.String("branch"))
			d.loadShared()
			d.loadTasks()
			d.loadCluster()
			d.privateKey = ctx.String("identity")
			d.exec()

			return nil
		},
	}
}

func (d *deploy) exec() error {
	// Init mfuncs
	d.mfuncs = map[string]interface{}{
		"deploy:prepare": cmdDep.Prepare,
		"deploy:fetch":   d.repo.Fetch,
		"deploy:shared":  d.shared.Run,
		"deploy:tasks":   d.tasks.Run,
		"deploy:cluster": d.cluster.Run,
		"deploy:publish": cmdDep.Publish,
	}

	var r remote.Remote

	if d.config.Server.Address == "127.0.0.1" || d.config.Server.Address == "localhost" {
		r = &remote.Localhost{}
	} else {
		r = &remote.Server{}
	}

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
	// Run Commands for deploy
	d.commands(t, "deploy:prepare")
	d.commands(t, "deploy:fetch")
	d.commands(t, "deploy:shared")
	d.commands(t, "deploy:tasks")
	d.commands(t, "deploy:publish")
	d.commands(t, "deploy:cluster")

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
	out, _ := utils.Call(d.mfuncs, cmds, t)
	if !out[0].IsNil() {
		log.Fatalf("Error: %v", out[0].Interface())
	}
	sp.Stop()

	return nil
}

func (d *deploy) loadRepo(tag string, branch string) *cmdDep.Repo {
	t := d.config.Repository.Tag
	if tag != "" {
		t = tag
	}

	b := d.config.Repository.Branch
	if branch != "" {
		b = branch
	}

	d.repo = cmdDep.NewRepo(d.config.Repository.Url, b, t)

	return d.repo
}

func (d *deploy) loadShared() *cmdDep.Shared {
	d.shared = cmdDep.NewShared(d.config.Shared.Folders, d.config.Shared.Files)

	return d.shared
}

func (d *deploy) loadTasks() *cmdDep.Tasks {
	d.tasks = cmdDep.NewTasks(d.config.Tasks)

	return d.tasks
}

func (d *deploy) loadCluster() *cmdDep.Cluster {
	d.cluster = cmdDep.NewCluster(
		d.config.Cluster.Hosts,
		d.config.Cluster.Rsync.Excludes,
		d.config.Cluster.Cmds,
	)

	return d.cluster
}
