package internal

import (
	"go-editor/config"
	"log"
)

type File struct {
	Lines      []*Line
	Name       string
	Modified   bool
	Readonly   bool
	CursorLine int
	ScrollX    int
	ScrollY    int
}

func (file *File) AdjustScroll(viewModel *ViewModel) {
	// adjust horizontal scroll
	line := file.GetCurrentLine()
	if line.Cursor-file.ScrollX < 0 {
		if config.DEBUG {
			log.Println("scrolling left")
		}
		file.ScrollX = line.Cursor
	} else if line.Cursor-file.ScrollX >= viewModel.GetMaxDisplayCols() {
		if config.DEBUG {
			log.Println("scrolling right")
		}
		file.ScrollX = line.Cursor - viewModel.GetMaxDisplayCols() + 1
	}

	// adjust vertical scroll
	if file.CursorLine-file.ScrollY < 0 {
		if config.DEBUG {
			log.Println("scrolling up")
		}
		file.ScrollY = file.CursorLine
	} else if file.CursorLine-file.ScrollY >= viewModel.GetMaxDisplayLines() {
		if config.DEBUG {
			log.Println("scrolling down")
		}
		file.ScrollY = file.CursorLine - viewModel.GetMaxDisplayLines() + 1
	}
}

func NewFile(filename string) *File {
	var file *File = new(File)
	file.Name = filename
	file.Modified = false
	file.Readonly = false
	file.Lines = make([]*Line, 0)
	file.CursorLine = 0
	file.ScrollX = 0
	file.ScrollY = 0
	return file
}

func (file *File) ReadFile() error {
	lines, err := ReadFileUtil(file.Name)
	if err != nil {
		return err
	}

	file.Lines = make([]*Line, 0)

	for _, line := range lines {
		file.Lines = append(file.Lines, NewLine1(line))
	}

	file.Modified = false
	return nil
}

func (file *File) WriteFile() error {
	if file.Readonly {
		return ErrorFileNotModifiable()
	}

	raw_lines := make([]string, 0)
	for _, line := range file.Lines {
		raw_lines = append(raw_lines, string(line.Text))
	}

	err := WriteFileUtil(file.Name, raw_lines)
	if err != nil {
		return err
	}

	file.Modified = true
	return nil
}

func (file *File) InsertLineAboveCursor() {
	newLines := make([]*Line, 0)
	newLines = append(newLines, file.Lines[:file.CursorLine]...)
	newLines = append(newLines, NewLine())
	newLines = append(newLines, file.Lines[file.CursorLine:]...)
	file.CursorLine += 1
	file.Lines = newLines
	file.Modified = true
}

func (file *File) InsertLineBelowCursor() {
	if file.CursorLine == len(file.Lines)-1 {
		file.Lines = append(file.Lines, NewLine())
	} else {
		file.CursorLine += 1
		file.InsertLineAboveCursor()
		file.CursorLine -= 2
	}
	file.Modified = true
}

// delete line at cursor
func (file *File) DeleteLine() {
	if len(file.Lines) == 0 {
		return
	}

	newLines := make([]*Line, 0)
	newLines = append(newLines, file.Lines[:file.CursorLine]...)

	if file.CursorLine+1 < len(file.Lines) {
		newLines = append(newLines, file.Lines[file.CursorLine+1:]...)
	}

	file.Lines = newLines

	if file.CursorLine > 0 && file.CursorLine == len(file.Lines)-1 {
		file.CursorLine -= 1
	}

	file.Modified = true
}

func (file *File) GetCurrentLine() *Line {
	if file.CursorLine < 0 || file.CursorLine >= len(file.Lines) {
		log.Println("no lines in file")
		return nil
	}
	return file.Lines[file.CursorLine]
}

func (file *File) MoveDown() {
	if file.CursorLine+1 >= len(file.Lines) {
		return
	}
	prevCursorX := file.GetCurrentLine().Cursor
	file.CursorLine += 1
	file.GetCurrentLine().SetCursor(min(prevCursorX, file.GetCurrentLine().Size()-1))
}

func (file *File) MoveUp() {
	if file.CursorLine-1 < 0 {
		return
	}
	prevCursorX := file.GetCurrentLine().Cursor
	file.CursorLine -= 1
	file.GetCurrentLine().SetCursor(min(prevCursorX, file.GetCurrentLine().Size()-1))
}

func (file *File) MoveForward() {
	file.GetCurrentLine().MoveForward()
}

func (file *File) MoveBackward() {
	file.GetCurrentLine().MoveBackward()
}

func (file *File) Insert(r rune) {
	file.GetCurrentLine().Insert(r)
}

func (file *File) GetXCursor() int {
	return file.GetCurrentLine().Cursor
}

func (file *File) SetXCursor(value int) {
	if value < 0 || value >= file.GetCurrentLine().Cursor {
		return
	}

	prevCursorX := file.GetCurrentLine().Cursor
	if value == prevCursorX { // noop
		return
	}

	file.GetCurrentLine().Cursor = value
}

func (file *File) SetYCursor(value int) {
	if value < 0 || value >= file.CountLines() {
		return
	}

	prevCursorY := file.CursorLine
	if value == prevCursorY { // noop
		return
	}

	file.CursorLine = value
}

func (file *File) CountLines() int {
	return len(file.Lines)
}

func (file *File) JumpToNextChar(c rune) {
	file.GetCurrentLine().JumpToNextChar(c)
}

// func (file *File) Paste(text string) {
// 	if len(text) == 0 {
// 		return
// 	}
// 	file.Modified = true

// 	insertLinesRaw := strings.Split(text, "\n")
// 	insertLines := make([]Line, 0)
// 	for _, line := range insertLinesRaw {
// 		insertLines = append(insertLines, Line{Text: line})
// 	}
// 	if len(insertLines) == 1 {
// 		var newLine bytes.Buffer
// 		newLine.WriteString(file.Lines[file.CursorLine].Text[:file.CursorX])
// 		newLine.WriteString(insertLines[0].Text)
// 		newLine.WriteString(file.Lines[file.CursorLine].Text[file.CursorX:])
// 		file.Lines[file.CursorLine].Text = newLine.String()
// 		return
// 	}

// 	var newLine bytes.Buffer
// 	newLine.WriteString(file.Lines[file.CursorLine].Text[:file.CursorX])
// 	newLine.WriteString(insertLines[0].Text)
// 	file.Lines[file.CursorLine].Text = newLine.String()

// 	brokenLinePart := file.Lines[file.CursorLine].Text[file.CursorX:]
// 	remainingLines := make([]Line, 0)
// 	remainingLines = append(remainingLines, Line{Text: brokenLinePart})
// 	remainingLines = append(remainingLines, insertLines[1:]...)

// 	newLines := make([]Line, 0)
// 	newLines = append(newLines, file.Lines[:file.CursorLine]...)
// 	newLines = append(newLines, remainingLines...)
// 	newLines = append(newLines, file.Lines[file.CursorLine+1:]...)

// 	file.Lines = newLines
// }
