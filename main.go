package main

import (
	"github.com/gdamore/tcell/v2"
	"go-editor/config"
	"go-editor/internal"
	"log"
  "os"
)

func main() {
  logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	// Redirect log output to the file
	log.SetOutput(logFile)

	var editor *internal.Application = internal.NewApplication()

	if err := editor.OpenFile(config.DEFAULT_FILENAME); err != nil {
		log.Fatal(err)
	}

	if err := editor.CurrentFile.ReadFile(); err != nil {
		log.Fatal(err)
	}

	defer editor.CloseAll()

	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

	// Initialize screen
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

	defer quit()

	// Event loop
	for {
		// Update screen
		refreshScreen(s, editor)

		// Poll event
		ev := s.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				return
			} else if ev.Key() == tcell.KeyCtrlL {
				s.Sync()
			} else {
				handleKey(ev.Rune(), editor)
			}
		case *tcell.EventMouse:
			log.Println("Mouse is not supported")
		}
	}
}
