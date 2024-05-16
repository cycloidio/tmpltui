# tmpltui 

This is a small TUI to pick a template file and render it, the original file and the rendered result are shown side by side.

The `data.json` file contains the data that is used to render the template.

You can change the template delimiters by specifying in `data.json` a `left_delimiter` and a `right_delimiter` string.


## Example

A simple example configuration is in `example_cycloid.yaml`, you can change the data in `data.json` and see how the rendering changes.


## Installation

Precompiled binaries for released versions are available in the [release](https://github.com/cycloidio/tmpltui/releases) page on Github.

Example:

```bash
wget https://github.com/cycloidio/tmpltui/releases/download/v0.2.0/tmpltui-v0.2.0-linux-amd64.tar.gz
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
./tmpltui
```
From the prompt, select your `example_cycloid.yaml` template file, press `enter`, and then `f` to display it in full screen.
