package internal

import (
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
)

func handleUserCommand(app *Application) {
	userCommand = strings.TrimSpace(userCommand)
	args := strings.Split(userCommand, " ")

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
				app.GotoLine(userNumber)
			}
		}
	}
}
