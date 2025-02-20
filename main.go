package main

import (
	"go-editor/controller"
	"log"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	// Redirect log output to the file
	log.SetOutput(logFile)

	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(logFile)

	filenames := os.Args[1:]

	// Run application
	controller.Start(filenames)
}
