package config

import "github.com/fatih/color"

// AppHelpTemplate is the text template for the Default help topic.
// cli.go uses text/template to render templates. You can
// render custom help text by setting this variable.

var (
	green  = color.New(color.FgHiGreen).SprintFunc()
	yellow = color.New(color.FgHiYellow).SprintFunc()
)

var AppHelpTemplate = yellow("NAME:") + green(`
	{{.Name}}{{if .Usage}} - {{.Usage}}{{end}}

`) + yellow("USAGE:") + green(`
	{{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Version}}{{if not .HideVersion}}

`) + yellow("VERSION:") + green(`
	{{.Version}}{{end}}{{end}}{{if .Description}}

`) + yellow("DESCRIPTION:") + green(`
	{{.Description | nindent 3 | trim}}{{end}}{{if len .Authors}}

`) + yellow("AUTHOR:") + green(`{{with $length := len .Authors}}{{if ne 1 $length}}S{{end}}{{end}}:
	{{range $index, $author := .Authors}}{{if $index}}
	{{end}}{{$author}}{{end}}{{end}}{{if .VisibleCommands}}

`) + yellow("COMMANDS:") + green(`{{range .VisibleCategories}}{{if .Name}}
	{{.Name}}:{{range .VisibleCommands}}
	  {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{else}}{{range .VisibleCommands}}
	{{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{end}}{{end}}{{end}}{{if .VisibleFlags}}

`) + yellow("GLOBAL OPTIONS:") + green(`
	{{range $index, $option := .VisibleFlags}}{{if $index}}
	{{end}}{{$option}}{{end}}{{end}}{{if .Copyright}}

`) + yellow("COPYRIGHT:") + green(`
	{{.Copyright}}{{end}}
`)

// CommandHelpTemplate is the text template for the command help topic.
// cli.go uses text/template to render templates. You can
// render custom help text by setting this variable.
var CommandHelpTemplate = yellow("NAME:") + green(`
	{{.HelpName}} - {{.Usage}}

`) + yellow("USAGE:") + green(`
	{{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}}{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Category}}

`) + yellow("CATEGORY:") + green(`
	{{.Category}}{{end}}{{if .Description}}

`) + yellow("DESCRIPTION:") + green(`
	{{.Description | nindent 3 | trim}}{{end}}{{if .VisibleFlags}}

`) + yellow("OPTIONS:") + green(`
	{{range .VisibleFlags}}{{.}}
	{{end}}{{end}}
`)
