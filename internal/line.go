package internal

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

func (this *Line) Size() int {
	return len(this.Text)
}

func (this *Line) SetCursor(cursor int) {
	if cursor < 0 || cursor >= len(this.Text) {
		return
	}

	this.Cursor = cursor
}

// Clear all characters from line
func (this *Line) Clear() {
	this.Text = []rune{}
	this.Cursor = 0
	this.Modified = true
}

// Move cursor to last position
func (this *Line) MoveToEnd() {
	this.Cursor = len(this.Text) - 1
}

// Move cursor to fist position
func (this *Line) MoveToStart() {
	this.Cursor = 0
}

func (this *Line) MoveBackward() {
	if this.Cursor > 0 {
		this.Cursor -= 1
	}
}

func (this *Line) MoveForward() {
	if this.Cursor+1 < len(this.Text) {
		this.Cursor += 1
	}
}

func (this *Line) MoveBackwardN(count int) {
	if this.Cursor-count < 0 {
		this.Cursor = 0
	} else {
		this.Cursor -= count
	}
}

func (this *Line) MoveForwardN(count int) {
	if this.Cursor+count >= len(this.Text) {
		this.MoveToEnd()
	} else {
		this.Cursor += count
	}
}

// Inserts given string at cursor up to the first newline
func (this *Line) InsertString(text string) {
	text = strings.Split(text, "\n")[0]
	this.Text = append(this.Text[:this.Cursor], append([]rune(text), this.Text[this.Cursor:]...)...)
	this.Cursor += len(text)
	this.Modified = true
}

// Insert character at cursor and advance
func (this *Line) Insert(r rune) {
	this.Text = append(this.Text[:this.Cursor], append([]rune{r}, this.Text[this.Cursor:]...)...)
	this.Cursor += 1
	this.Modified = true
}

// Remove character at cursor
func (this *Line) RemoveChar() {
	if len(this.Text) == 0 {
		return
	}

	this.Text = append(this.Text[:this.Cursor], this.Text[this.Cursor+1:]...)
	this.Modified = true
	if this.Cursor >= len(this.Text) {
		this.MoveToEnd()
	}
}

// Remove *count* characters from cursor
func (this *Line) RemoveN(count int) {
	maxRemovable := len(this.Text) - this.Cursor
	if count > maxRemovable {
		count = maxRemovable
	}
	if count <= 0 {
		return
	}
	this.Text = append(this.Text[:this.Cursor], this.Text[this.Cursor+count:]...)
	this.Modified = true
	if this.Cursor >= len(this.Text) {
		this.MoveToEnd()
	}
}

// Insert string after cursor and advance
func (this *Line) AppendString(text string) {
	this.Cursor += 1
	this.InsertString(text)
}

// Insert character after cursor and advance
func (this *Line) Append(r rune) {
	this.Cursor += 1 // to insert after cursor
	this.Insert(r)
}

func (this *Line) JumpToNextChar(c rune) {
	for i := this.Cursor; i < this.Size(); i++ {
		if this.Text[i] == c {
			this.Cursor = i
		}
	}
}
