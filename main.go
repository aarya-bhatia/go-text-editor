package main

import (
	"go-editor/internal"
	"log"
	"os"
)

func main() {
	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	// Redirect log output to the file
	log.SetOutput(logFile)

  filenames := os.Args[1:]

  // Run application
	internal.Start(filenames)
}
