package main

import (
	"log"
	"os"
)

func printRender() {
	data, err := loadTmplData()
	if err != nil {
		log.Fatalf("failed to load templating variables: '%v'", err)
	}
	_, rendered, err := loadFile(path, data)
	if err != nil {
		log.Fatal(err)
	}
	_, err = os.Stdout.WriteString(rendered)
	if err != nil {
		log.Fatal(err)
	}
}
