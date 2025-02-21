package controller

import (
	"go-editor/model"

	"github.com/gdamore/tcell/v2"
)

func handleKeyInCommandMode(event *tcell.EventKey, editor *model.Application) {
	if event.Key() == tcell.KeyEnter || event.Key() == tcell.KeyEscape {
		editor.Mode = model.NORMAL_MODE
		handleUserCommand(editor)
	} else if event.Key() == tcell.KeyBS || event.Key() == tcell.KeyBackspace2 {
		if len(editor.UserCommand) > 0 {
			editor.UserCommand = editor.UserCommand[:len(editor.UserCommand)-1]
		} else {
			editor.Mode = model.NORMAL_MODE
		}
	} else if event.Rune() != 0 {
		editor.UserCommand += string(event.Rune())
	}
}
