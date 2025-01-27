package internal

import (
	"fmt"
	"go-editor/config"
	"log"

	"github.com/gdamore/tcell/v2"
)

func NewScreen() (tcell.Screen, func()) {
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	s.SetStyle(defStyle)
	s.Clear()

  quit := func() {
		// You have to catch panics in a defer, clean up, and
		// re-raise them - otherwise your application can
		// die without leaving any diagnostic trace.
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}

	return s, quit
}

func getModeName(mode int) string {
	switch mode {
	case NORMAL_MODE:
		return "NORMAL"
	case COMMAND_MODE:
		return "COMMAND"
	case INSERT_MODE:
		return "INSERT"
	default:
		return "UNNAMED"
	}
}

func getVisibleText(file *File) [][]rune {
	displayLines := make([][]rune, 0)

	lines := file.Lines

	if file.ScrollY > 0 {
		if len(lines) > file.ScrollY {
			lines = lines[file.ScrollY:]
		} else {
			lines = make([]*Line, 0)
		}
	}
	if len(lines) > config.MAX_DISPLAY_LINES {
		lines = lines[:config.MAX_DISPLAY_LINES]
	}

	for _, line := range lines {
		text := line.Text
		if file.ScrollX > 0 {
			if len(text) > file.ScrollX {
				text = text[file.ScrollX:]
			} else {
				text = []rune{}
			}
		}

		blank_line := ""
		for i := 0; i < config.MAX_DISPLAY_COLS; i++ {
			blank_line += " "
		}

		text = append(text, []rune(blank_line)...) // pad line with blank spaces
		text = text[:config.MAX_DISPLAY_COLS]

		displayLines = append(displayLines, text)
	}

	return displayLines
}

func getDisplayCursor(file *File) (int, int) {
	displayCursorX := file.GetCurrentLine().Cursor - file.ScrollX
	displayCursorY := file.CursorLine - file.ScrollY

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

	return displayCursorX, displayCursorY
}

func displayFile(s tcell.Screen, file *File) {

	displayLines := getVisibleText(file)
	displayString := FlattenList(displayLines)

	DrawBox(s, config.EDITOR_BOX_LEFT, config.EDITOR_BOX_TOP, config.EDITOR_BOX_LEFT+config.EDITOR_BOX_WIDTH,
		config.EDITOR_BOX_TOP+config.EDITOR_BOX_HEIGHT, tcell.StyleDefault, displayString)
}

func refreshScreen(s tcell.Screen, editor *Application) {
	s.Clear()

	if editor.CurrentFile != nil {
		editor.CurrentFile.AdjustScroll()
		displayFile(s, editor.CurrentFile)

		statusLineWithModename := fmt.Sprintf("[%s] %s", getModeName(editor.Mode), editor.StatusLine)
		DrawBox(s, config.STATUS_BOX_LEFT, config.STATUS_BOX_TOP, config.STATUS_BOX_LEFT+config.STATUS_BOX_WIDTH,
			config.STATUS_BOX_TOP+config.STATUS_BOX_HEIGHT, tcell.StyleDefault, []rune(statusLineWithModename))

		cursorX, cursorY := getDisplayCursor(editor.CurrentFile)
		s.ShowCursor(cursorX, cursorY)

	} else {
		s.ShowCursor(config.EDITOR_BOX_LEFT+1, config.EDITOR_BOX_HEIGHT+1)
	}

	s.Show()
}
