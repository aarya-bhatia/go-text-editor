package controller

import (
	"go-editor/config"
	"go-editor/model"
	"log"

	"github.com/gdamore/tcell/v2"
)

func handleKeyEvent(event *tcell.EventKey, editor *model.Application, screen tcell.Screen) {
	if config.DEBUG {
		log.Println("Got key", event.Name())
	}

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
