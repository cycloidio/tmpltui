package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

var (
	//go:embed stack_templating.json
	stackTemplating embed.FS

	//go:embed blueprints.json
	blueprints embed.FS
)

func loadTmplData() (map[string]any, error) {
	values := map[string]any{}

	if config != "" {
		data, err := os.ReadFile(config)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(data, &values)
		if err != nil {
			return nil, err
		}

		return values, nil
	}

	var data []byte
	var err error
	switch tmplType {
	case "stacks":
		data, err = stackTemplating.ReadFile("stack_templating.json")
		if err != nil {
			return nil, err
		}
	case "blueprints":
		data, err = blueprints.ReadFile("blueprints.json")
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New(fmt.Sprintf("invalid templating type '%s', see usage for allowed values.", tmplType))
	}

	err = json.Unmarshal(data, &values)
	if err != nil {
		return nil, err
	}

	return values, nil
}

// loadFile reads the file and attempt to render it as a template, returns
// the original file, the rendered. In case of errors the returned value
// will contain the error text.
func loadFile(fpath string, tdata map[string]any) (string, string, error) {
	fdata, err := os.ReadFile(fpath)
	if err != nil {
		return fmt.Sprintf("error reading file %s: %+v\n", fpath, err), "", nil
	}
	t := template.New(fpath).Option("missingkey=zero").Funcs(sprig.FuncMap())

	setDelimiters(t, tdata)

	fstring := string(fdata)
	t, err = t.Parse(fstring)
	if err != nil {
		return fstring, fmt.Sprintf("error parsing template %s: %+v\n", fpath, err), err
	}

	var buff bytes.Buffer
	err = t.Execute(&buff, tdata)
	if err != nil {
		return fstring, fmt.Sprintf("error rendering template %s: %+v\n", fpath, err), err
	}

	return fstring, buff.String(), nil
}

const (
	leftDelim  = "left_delimiter"
	rightDelim = "right_delimiter"
)

// finds and sets the delimiters
func setDelimiters(t *template.Template, tdata map[string]any) {
	// default delimiters, left and right
	var ld, rd = "{{", "}}"
	if left, ok := tdata[leftDelim]; ok {
		ld = left.(string)
	}
	if right, ok := tdata[rightDelim]; ok {
		rd = right.(string)
	}
	t.Delims(ld, rd)
}
