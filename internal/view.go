package internal

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
)

type ViewModel struct {
	Screen tcell.Screen
	Width  int
	Height int
}

func NewViewModel(screen tcell.Screen) *ViewModel {
	viewModel := new(ViewModel)
	viewModel.Screen = screen
	viewModel.Width, viewModel.Height = screen.Size()
	return viewModel
}

func (this *ViewModel) GetMaxEditorLines() int {
	return this.Height - 3
}

func (this *ViewModel) GetMaxEditorCols() int {
	return this.Width - 2
}

func (this *ViewModel) drawBorder() {

	for col := 1; col < this.Width-1; col++ {
		this.Screen.SetContent(col, 0, tcell.RuneHLine, nil, tcell.StyleDefault)
		this.Screen.SetContent(col, this.Height-2, tcell.RuneHLine, nil, tcell.StyleDefault)
	}
	for row := 1; row < this.Height-2; row++ {
		this.Screen.SetContent(0, row, tcell.RuneVLine, nil, tcell.StyleDefault)
		this.Screen.SetContent(this.Width-1, row, tcell.RuneVLine, nil, tcell.StyleDefault)
	}

	this.Screen.SetContent(0, 0, tcell.RuneULCorner, nil, tcell.StyleDefault)
	this.Screen.SetContent(this.Width-1, 0, tcell.RuneURCorner, nil, tcell.StyleDefault)
	this.Screen.SetContent(0, this.Height-2, tcell.RuneLLCorner, nil, tcell.StyleDefault)
	this.Screen.SetContent(this.Width-1, this.Height-2, tcell.RuneLRCorner, nil, tcell.StyleDefault)
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

func (this *ViewModel) getVisibleText(file *File) [][]rune {
	displayLines := make([][]rune, 0)

	lines := file.Lines

	if file.ScrollY > 0 {
		if len(lines) > file.ScrollY {
			lines = lines[file.ScrollY:]
		} else {
			lines = make([]*Line, 0)
		}
	}
	if len(lines) > this.GetMaxEditorLines() {
		lines = lines[:this.GetMaxEditorLines()]
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
		for i := 0; i < this.GetMaxEditorCols(); i++ {
			blank_line += " "
		}

		text = append(text, []rune(blank_line)...) // pad line with blank spaces
		text = text[:this.GetMaxEditorCols()]

		displayLines = append(displayLines, text)
	}

	return displayLines
}

func (this *ViewModel) getDisplayCursor(file *File) (int, int) {
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

	displayCursorX += 1
	displayCursorY += 2

	return displayCursorX, displayCursorY
}

func (this *ViewModel) displayFile(file *File) {

	displayLines := this.getVisibleText(file)
	displayString := FlattenList(displayLines)
	this.renderEditorBox(displayString)
}

func (this *ViewModel) renderEditorBox(text []rune) {
	DrawBox(this.Screen, 0, 0, this.Width, this.Height-2, tcell.StyleDefault, text)
}

func getStatusLine(editor *Application) string {
	return fmt.Sprintf("[%s] %s", getModeName(editor.Mode), editor.StatusLine)
}

func (this *ViewModel) renderCursor(x int, y int) {
	this.Screen.ShowCursor(x, y)
}

func (this *ViewModel) render(editor *Application) {
	this.Screen.Clear()

	this.drawBorder()

	statusline := []rune{}
	for i := 0; i < this.Width; i++ {
		statusline = append(statusline, ' ')
	}

	statusvalue := []rune(getStatusLine(editor))
	for i := 0; i < min(len(statusline)-1, len(statusvalue)); i++ {
		statusline[i+1] = statusvalue[i]
	}

	DrawText(this.Screen, 0, this.Height-1, this.Width, this.Height-1,
		tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite),
		statusline)

	if editor.CurrentFile != nil {
		editor.CurrentFile.AdjustScroll(this.Height, this.Width)
		this.displayFile(editor.CurrentFile)

		cursorX, cursorY := this.getDisplayCursor(editor.CurrentFile)
		this.renderCursor(cursorX, cursorY)

	} else {
		this.renderEditorBox([]rune{})
		this.renderCursor(1, 1)
	}

	this.Screen.Show()
}
