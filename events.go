package main

import (
	"go-editor/internal"
	"log"
	"strconv"

	"github.com/gdamore/tcell/v2"
)

var currentMode int = internal.NORMAL_MODE
var userCommand string = ""

func handleKeyEvent(event *tcell.EventKey, editor *internal.Application, screen tcell.Screen) {
	log.Println("Got key", event.Name())

	if event.Key() == tcell.KeyCtrlC {
		log.Println("Setting quit signal")
		editor.QuitSignal = true
		return
	}

	if event.Key() == tcell.KeyCtrlL {
		screen.Sync()
		return
	}

	switch currentMode {
	case internal.NORMAL_MODE:
		handleKeyInNormalMode(event, editor)

	case internal.COMMAND_MODE:
		handleKeyInCommandMode(event, editor)

	case internal.INSERT_MODE:
		handleKeyInInsertMode(event, editor)
	}

	if editor.CurrentFile != nil {
		log.Printf("Cursor: %d,%d | Scroll: %d,%d",
			editor.CurrentFile.CursorLine,
			editor.CurrentFile.GetCurrentLine().Cursor,
			editor.CurrentFile.ScrollX,
			editor.CurrentFile.ScrollY)
	}
}

func handleKeyInNormalMode(event *tcell.EventKey, editor *internal.Application) {
	switch event.Rune() {
	case 'h':
		if editor.CurrentFile != nil {
			editor.CurrentFile.MoveBackward()
		}
	case 'l':
		if editor.CurrentFile != nil {
			editor.CurrentFile.MoveForward()
		}
	case 'j':
		if editor.CurrentFile != nil {
			editor.CurrentFile.MoveDown()
		}
	case 'k':
		if editor.CurrentFile != nil {
			editor.CurrentFile.MoveUp()
		}
	case ':':
		currentMode = internal.COMMAND_MODE
		editor.StatusLine = ":"
		userCommand = ""
	}
}

func handleKeyInCommandMode(event *tcell.EventKey, editor *internal.Application) {
	if event.Key() == tcell.KeyEnter || event.Key() == tcell.KeyEscape {
		log.Println("Enter pressed")
		currentMode = internal.NORMAL_MODE
		editor.StatusLine = ""
		handleUserCommand(editor)

	} else if event.Rune() != 0 {
		userCommand += string(event.Rune())
		editor.StatusLine = ":" + userCommand
	}
}

func handleKeyInInsertMode(event *tcell.EventKey, editor *internal.Application) {
}

func handleUserCommand(editor *internal.Application) {
	if userCommand == "q" || userCommand == "quit" || userCommand == "exit" {
		log.Println("Setting quit signal")
		editor.QuitSignal = true
		return
	}

	userNumber, err := strconv.Atoi(userCommand) // check if its a numeral
	if err == nil {
		if editor.CurrentFile != nil {
			log.Println("Move to line ", userNumber)
			editor.CurrentFile.MoveToLineNo(userNumber)
			return
		}
	}
}
