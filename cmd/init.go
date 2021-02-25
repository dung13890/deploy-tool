package cmd

import (
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/dung13890/deploy-tool/config"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

type answers struct {
	Project string
	Ip      string
	Repo    string
}

func NewInit() *cli.Command {
	return &cli.Command{
		Name:    "init",
		Aliases: []string{"i"},
		Usage:   "Setup new project",
		Action: func(ctx *cli.Context) error {
			a := &answers{}
			if err := a.exec(); err != nil {
				log.Fatal(err)
			}
			return nil
		},
	}
}

func (a *answers) exec() error {
	// Create question setup yml file
	qs := []*survey.Question{
		{
			Name:     "project",
			Prompt:   &survey.Input{Message: "Please setup project name:"},
			Validate: survey.Required,
		},
		{
			Name: "ip",
			Prompt: &survey.Input{
				Message: "Please setup IP Remote:",
				Default: "127.0.0.1",
			},
		},
		{
			Name: "repo",
			Prompt: &survey.Input{
				Message: "Please setup url repository:",
				Default: "git@github.com:repo/example.git",
			},
		},
	}

	if err := survey.Ask(qs, a); err != nil {
		return err
	}

	if err := a.createFile(a.Project); err != nil {
		return err
	}

	success := color.New(color.FgHiGreen, color.Bold).PrintlnFunc()
	success(fmt.Sprintf("Successfully generate file %s.yml", a.Project))

	return nil
}

func (a *answers) createFile(filename string) (err error) {
	// Get Current path
	dir, err := os.Getwd()
	if err != nil {
		return
	}

	filename = fmt.Sprintf("%s.yml", filename)

	// Check exists file config
	if _, err := os.Stat(filepath.Join(dir, filename)); err == nil || os.IsExist(err) {
		return errors.New(fmt.Sprintf("Warning: File config %s is exists.", filename))
	}
	// Create File config yml
	f, err := os.Create(filepath.Join(dir, filename))
	defer f.Close()
	if err != nil {
		return
	}
	// Parse template to template yml
	t, err := template.New("index").Parse(config.SourceYaml)
	if err != nil {
		return
	}
	t.Execute(f, a)

	return nil
}
