package model

import (
	"strings"
)

type Line struct {
	Text     []rune
	Cursor   int
	Modified bool
}

func NewLine() *Line {
	return &Line{Text: []rune{}, Cursor: 0, Modified: false}
}

func NewLine1(text string) *Line {
	return &Line{Text: []rune(text), Cursor: 0, Modified: false}
}

func (line *Line) Size() int {
	return len(line.Text)
}

func (line *Line) SetCursor(cursor int) {
	if cursor < 0 || cursor >= len(line.Text) {
		return
	}

	line.Cursor = cursor
}

// Clear all characters from line
func (line *Line) Clear() {
	line.Text = []rune{}
	line.Cursor = 0
	line.Modified = true
}

// Move cursor to last position
func (line *Line) MoveToEnd() {
	line.Cursor = len(line.Text) - 1
}

// Move cursor to fist position
func (line *Line) MoveToStart() {
	line.Cursor = 0
}

func (line *Line) MoveBackward() {
	if line.Cursor > 0 {
		line.Cursor -= 1
	}
}

func (line *Line) MoveForward() {
	if line.Cursor+1 < len(line.Text) {
		line.Cursor += 1
	}
}

func (line *Line) MoveBackwardN(count int) {
	if line.Cursor-count < 0 {
		line.Cursor = 0
	} else {
		line.Cursor -= count
	}
}

func (line *Line) MoveForwardN(count int) {
	if line.Cursor+count >= len(line.Text) {
		line.MoveToEnd()
	} else {
		line.Cursor += count
	}
}

// Inserts given string at cursor up to the first newline
func (line *Line) InsertString(text string) {
	text = strings.Split(text, "\n")[0]
	newText := []rune{}
	newText = append(newText, line.Text[:line.Cursor]...)
	newText = append(newText, []rune(text)...)
	newText = append(newText, line.Text[line.Cursor:]...)

	line.Text = newText
	line.Cursor += len(text)
	line.Modified = true
}

// Insert character at cursor and advance
func (line *Line) Insert(r rune) {
	newText := []rune{}
	newText = append(newText, line.Text[:line.Cursor]...)
	newText = append(newText, r)
	newText = append(newText, line.Text[line.Cursor:]...)

	line.Text = newText
	line.Cursor += 1
	line.Modified = true
}

// Remove character at cursor
func (line *Line) RemoveChar() {
	if len(line.Text) == 0 || line.Cursor == 0 {
		return
	}

	if line.Cursor >= len(line.Text) {
		line.Text = line.Text[:len(line.Text)]
		line.MoveToEnd()
		return
	}

	line.Text = append(line.Text[:line.Cursor], line.Text[line.Cursor+1:]...)
	line.Modified = true
	if line.Cursor >= len(line.Text) {
		line.MoveToEnd()
	}
}

// Remove *count* characters from cursor
func (line *Line) RemoveN(count int) {
	maxRemovable := len(line.Text) - line.Cursor
	if count > maxRemovable {
		count = maxRemovable
	}
	if count <= 0 {
		return
	}
	line.Text = append(line.Text[:line.Cursor], line.Text[line.Cursor+count:]...)
	line.Modified = true
	if line.Cursor >= len(line.Text) {
		line.MoveToEnd()
	}
}

// Insert string after cursor and advance
func (line *Line) AppendString(text string) {
	if line.Cursor+1 >= len(line.Text) {
		line.Text = append(line.Text, []rune(text)...)
	} else {
		newText := []rune{}
		newText = append(newText, line.Text[:line.Cursor+1]...)
		newText = append(newText, []rune(text)...)
		newText = append(newText, line.Text[line.Cursor+1+len(text):]...)
		line.Text = newText
	}

	line.Cursor += len(text)
	line.Modified = true

}

// Insert character after cursor and advance
func (line *Line) Append(r rune) {
	if line.Cursor+1 >= len(line.Text) {
		line.Text = append(line.Text, r)
	} else {
		newText := []rune{}
		newText = append(newText, line.Text[:line.Cursor+1]...)
		newText = append(newText, r)
		newText = append(newText, line.Text[line.Cursor+2:]...)
		line.Text = newText
	}

	line.Cursor += 1
	line.Modified = true
}

func (line *Line) JumpToNextChar(c rune) {
	for i := line.Cursor + 1; i < line.Size(); i++ {
		if line.Text[i] == c {
			line.Cursor = i
			break
		}
	}
}

func (line *Line) GetVisibleText(nCols int, scrollX int) []rune {
	if scrollX >= len(line.Text) {
		return []rune{}
	}

	if scrollX+nCols >= len(line.Text) {
		return line.Text[scrollX:]
	}

	return line.Text[scrollX : scrollX+nCols]
}
