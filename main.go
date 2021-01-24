package main

import (
	"github.com/dung13890/deploy-tool/cmd"
	"github.com/dung13890/deploy-tool/config"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	cli.AppHelpTemplate = config.AppHelpTemplate
	cli.CommandHelpTemplate = config.CommandHelpTemplate
	ping := cmd.PingInit()
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		ping,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
