package deploy

import (
	"github.com/urfave/cli/v2"
)

var Branch = &cli.StringFlag{
	Name:    "branch",
	Aliases: []string{"b"},
	Value:   "master",
	Usage:   "deploy with branch default `master`",
}

var Tag = &cli.StringFlag{
	Name:    "tag",
	Aliases: []string{"t"},
	Value:   "1.0.0",
	Usage:   "deploy with tag default `1.0.0`",
}
