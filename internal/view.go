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

	const gapX = 2
	const gapY = 2

	viewModel := new(ViewModel)
	viewModel.Screen = screen
	viewModel.StatusBoxHeight = 2
	viewModel.EditorBoxTop = gapY
	viewModel.EditorBoxLeft = gapX
	viewModel.EditorBoxWidth = width - 2*gapX
	viewModel.EditorBoxHeight = height - 3*gapY - viewModel.StatusBoxHeight
	viewModel.StatusBoxTop = viewModel.EditorBoxHeight + 2*gapY
	viewModel.StatusBoxLeft = gapX
	viewModel.StatusBoxWidth = width - 2*gapX

	return viewModel
}

func (this *ViewModel) GetMaxDisplayLines() int {
	return this.EditorBoxHeight - 2
}

func (this *ViewModel) GetMaxDisplayCols() int {
	return this.EditorBoxWidth - 2
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
	if len(lines) > this.GetMaxDisplayLines() {
		lines = lines[:this.GetMaxDisplayLines()]
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
		for i := 0; i < this.GetMaxDisplayCols(); i++ {
			blank_line += " "
		}

		text = append(text, []rune(blank_line)...) // pad line with blank spaces
		text = text[:this.GetMaxDisplayCols()]

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

	displayCursorX += this.EditorBoxLeft + 1
	displayCursorY += this.EditorBoxTop + 1

	return displayCursorX, displayCursorY
}

func (this *ViewModel) displayFile(file *File) {

	displayLines := this.getVisibleText(file)
	displayString := FlattenList(displayLines)
	this.renderEditorBox(displayString)
}

func (this *ViewModel) renderEditorBox(text []rune) {
	DrawBox(this.Screen, this.EditorBoxLeft,
		this.EditorBoxTop, this.EditorBoxLeft+this.EditorBoxWidth,
		this.EditorBoxTop+this.EditorBoxHeight, tcell.StyleDefault, text)
}

func (this *ViewModel) renderStatusBox(text []rune) {
	DrawBox(this.Screen, this.StatusBoxLeft,
		this.StatusBoxTop, this.StatusBoxLeft+this.StatusBoxWidth,
		this.StatusBoxTop+this.StatusBoxHeight, tcell.StyleDefault, text)
}

func getStatusLine(editor *Application) string {
	return fmt.Sprintf("[%s] %s", getModeName(editor.Mode), editor.StatusLine)
}

func (this *ViewModel) renderCursor(x int, y int) {
	this.Screen.ShowCursor(x, y)
}

func (this *ViewModel) render(editor *Application) {
	this.Screen.Clear()

	if editor.CurrentFile != nil {
		editor.CurrentFile.AdjustScroll(this)
		this.displayFile(editor.CurrentFile)

		this.renderStatusBox([]rune(getStatusLine(editor)))

		cursorX, cursorY := this.getDisplayCursor(editor.CurrentFile)
		this.renderCursor(cursorX, cursorY)

	} else {
		this.renderEditorBox([]rune{})
		this.renderStatusBox([]rune{})
		this.renderStatusBox([]rune(getStatusLine(editor)))
		this.renderCursor(this.EditorBoxLeft+1, this.EditorBoxTop+1)
	}

	this.Screen.Show()
}
