package internal

import (
	"strconv"
)

func handleUserCommand(app *Application) {
	if userCommand == "q" || userCommand == "quit" || userCommand == "exit" {
		app.Quit()
		return
	}

	if userCommand == "next" {
		app.OpenNextFile()
		return
	}

	if userCommand == "prev" {
		app.OpenPrevFile()
		return
	}

	userNumber, err := strconv.Atoi(userCommand) // check if its a numeral
	if err == nil {
		app.GotoLine(userNumber)
		return
	}
}
