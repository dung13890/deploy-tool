package main

import (
	"github.com/dung13890/deploy-tool/cmd"
	"github.com/dung13890/deploy-tool/config"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"time"
)

func main() {
	cli.AppHelpTemplate = config.AppHelpTemplate
	cli.CommandHelpTemplate = config.CommandHelpTemplate
	ping := cmd.NewPing()
	deploy := cmd.NewDeploy()
	init := cmd.NewInit()

	app := &cli.App{
		Name:                 "doo",
		Usage:                "Deployment for your project",
		HelpName:             "doo",
		Compiled:             time.Now(),
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			init,
			ping,
			deploy,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
