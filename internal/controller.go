package internal

import (
	"log"

	"github.com/gdamore/tcell/v2"
)

func Start(fileNames []string) {
	screen, quit := NewScreen()
	defer quit()

	app := NewApplication()
  app.OpenAll(fileNames)
	defer app.CloseAll()

	// Event loop
	for !app.QuitSignal {
		// Update screen
		refreshScreen(screen, app)

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
