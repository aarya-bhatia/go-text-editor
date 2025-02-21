package controller

import (
	"go-editor/model"

	"github.com/gdamore/tcell/v2"
)

func handleKeyInNormalMode(event *tcell.EventKey, editor *model.Application) {
	switch event.Rune() {
	case '0':
		if editor.CurrentFile != nil {
			editor.CurrentFile.SetXCursor(0)
		}

	case 'a':
		if editor.CurrentFile != nil {
			editor.CurrentFile.GetCurrentLine().Cursor =
				min(1+editor.CurrentFile.GetCurrentLine().Cursor,
					editor.CurrentFile.GetCurrentLine().Size())
			editor.Mode = model.INSERT_MODE
		}

	case 'A':
		if editor.CurrentFile != nil {
			editor.CurrentFile.GetCurrentLine().Cursor = editor.CurrentFile.GetCurrentLine().Size()
			editor.Mode = model.INSERT_MODE
		}

	case 'I':
		if editor.CurrentFile != nil {
			editor.CurrentFile.GetCurrentLine().Cursor = 0 // TODO: should move to first non-whitespace char instead.
			editor.Mode = model.INSERT_MODE
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
		editor.UserCommand = ""
	case 'i':
		editor.Mode = model.INSERT_MODE
	case 'f':
		if editor.CurrentFile != nil {
			editor.Mode = model.NORMAL_MODE_ARG_PENDING
		}

	case 'G':
		if editor.CurrentFile != nil && editor.CurrentFile.CountLines() > 0 {
			editor.CurrentFile.CursorLine = editor.CurrentFile.CountLines() - 1
		}

	case 'g':
		if editor.CurrentFile != nil {
			editor.CurrentFile.CursorLine = 0
		}
	}
}
