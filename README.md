# tmpl


This is a small CLI to explore and render tmaplated files. 
You can either start the TUI with the `tui` subcommand or print to stdout the rendered file with the `render` subcommand.

The TUI allows you to pick a template file and render it, the original file and the rendered result are shown side by side.

The `data.json` file contains the data that is used to render the template, you can change the filename and path with the config flag.

You can change the template delimiters by specifying in `data.json` a `left_delimiter` and a `right_delimiter` string.


## Example

A simple example configuration is in `example_cycloid.yaml`, you can change the data in `data.json` and see how the rendering changes.


## Installation

Precompiled binaries for released versions are available in the [release](https://github.com/cycloidio/tmpltui/releases) page on Github.

Example:

```bash
LATEST_RELEASE=$(curl -s https://api.github.com/repos/cycloidio/tmpltui/releases/latest | jq -r '.tag_name')
wget https://github.com/cycloidio/tmpltui/releases/download/$LATEST_RELEASE/tmpltui-$LATEST_RELEASE-linux-amd64.tar.gz
tar xf tmpltui-*.tar.gz
```

## Getting started

Get the `data.json` file that contains default Cycloid special variables.

```bash
wget https://raw.githubusercontent.com/cycloidio/tmpltui/main/data.json
```

Get an example of a template.

```bash
wget https://raw.githubusercontent.com/cycloidio/tmpltui/main/example_cycloid.yaml
```

```bash
./tmpl tui
```
From the prompt, select your `example_cycloid.yaml` template file, press `enter`, and then `f` to display it in full screen.

To print the rendered file on stdout:
```bash
tmpl render -f example_cycloid.yaml
```

## Usage 

```
usage: tmpl [-h | -help] [-c | -config <path>] <command> [<args>]
	
options: 
	-h, -help displays this text
	-c, -config data and config used to render the template
		default: data.json

commands:

render
	renders the template to stdout
	usage: tmpl render [-h | -help] [-f | -file <path>]
	options: 
		-h, -help displays this text
		-f, -file path to template file
			default: .cycloid.yaml

tui
	starts a tui to view the template and the rendered file side by side
```

