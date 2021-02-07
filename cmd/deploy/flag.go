package deploy

import (
	"github.com/urfave/cli/v2"
)

var Branch = &cli.StringFlag{
	Name:    "branch",
	Aliases: []string{"b"},
	Usage:   "deploy with branch `master`",
}

var Tag = &cli.StringFlag{
	Name:    "tag",
	Aliases: []string{"t"},
	Usage:   "deploy with tag `1.0.0`",
}
