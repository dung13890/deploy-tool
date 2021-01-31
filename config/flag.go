package config

import (
	"github.com/urfave/cli/v2"
)

var Load = &cli.StringFlag{
	Name:    "config",
	Aliases: []string{"c"},
	Value:   "config.yml",
	Usage:   "Load configuration from `FILE`",
}

var Identity = &cli.StringFlag{
	Name:    "identity",
	Aliases: []string{"i"},
	Value:   "~/.ssh/id_rsa",
	Usage:   "Identity (private key) for RSA or DSA authentication",
}

var EnableLog = &cli.BoolFlag{
	Name:    "log",
	Aliases: []string{"l"},
	Usage:   "Enable log detail",
}
