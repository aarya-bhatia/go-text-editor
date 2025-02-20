package controller

import (
	"log"
	"strconv"

	"go-editor/model"
	"go-editor/view"

	"github.com/gdamore/tcell/v2"
)

var editorWidth = 0
var editorHeight = 0
var editorTopX = 0
var editorTopY = 0

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

func setupBaseView(nRows int, nCols int) *view.ViewBuffer {
	viewBuffer := view.NewViewBuffer(0, 0, nCols, nRows)
	viewBuffer.Add(view.NewViewBuffer(0, 0, nCols, nRows-2).AddBorder())

	editorWidth = nCols - 2
	editorHeight = nRows - 4
	editorTopX = 1
	editorTopY = 1

	lineNumWidth := len(strconv.Itoa(nRows - 4))

	if editorWidth-lineNumWidth >= 10 { // hide line numbers on very narrow screen.
		lineNumberView := view.NewViewBuffer(editorTopX, editorTopY, lineNumWidth, editorHeight).AddLineNumbers()
		viewBuffer.Add(lineNumberView)

		separatorView := view.NewViewBuffer(editorTopX+lineNumWidth, editorTopY, 1, editorHeight)
		viewBuffer.Add(separatorView)

		editorWidth -= lineNumWidth + 1
		editorTopX += lineNumWidth + 1
	}

	return viewBuffer
}

func Start(fileNames []string) {
	screen, quit := NewScreen()
	defer quit()

	app := model.NewApplication()

	if len(fileNames) == 0 {
		app.OpenTempFile()
	} else {
		app.OpenAll(fileNames)
	}

	defer app.CloseAll()

	nCols, nRows := screen.Size()
	viewBuffer := setupBaseView(nRows, nCols)

	// Event loop
	for !app.QuitSignal {
		// Add status
		statusView := view.NewViewBuffer(0, nRows-2, nCols, 1)
		statusView.AddText([][]rune{app.GetStatusLine()})
		viewBuffer.Add(statusView)

		commandView := view.NewViewBuffer(0, nRows-1, nCols, 1)
		if app.Mode == model.COMMAND_MODE {
			commandView.AddText([][]rune{[]rune(":" + app.UserCommand)})
		} else {
			commandView.AddText([][]rune{[]rune(app.StatusLine)})
		}
		viewBuffer.Add(commandView)

		editorView := view.NewViewBuffer(editorTopX, editorTopY, editorWidth, editorHeight)

		if app.CurrentFile != nil {
			app.CurrentFile.AdjustScroll(editorView.Height, editorView.Width)

			editorView.AddText(app.CurrentFile.GetVisibleText(editorHeight, editorWidth))
			viewBuffer.Add(editorView)

			cursorX, cursorY := app.CurrentFile.GetCursor()
			screen.ShowCursor(editorView.TopX+cursorX, editorView.TopY+cursorY)
		} else {
			screen.ShowCursor(editorView.TopX, editorView.TopY)
		}

		// Update screen content
		viewBuffer.Render(screen, tcell.StyleDefault)

		// Poll event
		ev := screen.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			screen.Sync()
			nCols, nRows = screen.Size()
			viewBuffer = setupBaseView(nRows, nCols)
		case *tcell.EventKey:
			handleKeyEvent(ev, app, screen)
		case *tcell.EventMouse:
			log.Println("Mouse is not supported")
		}
	}
}
