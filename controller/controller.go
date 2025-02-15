package controller

import (
	"log"

	"go-editor/model"
	"go-editor/view"

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
	viewBuffer := view.NewViewBuffer(0, 0, nCols, nRows)
	viewBuffer.Add(view.NewViewBuffer(0, 0, nCols, nRows-1).AddBorder())

	// Event loop
	for !app.QuitSignal {
		// Add status
		statusView := view.NewViewBuffer(0, nRows-1, nCols, 1)
		statusView.AddText([][]rune{app.GetStatusLine()})
		viewBuffer.Add(statusView)

		editorView := view.NewViewBuffer(1, 1, nCols-2, nRows-3)

		if app.CurrentFile != nil {
			app.CurrentFile.AdjustScroll(editorView.Height, editorView.Width)

			editorView.AddText(app.CurrentFile.GetVisibleText(nRows, nCols))
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
		case *tcell.EventKey:
			handleKeyEvent(ev, app, screen)
		case *tcell.EventMouse:
			log.Println("Mouse is not supported")
		}
	}
}
