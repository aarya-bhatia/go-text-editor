package main

import (
	"github.com/gdamore/tcell/v2"
	"go-editor/internal"
	"log"
)

func main() {
	var editor *internal.Application = internal.NewApplication()

	if err := editor.OpenFile("hello.txt"); err != nil {
		log.Fatal(err)
	}

  editor.CurrentFile.WriteFile()

	defer editor.CloseAll()

	// Initialize the tcell screen
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("Failed to create screen: %v", err)
	}
	defer screen.Fini()

	if err := screen.Init(); err != nil {
		log.Fatalf("Failed to initialize screen: %v", err)
	}

	// Set screen style
	screen.Clear()

	// Box dimensions
	editorBoxTop := 2
	editorBoxLeft := 2
	editorBoxWidth := 50
	editorBoxHeight := 10

	statusBoxTop := editorBoxTop + editorBoxHeight + 1
	statusBoxLeft := editorBoxLeft
	statusBoxWidth := editorBoxWidth
	statusBoxHeight := 3

	// Text and cursor position
	editorText := []rune{}
	cursorX, cursorY := 0, 0

	// Helper to draw a box
	drawBox := func(x, y, width, height int) {
		for i := 0; i < width; i++ {
			screen.SetContent(x+i, y, tcell.RuneHLine, nil, tcell.StyleDefault) // Top border
			screen.SetContent(x+i, y+height-1, tcell.RuneHLine, nil, tcell.StyleDefault) // Bottom border
		}
		for i := 0; i < height; i++ {
			screen.SetContent(x, y+i, tcell.RuneVLine, nil, tcell.StyleDefault) // Left border
			screen.SetContent(x+width-1, y+i, tcell.RuneVLine, nil, tcell.StyleDefault) // Right border
		}
		screen.SetContent(x, y, tcell.RuneULCorner, nil, tcell.StyleDefault)            // Top-left corner
		screen.SetContent(x+width-1, y, tcell.RuneURCorner, nil, tcell.StyleDefault)    // Top-right corner
		screen.SetContent(x, y+height-1, tcell.RuneLLCorner, nil, tcell.StyleDefault)   // Bottom-left corner
		screen.SetContent(x+width-1, y+height-1, tcell.RuneLRCorner, nil, tcell.StyleDefault) // Bottom-right corner
	}

	// Helper to render text
	renderText := func(x, y int, text []rune) {
		for i, r := range text {
			screen.SetContent(x+i, y, r, nil, tcell.StyleDefault)
		}
	}

	// Render the cursor
	renderCursor := func(x, y int) {
		screen.ShowCursor(x, y)
	}

	// Initial draw
	drawBox(editorBoxLeft, editorBoxTop, editorBoxWidth, editorBoxHeight)
	drawBox(statusBoxLeft, statusBoxTop, statusBoxWidth, statusBoxHeight)
	renderText(statusBoxLeft+1, statusBoxTop+1, []rune("Press Ctrl+C to exit."))
	screen.Show()

	// Event loop
	for {
		// Render editor content
		for i := 0; i < len(editorText)/editorBoxWidth+1; i++ {
			lineStart := i * (editorBoxWidth - 2)
			lineEnd := lineStart + (editorBoxWidth - 2)
			if lineEnd > len(editorText) {
				lineEnd = len(editorText)
			}
			renderText(editorBoxLeft+1, editorBoxTop+1+i, editorText[lineStart:lineEnd])
		}
		renderCursor(editorBoxLeft+1+cursorX, editorBoxTop+1+cursorY)
		screen.Show()

		// Wait for an event
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyCtrlC: // Exit on Ctrl+C
				return
			case tcell.KeyBackspace, tcell.KeyBackspace2: // Handle backspace
				if len(editorText) > 0 {
					editorText = editorText[:len(editorText)-1]
					if cursorX > 0 {
						cursorX--
					} else if cursorY > 0 {
						cursorY--
						cursorX = editorBoxWidth - 3
					}
				}
			case tcell.KeyRune: // Handle regular character input
				editorText = append(editorText, ev.Rune())
				if cursorX < editorBoxWidth-3 {
					cursorX++
				} else {
					cursorX = 0
					cursorY++
				}
			}
		case *tcell.EventResize:
			screen.Sync() // Handle terminal resize
		}
	}
}

