package view

import (
	"fmt"
	"go-editor/model"
	"go-editor/utils"
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

func NewViewModel(screen tcell.Screen) *ViewModel {
	width, height := screen.Size()

	viewModel := new(ViewModel)
	viewModel.Screen = screen

	viewModel.EditorBoxTop = 0
	viewModel.EditorBoxLeft = 0
	viewModel.EditorBoxWidth = width
	viewModel.EditorBoxHeight = height - 2

	viewModel.StatusBoxTop = viewModel.EditorBoxHeight + 1
	viewModel.StatusBoxLeft = 0
	viewModel.StatusBoxWidth = width
	viewModel.StatusBoxHeight = height - viewModel.EditorBoxHeight

	return viewModel
}

func getModeName(mode int) string {
	switch mode {
	case model.NORMAL_MODE:
		return "NORMAL"
	case model.COMMAND_MODE:
		return "COMMAND"
	case model.INSERT_MODE:
		return "INSERT"
	default:
		return "UNNAMED"
	}
}

func getStringMatrix(nRows int, nCols int, char rune) [][]rune {
	lines := make([][]rune, nRows)

	for y := 0; y < nRows; y++ {
		lines[y] = make([]rune, nCols)
		for x := 0; x < nCols; x++ {
			lines[y][x] = char
		}
	}

	return lines
}

func (view *ViewModel) renderFile(file *model.File) {
	nRows, nCols := view.EditorBoxHeight, view.EditorBoxWidth

	var displayLines [][]rune = getStringMatrix(nRows, nCols, ' ')

	utils.Assert(len(displayLines) == nRows)
	utils.Assert(len(displayLines[0]) == nCols)

	if file != nil {
		for y := 0; y < nRows; y++ {
			if y+file.ScrollY >= len(file.Lines) {
				break
			}

			curLine := file.Lines[y+file.ScrollY].Text

			for x := 0; x < nCols; x++ {
				if x+file.ScrollX >= len(curLine) {
					break
				}

				displayLines[y][x] = curLine[x+file.ScrollX]
			}
		}
	}

	for y, line := range displayLines {
		for x, char := range line {
			drawY := view.EditorBoxTop + y
			drawX := view.EditorBoxLeft + x

			view.Screen.SetContent(drawX, drawY, char, nil, tcell.StyleDefault)
		}
	}
}

func (view *ViewModel) getRelativeCursor(file *model.File) (cursorX int, cursorY int) {
	cursorX = 0
	cursorY = 0

	if file == nil {
		return
	}

	utils.Assert(file.GetCurrentLine() != nil)

	cursorX = file.GetCurrentLine().Cursor - file.ScrollX
	cursorY = file.CursorLine - file.ScrollY

	if cursorY < 0 {
		log.Print("WARN: cursor out of bounds")
		cursorY = 0
	}

	if cursorX < 0 {
		log.Print("WARN: cursor out of bounds")
		cursorX = 0
	}

	return
}

func (view *ViewModel) getAbsoluteCursor(file *model.File) (cursorX int, cursorY int) {
	cursorX, cursorY = view.getRelativeCursor(file)
	cursorX += view.EditorBoxLeft
	cursorY += view.EditorBoxTop
	return
}

func (view *ViewModel) renderStatus(editor *model.Application) {
	filename := "No File"
	if editor.CurrentFile != nil {
		filename = editor.CurrentFile.Name
		if editor.CurrentFile.Modified {
			filename += " [+]"
		}
	}

	mode := getModeName(editor.Mode)

	status := fmt.Sprintf("[%s] | %s | Ln %d, Col %d | file %d of %d",
		mode,
		filename,
		editor.CurrentFile.CursorLine,
		editor.CurrentFile.GetCurrentLine().Cursor,
		editor.GetCurrentFileIndex()+1,
		len(editor.Files))

	line := []rune(status)

	DrawText(view.Screen, view.StatusBoxLeft, view.StatusBoxTop,
		view.StatusBoxLeft+view.StatusBoxWidth,
		view.StatusBoxTop+view.StatusBoxHeight, tcell.StyleDefault, line)
}

func (view *ViewModel) renderCommandPrompt(command string) {
	commandBoxWidth := view.EditorBoxWidth / 2
	commandBoxHeight := 5 // TODO
	commandBoxLeft := view.EditorBoxLeft + (view.EditorBoxWidth/2 - commandBoxWidth/2)
	commandBoxTop := view.EditorBoxTop + (view.EditorBoxHeight/2 - commandBoxHeight/2)

	DrawBox(view.Screen, commandBoxLeft, commandBoxTop, commandBoxLeft+commandBoxWidth, commandBoxTop+commandBoxHeight, tcell.StyleDefault, []rune(command))
}

func (view *ViewModel) Render(editor *model.Application) {
	view.Screen.Clear()

	if editor.CurrentFile != nil {
		editor.CurrentFile.AdjustScroll(view.EditorBoxHeight, view.EditorBoxWidth)
	}

	view.renderFile(editor.CurrentFile)
	view.renderStatus(editor)

	cursorX, cursorY := view.getAbsoluteCursor(editor.CurrentFile)
	view.Screen.ShowCursor(cursorX, cursorY)

	if editor.Mode == model.COMMAND_MODE {
		view.renderCommandPrompt(":" + editor.UserCommand)
	}

	view.Screen.Show()
}
