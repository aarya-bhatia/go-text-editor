package main

import (
	"log"

	"go-editor/internal"
)

func main() {
	var app *internal.Application = internal.NewApplication()

	if err := app.OpenFile("hello.txt"); err != nil {
		log.Fatal(err)
	}
}
