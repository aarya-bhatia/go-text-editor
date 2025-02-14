package internal

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
)

type ViewModel struct {
	Screen tcell.Screen

	EditorBoxTop    int
	EditorBoxLeft   int
	EditorBoxWidth  int
	EditorBoxHeight int

	StatusBoxTop    int
	StatusBoxLeft   int
	StatusBoxWidth  int
	StatusBoxHeight int
}

func NewViewModel(screen tcell.Screen) *ViewModel {
	width, height := screen.Size()

	viewModel := new(ViewModel)
	viewModel.Screen = screen

	viewModel.EditorBoxTop = 0
	viewModel.EditorBoxLeft = 0
	viewModel.EditorBoxWidth = width - 1
	viewModel.EditorBoxHeight = height - 2

	viewModel.StatusBoxTop = viewModel.EditorBoxHeight + 1
	viewModel.StatusBoxLeft = 0
	viewModel.StatusBoxWidth = width
	viewModel.StatusBoxHeight = 4

	return viewModel
}

func (view *ViewModel) GetMaxDisplayLines() int {
	return view.EditorBoxHeight - 2
}

func (view *ViewModel) GetMaxDisplayCols() int {
	return view.EditorBoxWidth - 2
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

func (view *ViewModel) getVisibleText(file *File) [][]rune {
	displayLines := make([][]rune, 0)

	lines := file.Lines

	if file.ScrollY > 0 {
		if len(lines) > file.ScrollY {
			lines = lines[file.ScrollY:]
		} else {
			lines = make([]*Line, 0)
		}
	}
	if len(lines) > view.GetMaxDisplayLines() {
		lines = lines[:view.GetMaxDisplayLines()]
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
		for i := 0; i < view.GetMaxDisplayCols(); i++ {
			blank_line += " "
		}

		text = append(text, []rune(blank_line)...) // pad line with blank spaces
		text = text[:view.GetMaxDisplayCols()]

		displayLines = append(displayLines, text)
	}

	return displayLines
}

func (view *ViewModel) getDisplayCursor(file *File) (int, int) {
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

	displayCursorX += view.EditorBoxLeft + 1
	displayCursorY += view.EditorBoxTop + 1

	return displayCursorX, displayCursorY
}

func (view *ViewModel) displayFile(file *File) {

	displayLines := view.getVisibleText(file)
	displayString := FlattenList(displayLines)
	view.renderEditorBox(displayString)
}

func (view *ViewModel) renderEditorBox(text []rune) {
	DrawBox(view.Screen, view.EditorBoxLeft,
		view.EditorBoxTop, view.EditorBoxLeft+view.EditorBoxWidth,
		view.EditorBoxTop+view.EditorBoxHeight, tcell.StyleDefault, text)
}

func (view *ViewModel) renderStatus(editor *Application) {
	filename := "No File"
	if editor.CurrentFile != nil {
		filename = editor.CurrentFile.Name
		if editor.CurrentFile.Modified {
			filename += " [+]"
		}
	}

	status := fmt.Sprintf("[%s] | %s | Ln %d, Col %d", getModeName(editor.Mode), filename,
		editor.CurrentFile.CursorLine, editor.CurrentFile.GetCurrentLine().Cursor)

	line := []rune(status)

	DrawText(view.Screen, view.StatusBoxLeft, view.StatusBoxTop,
		view.StatusBoxLeft+view.StatusBoxWidth,
		view.StatusBoxTop+view.StatusBoxHeight, tcell.StyleDefault, line)
}

func (view *ViewModel) renderCursor(x int, y int) {
	view.Screen.ShowCursor(x, y)
}

func (view *ViewModel) render(editor *Application) {
	view.Screen.Clear()

	if editor.CurrentFile != nil {
		editor.CurrentFile.AdjustScroll(view)
		view.displayFile(editor.CurrentFile)

		view.renderStatus(editor)

		cursorX, cursorY := view.getDisplayCursor(editor.CurrentFile)
		view.renderCursor(cursorX, cursorY)

	} else {
		view.renderEditorBox([]rune{})
		view.renderStatus(editor)
		view.renderCursor(view.EditorBoxLeft+1, view.EditorBoxTop+1)
	}

	view.Screen.Show()
	if editor.StatusLine != "" {
		// TODO: view.Screen.ShowMessage(editor.StatusLine)
		log.Println("NOTICE:", editor.StatusLine)
	}
}
