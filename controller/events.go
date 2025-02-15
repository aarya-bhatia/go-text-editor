package controller

import (
	"go-editor/config"
	"go-editor/model"
	"log"

	"github.com/gdamore/tcell/v2"
)

func handleKeyEvent(event *tcell.EventKey, editor *model.Application, screen tcell.Screen) {
	// if config.DEBUG {
	// 	log.Println("Got key", event.Name())
	// }

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
	case model.NORMAL_MODE:
		handleKeyInNormalMode(event, editor)

	case model.COMMAND_MODE:
		handleKeyInCommandMode(event, editor)

	case model.INSERT_MODE:
		handleKeyInInsertMode(event, editor)

	case model.NORMAL_MODE_ARG_PENDING:
		handleKeyInNormalModeArgPending(event, editor)
	}

	if config.DEBUG {
		if editor.CurrentFile != nil {
			log.Printf("Cursor: %d,%d | Scroll: %d,%d",
				editor.CurrentFile.CursorLine,
				editor.CurrentFile.GetCurrentLine().Cursor,
				editor.CurrentFile.ScrollX,
				editor.CurrentFile.ScrollY)
		}
	}
}

func handleKeyInNormalModeArgPending(event *tcell.EventKey, editor *model.Application) {
	editor.Mode = model.NORMAL_MODE
	if event.Rune() != 0 {
		editor.CurrentFile.JumpToNextChar(event.Rune())
	}
}

func handleKeyInNormalMode(event *tcell.EventKey, editor *model.Application) {
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
		editor.Mode = model.COMMAND_MODE
		editor.StatusLine = ""
		editor.UserCommand = ""
	case 'i':
		editor.Mode = model.INSERT_MODE
		editor.StatusLine = ""
	case 'f':
		if editor.CurrentFile != nil {
			editor.Mode = model.NORMAL_MODE_ARG_PENDING
		}
	}
}

func handleKeyInCommandMode(event *tcell.EventKey, editor *model.Application) {
	if event.Key() == tcell.KeyEnter || event.Key() == tcell.KeyEscape {
		editor.Mode = model.NORMAL_MODE
		handleUserCommand(editor)
	} else if event.Key() == tcell.KeyBS || event.Key() == tcell.KeyBackspace2 {
		editor.UserCommand = editor.UserCommand[:len(editor.UserCommand)-1]
		editor.StatusLine = editor.UserCommand
	} else if event.Rune() != 0 {
		editor.UserCommand += string(event.Rune())
		editor.StatusLine = editor.UserCommand
	}
}

func handleKeyInInsertMode(event *tcell.EventKey, editor *model.Application) {
	if event.Key() == tcell.KeyEscape {
		editor.Mode = model.NORMAL_MODE
	} else if event.Key() == tcell.KeyBS || event.Key() == tcell.KeyBackspace2 {
		if editor.CurrentFile.GetCurrentLine().Size() == 0 {
			editor.CurrentFile.DeleteLine()
			editor.CurrentFile.GetCurrentLine().MoveToEnd()
		} else {
			editor.CurrentFile.GetCurrentLine().RemoveChar()
		}
	} else if event.Key() == tcell.KeyEnter {
		editor.CurrentFile.InsertLineBelowCursor()
		editor.CurrentFile.CursorLine += 1
	} else if event.Rune() != 0 {
		if editor.CurrentFile != nil {
			editor.CurrentFile.Insert(event.Rune())
		}
	}
}
