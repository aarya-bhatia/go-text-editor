package model

import (
	"fmt"
	"go-editor/utils"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

const NORMAL_MODE int = 0
const INSERT_MODE int = 1
const COMMAND_MODE int = 2
const NORMAL_MODE_ARG_PENDING = 3

type Application struct {
	QuitSignal  bool
	Files       []*File
	CurrentFile *File
	StatusLine  string
	Mode        int
	UserCommand string
}

// Returns new application instance
func NewApplication() *Application {
	var app *Application = new(Application)
	app.QuitSignal = false
	app.Files = make([]*File, 0)
	app.CurrentFile = nil
	app.StatusLine = ""
	app.Mode = NORMAL_MODE
	app.UserCommand = ""
	return app
}

func (app *Application) getModeName() string {
	switch app.Mode {
	case COMMAND_MODE:
		return "COMMAND"
	case INSERT_MODE:
		return "INSERT"
	default:
		return "NORMAL"
	}
}

func (app *Application) OpenTempFile() {
	f, err := os.CreateTemp("", "unnamed")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	app.AddFile(f.Name())
}

func (app *Application) GetCurrentFileIndex() int {
	for i, file := range app.Files {
		if file == app.CurrentFile {
			return i
		}
	}
	return -1
}

// Create new file and make it current
func (app *Application) AddFile(filename string) error {
	var file *File = NewFile(filename)
	err := file.ReadFile()
	if err != nil {
		return err
	}

	app.Files = append(app.Files, file)
	app.CurrentFile = file
	log.Info("new file added: ", filename)

	return nil
}

// Find file in buffer list or create it and make it current
func (app *Application) OpenFile(filename string) error {
	for _, file := range app.Files {
		if file.Name == filename {
			app.CurrentFile = file
			return nil
		}
	}

	return app.AddFile(filename)
}

// Close current file and replace it with the next file in buffer list.
func (app *Application) CloseFile() {
	if app.CurrentFile == nil || len(app.Files) == 0 {
		log.Debug("buffer list is empty")
		return
	}

	if len(app.Files) == 1 {
		app.CurrentFile = nil
		app.Files = make([]*File, 0)
		return
	}

	for i, file := range app.Files {
		if file == app.CurrentFile {

			if i+1 < len(app.Files) {
				app.CurrentFile = app.Files[i+1]
			} else {
				app.CurrentFile = app.Files[0]
			}

			app.Files = utils.ArrayRemove(app.Files, i)
			return
		}
	}
}

// Close all files
func (app *Application) CloseAll() {
	app.Files = make([]*File, 0)
	app.CurrentFile = nil
}

// Open the next file in buffer list
func (app *Application) OpenNextFile() {
	if app.CurrentFile == nil && len(app.Files) == 0 {
		return
	}

	if app.CurrentFile == nil {
		app.CurrentFile = app.Files[0]
		return
	}

	if app.CurrentFile == app.Files[len(app.Files)-1] {
		app.CurrentFile = app.Files[0]
		return
	}

	for i, file := range app.Files[:len(app.Files)-1] {
		if file == app.CurrentFile {
			app.CurrentFile = app.Files[i+1]
			return
		}
	}

	log.Debug("No next file")
}

// Open the previous file in buffer list
func (app *Application) OpenPrevFile() {
	if app.CurrentFile == nil && len(app.Files) == 0 {
		return
	}

	if app.CurrentFile == nil {
		app.CurrentFile = app.Files[len(app.Files)-1]
		return
	}

	if app.CurrentFile == app.Files[0] {
		app.CurrentFile = app.Files[len(app.Files)-1]
		return
	}

	for i, file := range app.Files[1:] {
		if file == app.CurrentFile {
			app.CurrentFile = app.Files[(i+1)-1] // offset i because it is relative to 1-indexed array
			return
		}
	}

	log.Debug("No prev file")
}

func (app *Application) OpenAll(filenames []string) {
	for _, filename := range filenames {
		file := NewFile(filename)

		if err := file.ReadFile(); err != nil {
			log.Warnf("Failed to open file: %s", filename)
			continue
		}

		log.Debugf("Opened file '%s' with %d lines.", filename, file.CountLines())
		app.Files = append(app.Files, file)
	}

	if len(app.Files) > 0 {
		app.CurrentFile = app.Files[0]
	}
}

func (app *Application) Quit() {
	log.Info("Setting quit signal")
	app.QuitSignal = true
}

func (app *Application) GotoLine(lineNo int) {
	if app.CurrentFile != nil {
		log.Debug("Move to line ", lineNo)
		app.CurrentFile.SetYCursor(lineNo)
		return
	}
}

func (app *Application) GetStatusLine() []rune {
	widgets := []string{}
	widgets = append(widgets, fmt.Sprintf("[%s]", app.getModeName()))

	filename := "No File"
	if app.CurrentFile != nil {
		filename = app.CurrentFile.Name
		if app.CurrentFile.Modified {
			filename += " [+]"
		}
	}

	widgets = append(widgets, filename)

	if app.CurrentFile != nil {
		widgets = append(widgets, fmt.Sprintf("Ln %d, Col %d",
			app.CurrentFile.CursorLine, app.CurrentFile.GetCurrentLine().Cursor))
	}

	widgets = append(widgets, fmt.Sprintf("file %d of %d",
		app.GetCurrentFileIndex()+1,
		len(app.Files)),
	)

	status := strings.Join(widgets, " | ")
	log.Debug(status)

	return []rune(status)
}
