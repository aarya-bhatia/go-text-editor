package internal

type File struct {
	Lines    []*Line
	Name     string
	Modified bool
	Readonly bool
	CursorLine int

	ScrollX  int
	ScrollY  int
}

func NewFile(filename string) *File {
	var this *File = new(File)
	this.Name = filename
	this.Modified = false
	this.Readonly = false
	this.Lines = make([]*Line, 0)
	this.CursorLine = 0
	this.ScrollX = 0
	this.ScrollY = 0
	return this
}

func (this *File) ReadFile() error {
	lines, err := ReadFileUtil(this.Name)
	if err != nil {
		return err
	}

	this.Lines = make([]*Line, 0)

	for _, line := range lines {
		this.Lines = append(this.Lines, NewLine1(line))
	}

	this.Modified = false
	return nil
}

func (this *File) WriteFile() error {
	if this.Readonly {
		return ErrorFileNotModifiable()
	}

	raw_lines := make([]string, 0)
	for _, line := range this.Lines {
		raw_lines = append(raw_lines, line.Text)
	}

	err := WriteFileUtil(this.Name, raw_lines)
	if err != nil {
		return err
	}

	this.Modified = true
	return nil
}

func (this *File) InsertLineAboveCursor() {
	newLines := make([]*Line, 0)
	newLines = append(newLines, this.Lines[:this.CursorLine]...)
	newLines = append(newLines, NewLine())
	newLines = append(newLines, this.Lines[this.CursorLine:]...)
	this.CursorLine += 1
	this.Lines = newLines
	this.Modified = true
}

func (this *File) InsertLineBelowCursor() {
	newLines := make([]*Line, 0)
	newLines = append(newLines, this.Lines[:this.CursorLine+1]...)
	newLines = append(newLines, NewLine())
	newLines = append(newLines, this.Lines[this.CursorLine+1:]...)
	this.Lines = newLines
	this.Modified = true
}

// func (this *File) Paste(text string) {
// 	if len(text) == 0 {
// 		return
// 	}
// 	this.Modified = true

// 	insertLinesRaw := strings.Split(text, "\n")
// 	insertLines := make([]Line, 0)
// 	for _, line := range insertLinesRaw {
// 		insertLines = append(insertLines, Line{Text: line})
// 	}
// 	if len(insertLines) == 1 {
// 		var newLine bytes.Buffer
// 		newLine.WriteString(this.Lines[this.CursorLine].Text[:this.CursorX])
// 		newLine.WriteString(insertLines[0].Text)
// 		newLine.WriteString(this.Lines[this.CursorLine].Text[this.CursorX:])
// 		this.Lines[this.CursorLine].Text = newLine.String()
// 		return
// 	}

// 	var newLine bytes.Buffer
// 	newLine.WriteString(this.Lines[this.CursorLine].Text[:this.CursorX])
// 	newLine.WriteString(insertLines[0].Text)
// 	this.Lines[this.CursorLine].Text = newLine.String()

// 	brokenLinePart := this.Lines[this.CursorLine].Text[this.CursorX:]
// 	remainingLines := make([]Line, 0)
// 	remainingLines = append(remainingLines, Line{Text: brokenLinePart})
// 	remainingLines = append(remainingLines, insertLines[1:]...)

// 	newLines := make([]Line, 0)
// 	newLines = append(newLines, this.Lines[:this.CursorLine]...)
// 	newLines = append(newLines, remainingLines...)
// 	newLines = append(newLines, this.Lines[this.CursorLine+1:]...)

// 	this.Lines = newLines
// }

func (this *File) DeleteLine() {
	newLines := make([]*Line, 0)
	newLines = append(newLines, this.Lines[:this.CursorLine]...)
	newLines = append(newLines, this.Lines[this.CursorLine+1:]...)
	this.Lines = newLines
	this.Modified = true
}

func (this *File) GetCurrentLine() *Line {
	return this.Lines[this.CursorLine]
}
