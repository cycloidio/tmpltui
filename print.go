package main

import (
	"log"
	"os"
)

func printRender() {
	data := loadTmplData()
	_, rendered, err := loadFile(path, data)
	if err != nil {
		log.Fatal(err)
	}
	_, err = os.Stdout.WriteString(rendered)
	if err != nil {
		log.Fatal(err)
	}
}
