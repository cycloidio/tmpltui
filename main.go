package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	usage = `usage: tmpl [-h | -help] [-c | -config <path>] <command> [<args>]
	
options: 
	-h, -help displays this text
	-c, -config override the config for templating.
	-t, -type the type of templating, valid values are 'stacks' or 'blueprints'. default 'stacks'

commands:

render
	renders the template to stdout
	usage: tmpl render [-h | -help] [-f | -file <path>]
	options: 
		-h, -help displays this text
		-f, -file path to template file. default: .cycloid.yaml

tui
	starts a tui to view the template and the rendered file side by side
`

	renderUsage = `usage: tmpl render [-h | -help] [-f | -file <path>]
options: 
	-h, -help displays this text
	-f, -file path to template file. default: .cycloid.yaml
	`
)

var (
	tmplType, config, path string
)

func main() {

	flag.StringVar(&tmplType, "t", "stacks", "the type of templating, valid values are 'stacks' or 'blueprints'")
	flag.StringVar(&tmplType, "type", "stacks", "the type of templating, valid values are 'stacks' or 'blueprints'")
	flag.Usage = func() { fmt.Print(usage) }

	// global option
	flag.StringVar(&config, "c", "", "data and config used to render the template")
	flag.StringVar(&config, "config", "", "data and config used to render the template")
	flag.Usage = func() { fmt.Print(usage) }

	// render options
	renderFlags := flag.NewFlagSet("render", flag.ExitOnError)
	renderFlags.StringVar(&path, "f", ".cycloid.yaml", "path to template file")
	renderFlags.StringVar(&path, "file", ".cycloid.yaml", "path to template file")
	renderFlags.Usage = func() { fmt.Print(renderUsage) }

	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Print(usage)
		os.Exit(1)
	}

	cmd, args := args[0], args[1:]
	switch cmd {
	case "tui":
		tui()
	case "render":
		if err := renderFlags.Parse(args); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		printRender()
	default:
		fmt.Printf("unrecognized command %q, expected 'tui' or 'render' subcommands", cmd)
		os.Exit(1)
	}
}
