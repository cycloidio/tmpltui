package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

func loadTmplData() map[string]any {
	ret := map[string]any{}
	data, err := os.ReadFile("data.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, &ret)
	if err != nil {
		log.Fatal(err)
	}

	return ret
}

// loadFile reads the file and attempt to render it as a template, returns
// the original file, the rendered. In case of errors the returned value
// will contain the error text.
func loadFile(fpath string, tdata map[string]any) (string, string) {
	fdata, err := os.ReadFile(fpath)
	if err != nil {
		return fmt.Sprintf("error reading file %s: %+v\n", fpath, err), ""
	}
	t := template.New(fpath).Option("missingkey=zero").Funcs(sprig.FuncMap())

	setDelimiters(t, tdata)

	fstring := string(fdata)
	t, err = t.Parse(fstring)
	if err != nil {
		return fstring, fmt.Sprintf("error parsing template %s: %+v\n", fpath, err)
	}

	var buff bytes.Buffer
	err = t.Execute(&buff, tdata)
	if err != nil {
		return fstring, fmt.Sprintf("error rendering template %s: %+v\n", fpath, err)
	}

	return fstring, buff.String()
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
