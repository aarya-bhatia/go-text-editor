package main

import "go-editor/internal"

func handleKey(key rune, editor *internal.Application) {
	switch key {
	case 'h':
		if editor.CurrentFile != nil {
			editor.CurrentFile.GetCurrentLine().MoveBackward()
		}
	case 'l':
		if editor.CurrentFile != nil {
			editor.CurrentFile.GetCurrentLine().MoveForward()
		}
	case 'j':
		if editor.CurrentFile != nil {
			editor.CurrentFile.MoveDown()
		}
	case 'k':
		if editor.CurrentFile != nil {
			editor.CurrentFile.MoveUp()
		}
	}
}
