package cmd

import (
	"fmt"
	"github.com/briandowns/spinner"
	cmdDep "github.com/dung13890/deploy-tool/cmd/deploy"
	"github.com/dung13890/deploy-tool/cmd/task"
	"github.com/dung13890/deploy-tool/config"
	"github.com/dung13890/deploy-tool/remote"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"log"
	"strings"
	"time"
)

type git struct {
	config     config.Configuration
	privateKey string
	log        bool
	branch     string
}

func NewGit() *cli.Command {
	return &cli.Command{
		Name:    "git",
		Aliases: []string{"g"},
		Usage:   "Use git pull into servers",
		Flags: []cli.Flag{
			config.Load,
			config.Identity,
			config.EnableLog,
			cmdDep.Branch,
		},
		Action: func(ctx *cli.Context) error {
			g := &git{}
			err := g.config.ReadFile(ctx.String("config"))
			if err != nil {
				log.Fatal(err)
			}
			g.log = ctx.Bool("log")
			g.branch = ctx.String("branch")
			g.privateKey = ctx.String("identity")
			g.exec()

			return nil
		},
	}
}

func (g *git) exec() error {
	var r remote.Remote

	if g.config.Server.Address == "127.0.0.1" || g.config.Server.Address == "localhost" {
		r = &remote.Localhost{}
	} else {
		r = &remote.Server{}
	}

	defer r.Close()
	r.Load(
		g.config.Server.Address,
		g.config.Server.User,
		g.config.Server.Group,
		g.config.Server.Port,
		g.config.Server.Dir,
		g.config.Server.Project,
	)
	fmt.Println("Use git pull into servers:")
	green := color.New(color.FgHiGreen).SprintFunc()
	sp := spinner.New(spinner.CharSets[50], 100*time.Millisecond)

	sp.Suffix = fmt.Sprintf(" %s: Processing...", r.Prefix())
	sp.Color("fgHiGreen")
	sp.FinalMSG = fmt.Sprintf("%s %s: Completed!\n", green("✔"), r.Prefix())
	sp.Start()
	if err := r.Connect(g.privateKey); err != nil {
		log.Fatalf("Error: %s", err)
	}
	t := task.NewTask(r, g.log)
	if err := g.command(t, "git_pull"); err != nil {
		log.Fatalf("Error: %s", err)
	}
	if err := g.command(t, "tasks"); err != nil {
		log.Fatalf("Error: %s", err)
	}
	sp.Stop()
	fmt.Println(green("Successfully deployed!"))

	return nil
}

func (g *git) command(t *task.Task, cmds string) error {
	path := t.GetDirectory()
	if cmds == "git_pull" {
		cmd := fmt.Sprintf("cd %s && git fetch origin %s", path, g.branch)
		if err := t.Run(cmd); err != nil {
			return err
		}
		cmd = fmt.Sprintf("cd %s && git checkout %s", path, g.branch)
		if err := t.Run(cmd); err != nil {
			return err
		}
		cmd = fmt.Sprintf("cd %s && git reset --hard origin/%s", path, g.branch)
		if err := t.Run(cmd); err != nil {
			return err
		}
	}

	if cmds == "tasks" {
		// Loop list unique tasks
		for _, v := range g.config.Tasks {
			v = strings.TrimSpace(v)
			cmd := fmt.Sprintf("cd %s && %s", path, v)
			if err := t.Run(cmd); err != nil {
				return err
			}
		}
	}
	return nil
}
