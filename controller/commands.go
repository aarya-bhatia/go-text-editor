package controller

import (
	"errors"
	"go-editor/model"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

func handleUserCommand(app *model.Application) {
	app.UserCommand = strings.TrimSpace(app.UserCommand)
	args := strings.Split(app.UserCommand, " ")

	switch args[0] {
	case "q", "quit", "exit":
		app.Quit()
	case "next":
		app.OpenNextFile()
	case "prev":
		app.OpenPrevFile()
	case "open", "edit":
		if len(args) != 2 {
			app.StatusLine = "Expected 2 arguments"
			break
		}
		app.OpenFile(args[1])
		if app.CurrentFile == nil {
			panic("file not open")
		}
		err := app.CurrentFile.ReadFile()
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			log.Println(err)
			app.CloseFile()
		}
	case "close":
		if app.CurrentFile == nil {
			break
		}
		app.CloseFile()
	case "closeall":
		app.CloseAll()
	case "ls":
		if len(app.Files) == 0 {
			log.Println("no files are open")
		}
		for _, file := range app.Files {
			log.Println(file.Name)
		}
	default:
		if len(args) == 1 {
			userNumber, err := strconv.Atoi(args[0]) // check if its a numeral
			if err == nil {
				log.Infof("Goto line: %d", userNumber)
				app.GotoLine(userNumber - 1) // make line number zero-indexed
			}
		}
	}
}
