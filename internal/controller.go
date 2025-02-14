package internal

import (
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

func Start(fileNames []string) {
	screen, quit := NewScreen()
	defer quit()

	viewModel := NewViewModel(screen)

	app := NewApplication()

	if len(fileNames) == 0 {
		app.OpenTempFile()
	} else {
		app.OpenAll(fileNames)
	}

	defer app.CloseAll()

	// Event loop
	for !app.QuitSignal {
		// Update screen
		viewModel.render(app)

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
