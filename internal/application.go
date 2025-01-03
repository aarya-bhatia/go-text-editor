package internal

import "log"

const NORMAL_MODE int = 0
const INSERT_MODE int = 1
const COMMAND_MODE int = 2

type Application struct {
	QuitSignal  bool
	Files       []*File
	CurrentFile *File
	StatusLine  string
	Mode        int
}

// Returns new application instance
func NewApplication() *Application {
	var this *Application = new(Application)
	this.QuitSignal = false
	this.Files = make([]*File, 0)
	this.CurrentFile = nil
	this.StatusLine = ""
	this.Mode = NORMAL_MODE
	return this
}

// Create new file and make it current
func (this *Application) AddFile(filename string) error {
	var file *File = NewFile(filename)
	this.Files = append(this.Files, file)
	this.CurrentFile = file
  log.Println("new file added: ", filename)

	return nil
}

// Find file in buffer list or create it and make it current
func (this *Application) OpenFile(filename string) error {
	for _, file := range this.Files {
		if file.Name == filename {
			this.CurrentFile = file
			return nil
		}
	}

	return this.AddFile(filename)
}

// Close current file and replace it with the next file in buffer list.
func (this *Application) CloseFile() {
	if this.CurrentFile == nil || len(this.Files) == 0 {
    log.Println("buffer list is empty")
		return
	}

	if len(this.Files) == 1 {
		this.CurrentFile = nil
		this.Files = make([]*File, 0)
		return
	}

	for i, file := range this.Files {
		if file == this.CurrentFile {

			if i+1 < len(this.Files) {
				this.CurrentFile = this.Files[i+1]
			} else {
				this.CurrentFile = this.Files[0]
			}

			newFiles := make([]*File, 0)
			newFiles = append(newFiles, this.Files[:i]...)
			newFiles = append(newFiles, this.Files[i+1:]...)

      this.Files = newFiles
			return
		}
	}
}

// Open the next file in buffer list
func (this *Application) OpenNextFile() {
	if this.CurrentFile == nil && len(this.Files) == 0 {
		return
	}

	if this.CurrentFile == nil {
		this.CurrentFile = this.Files[0]
		return
	}

	if this.CurrentFile == this.Files[len(this.Files)-1] {
		this.CurrentFile = this.Files[0]
		return
	}

	for i, file := range this.Files[:len(this.Files)-1] {
		if file == this.CurrentFile {
			this.CurrentFile = this.Files[i+1]
			return
		}
	}

  log.Println("No next file")
}

// Open the previous file in buffer list
func (this *Application) OpenPrevFile() {
	if this.CurrentFile == nil && len(this.Files) == 0 {
		return
	}

	if this.CurrentFile == nil {
		this.CurrentFile = this.Files[len(this.Files)-1]
		return
	}

	if this.CurrentFile == this.Files[0] {
		this.CurrentFile = this.Files[len(this.Files)-1]
		return
	}

	for i, file := range this.Files[1:] {
		if file == this.CurrentFile {
			this.CurrentFile = this.Files[(i+1)-1] // offset i because it is relative to 1-indexed array
			return
		}
	}

  log.Println("No prev file")
}
