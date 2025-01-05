package main

import (
	"go-editor/config"
	"go-editor/internal"
	"log"
	"strings"

	"github.com/gdamore/tcell/v2"
)

func refreshScreen(s tcell.Screen, editor *internal.Application) {
	s.Clear()

	displayLines := make([]string, 0)

	if editor.CurrentFile != nil {
		lines := editor.CurrentFile.Lines
		if editor.CurrentFile.ScrollY > 0 {
			if len(lines) > editor.CurrentFile.ScrollY {
				lines = lines[editor.CurrentFile.ScrollY:]
			} else {
				lines = make([]*internal.Line, 0)
			}
		}
		if len(lines) > config.MAX_DISPLAY_LINES {
			lines = lines[:config.MAX_DISPLAY_LINES]
		}

		for _, line := range lines {
			text := line.Text
			if editor.CurrentFile.ScrollX > 0 {
				if len(text) > editor.CurrentFile.ScrollX {
					text = text[editor.CurrentFile.ScrollX:]
				} else {
					text = ""
				}
			}

      blank_line := ""
      for i := 0; i < config.MAX_DISPLAY_COLS; i++ {
        blank_line += " "
      }

      text = text + blank_line // pad line with blank spaces
      text = text[:config.MAX_DISPLAY_COLS]

			displayLines = append(displayLines, text)
		}

		displayString := strings.Join(displayLines, "")

		DrawBox(s, config.EDITOR_BOX_LEFT, config.EDITOR_BOX_TOP, config.EDITOR_BOX_LEFT+config.EDITOR_BOX_WIDTH,
			config.EDITOR_BOX_TOP+config.EDITOR_BOX_HEIGHT, tcell.StyleDefault, displayString)

		DrawBox(s, config.STATUS_BOX_LEFT, config.STATUS_BOX_TOP, config.STATUS_BOX_LEFT+config.STATUS_BOX_WIDTH,
			config.STATUS_BOX_TOP+config.STATUS_BOX_HEIGHT, tcell.StyleDefault, editor.StatusLine)

		displayCursorX := editor.CurrentFile.GetCurrentLine().Cursor - editor.CurrentFile.ScrollX
		displayCursorY := editor.CurrentFile.CursorLine - editor.CurrentFile.ScrollY

		if displayCursorY < 0 {
			log.Print("WARN: cursor out of bounds")
			displayCursorY = 0
		}

		if displayCursorX < 0 {
			log.Print("WARN: cursor out of bounds")
			displayCursorX = 0
		}

		displayCursorX += config.EDITOR_BOX_LEFT + 1
		displayCursorY += config.EDITOR_BOX_TOP + 1

		s.ShowCursor(displayCursorX, displayCursorY)

	} else {
		s.ShowCursor(config.EDITOR_BOX_LEFT+1, config.EDITOR_BOX_HEIGHT+1)
	}

	s.Show()
}
