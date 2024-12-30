package main

import (
	"fmt"
	"log"

	"go-editor/internal"
	"go-editor/util"
)

func main() {
	fmt.Println("Hello World")
	util.SomeNumber()

	var app *internal.Application = internal.NewApplication()

	if err := app.OpenFile("hello.txt"); err != nil {
		log.Fatal(err)
	}
}
