package controller

import (
	"go-editor/model"

	"github.com/gdamore/tcell/v2"
)

func handleKeyInInsertMode(event *tcell.EventKey, editor *model.Application) {
	if event.Key() == tcell.KeyEscape {
		editor.Mode = model.NORMAL_MODE
		if editor.CurrentFile != nil {
			if editor.CurrentFile.GetCurrentLine().Cursor > 0 {
				editor.CurrentFile.GetCurrentLine().Cursor -= 1
			}
		}
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
