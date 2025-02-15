package view

import (
	"go-editor/utils"

	"github.com/gdamore/tcell/v2"
)

type ViewBuffer struct {
	TopX   int
	TopY   int
	Width  int
	Height int
	Text   [][]rune
}

func NewViewBuffer(topX int, topY int, width int, height int) *ViewBuffer {
	buf := &ViewBuffer{}
	buf.TopX = topX
	buf.TopY = topY
	buf.Width = width
	buf.Height = height
	buf.Text = utils.GetStringMatrix(height, width, ' ')

	return buf
}

func (buf *ViewBuffer) Render(screen tcell.Screen, style tcell.Style) {
	for y := 0; y < buf.Height; y++ {
		for x := 0; x < buf.Width; x++ {
			screen.SetContent(x+buf.TopX, y+buf.TopY, buf.Text[y][x], nil, style)
		}
	}

	screen.Show()
}

func (buf *ViewBuffer) Add(other *ViewBuffer) *ViewBuffer {
	relX := other.TopX - buf.TopX
	relY := other.TopY - buf.TopY

	for y := 0; y < other.Height; y++ {
		if y+relY >= buf.Height {
			break
		}
		for x := 0; x < other.Width; x++ {
			if x+relX >= buf.Width {
				break
			}
			buf.Text[y+relY][x+relX] = other.Text[y][x]
		}
	}

	return buf
}

func (buf *ViewBuffer) AddBorder() *ViewBuffer {
	for y := 0; y < buf.Height; y++ {
		buf.Text[y][0] = '|'
		buf.Text[y][buf.Width-1] = '|'
	}

	for x := 0; x < buf.Width; x++ {
		buf.Text[0][x] = '-'
		buf.Text[buf.Height-1][x] = '-'
	}

	return buf
}

func (buf *ViewBuffer) AddText(lines [][]rune) *ViewBuffer {
	for y := 0; y < len(lines); y++ {
		if y >= buf.Height {
			break
		}
		for x := 0; x < len(lines[y]); x++ {
			if x >= buf.Width {
				break
			}
			buf.Text[y][x] = lines[y][x]
		}
	}

	return buf
}
