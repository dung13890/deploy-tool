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
	init := cmd.NewInit()
	openssl := cmd.NewOpenssl()

	app := &cli.App{
		Name:                 "euv3",
		Usage:                "Check version the SSL/TLS On Server",
		HelpName:             "euv3",
		Compiled:             time.Now(),
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			init,
			openssl,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
