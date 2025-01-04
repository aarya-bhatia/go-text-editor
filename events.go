package main

import (
	"go-editor/internal"
	"log"
)

func handleKey(key rune, editor *internal.Application) {
	log.Println("Got key", key)

	switch key {
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
	}

	if editor.CurrentFile != nil {
		log.Printf("Cursor: %d,%d | Scroll: %d,%d",
			editor.CurrentFile.CursorLine,
			editor.CurrentFile.GetCurrentLine().Cursor,
			editor.CurrentFile.ScrollX,
			editor.CurrentFile.ScrollY)
	}
}
