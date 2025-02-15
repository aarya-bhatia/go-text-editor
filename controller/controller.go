package controller

import (
	"log"

	"go-editor/model"
	"go-editor/view"

	"github.com/gdamore/tcell/v2"
)

func Start(fileNames []string) {
	screen, quit := view.NewScreen()
	defer quit()

	viewModel := view.NewViewModel(screen)

	app := model.NewApplication()

	if len(fileNames) == 0 {
		app.OpenTempFile()
	} else {
		app.OpenAll(fileNames)
	}

	defer app.CloseAll()

	// Event loop
	for !app.QuitSignal {
		// Update screen
		viewModel.Render(app)

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
