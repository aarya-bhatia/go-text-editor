package main

import (
	"go-editor/internal"
	"log"
	"strconv"

	"github.com/gdamore/tcell/v2"
)

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

	switch editor.Mode {
	case internal.NORMAL_MODE:
		handleKeyInNormalMode(event, editor)

	case internal.COMMAND_MODE:
		handleKeyInCommandMode(event, editor)

	case internal.INSERT_MODE:
		handleKeyInInsertMode(event, editor)

	case internal.NORMAL_MODE_ARG_PENDING:
		handleKeyInNormalModeArgPending(event, editor)
	}

	if editor.CurrentFile != nil {
		log.Printf("Cursor: %d,%d | Scroll: %d,%d",
			editor.CurrentFile.CursorLine,
			editor.CurrentFile.GetCurrentLine().Cursor,
			editor.CurrentFile.ScrollX,
			editor.CurrentFile.ScrollY)
	}
}

func handleKeyInNormalModeArgPending(event *tcell.EventKey, editor *internal.Application) {
	editor.Mode = internal.NORMAL_MODE
	if event.Rune() != 0 {
		editor.CurrentFile.JumpToNextChar(event.Rune())
	}
}

func handleKeyInNormalMode(event *tcell.EventKey, editor *internal.Application) {
	switch event.Rune() {
	case '0':
		if editor.CurrentFile != nil {
			editor.CurrentFile.SetXCursor(0)
		}

	case '$':
		if editor.CurrentFile != nil {
			editor.CurrentFile.GetCurrentLine().MoveToEnd()
		}
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
		editor.Mode = internal.COMMAND_MODE
		editor.StatusLine = ""
		userCommand = ""
	case 'i':
		editor.Mode = internal.INSERT_MODE
		editor.StatusLine = ""
	case 'f':
		if editor.CurrentFile != nil {
			editor.Mode = internal.NORMAL_MODE_ARG_PENDING
		}
	}
}

func handleKeyInCommandMode(event *tcell.EventKey, editor *internal.Application) {
	if event.Key() == tcell.KeyEnter || event.Key() == tcell.KeyEscape {
		editor.Mode = internal.NORMAL_MODE
		handleUserCommand(editor)
	} else if event.Key() == tcell.KeyBS || event.Key() == tcell.KeyBackspace2 {
    userCommand = userCommand[:len(userCommand)-1]
		editor.StatusLine = userCommand
	} else if event.Rune() != 0 {
		userCommand += string(event.Rune())
		editor.StatusLine = userCommand
	}
}

func handleKeyInInsertMode(event *tcell.EventKey, editor *internal.Application) {
	if event.Key() == tcell.KeyEnter || event.Key() == tcell.KeyEscape {
		editor.Mode = internal.NORMAL_MODE
		handleUserCommand(editor)
	} else if event.Rune() != 0 {
		if editor.CurrentFile != nil {
			editor.CurrentFile.InsertChar(byte(event.Rune())) // TODO: add support for runes
		}
	}
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
			editor.CurrentFile.SetYCursor(userNumber)
			return
		}
	}
}
