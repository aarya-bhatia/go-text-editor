package internal

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const TEMP_FILE = "/tmp/testfile"

func TestApplication_NoFiles(t *testing.T) {
	var app *Application = NewApplication()
	assert.True(t, app != nil)

	// Verify the app is initialized with correct values and there are no files
	assert.True(t, app.CurrentFile == nil)
	assert.True(t, len(app.Files) == 0)
	assert.True(t, app.QuitSignal == false)
	assert.True(t, app.Mode == NORMAL_MODE)

	// Verify app can handle edge cases
	app.OpenPrevFile()
	assert.True(t, app.CurrentFile == nil)
	app.OpenNextFile()
	assert.True(t, app.CurrentFile == nil)
	app.CloseFile()
	assert.True(t, app.CurrentFile == nil)
}

func TestApplication_SingleFile(t *testing.T) {
	defer os.Remove(TEMP_FILE)

	var app *Application = NewApplication()

	// Verify it can open one file
	app.OpenFile(TEMP_FILE)
	assert.True(t, app.CurrentFile != nil)
	assert.True(t, len(app.Files) == 1)
	assert.True(t, app.CurrentFile.Name == TEMP_FILE)
	file := app.CurrentFile

	// Verify it does not create new buffer if file exists
	app.OpenFile(TEMP_FILE)
	assert.True(t, app.CurrentFile == file)
	assert.True(t, len(app.Files) == 1)
	assert.True(t, app.CurrentFile.Name == TEMP_FILE)

	// Verify it does not change current file when moving back and forth
	app.OpenNextFile()
	assert.True(t, app.CurrentFile == file)
	app.OpenPrevFile()
	assert.True(t, app.CurrentFile == file)
	app.CloseFile()

	// Verify it closes the file
	assert.True(t, app.CurrentFile == nil)
	assert.True(t, len(app.Files) == 0)

}

func TestApplication_MultipleFiles(t *testing.T) {
	f1, _ := os.CreateTemp("", "*.0")
	f2, _ := os.CreateTemp("", "*.1")
	f3, _ := os.CreateTemp("", "*.2")

	defer f1.Close()
	defer f2.Close()
	defer f3.Close()

	defer os.Remove(f1.Name())
	defer os.Remove(f2.Name())
	defer os.Remove(f3.Name())

	var app *Application = NewApplication()

	// Verify it opens first file
	app.OpenFile(f1.Name())
	assert.True(t, app.CurrentFile != nil)
	assert.True(t, len(app.Files) == 1)
	assert.True(t, app.CurrentFile.Name == f1.Name())
	f1_file := app.CurrentFile

	// Verify it opens second file
	app.OpenFile(f1.Name())
	app.OpenFile(f2.Name())
	assert.True(t, app.CurrentFile != nil)
	assert.True(t, len(app.Files) == 2)
	assert.True(t, app.CurrentFile.Name == f2.Name())
	f2_file := app.CurrentFile

	// Verify it opens third file
	app.OpenFile(f3.Name())
	assert.True(t, app.CurrentFile != nil)
	assert.True(t, len(app.Files) == 3)
	assert.True(t, app.CurrentFile.Name == f3.Name())
	f3_file := app.CurrentFile

	// Verify it opens the first buffer when opening the first file
	app.OpenFile(f1.Name())
	assert.True(t, app.CurrentFile != nil)
	assert.True(t, len(app.Files) == 3)
	assert.True(t, app.CurrentFile.Name == f1.Name())
	assert.True(t, app.CurrentFile == f1_file)

	// Verify it opens the second buffer when opening the second file
	app.OpenFile(f2.Name())
	assert.True(t, app.CurrentFile == f2_file)
	assert.True(t, len(app.Files) == 3)

	// Verify it opens the third buffer when opening the third file
	app.OpenFile(f3.Name())
	assert.True(t, app.CurrentFile == f3_file)
	assert.True(t, len(app.Files) == 3)

	log.Printf("f1: %s, f2: %s, f3: %s", f1.Name(), f2.Name(), f3.Name())

	// Verify it opens the files in the sequence they were added - 1 - 2 - 3
	// and it should circle back to the start
	app.OpenNextFile()
	assert.True(t, app.CurrentFile == f1_file)
	app.OpenNextFile()
	assert.True(t, app.CurrentFile == f2_file)
	app.OpenNextFile()
	assert.True(t, app.CurrentFile == f3_file)
	app.OpenPrevFile()
	log.Println("Current file: ", app.CurrentFile.Name)
	assert.True(t, app.CurrentFile == f2_file)
	app.OpenPrevFile()
	assert.True(t, app.CurrentFile == f1_file)
	app.OpenPrevFile()
	assert.True(t, app.CurrentFile == f3_file)

	// Verify it closes the third file and replaces current file with the next one, which is f1
	app.CloseFile()
	assert.True(t, len(app.Files) == 2)
	assert.True(t, app.CurrentFile == f1_file)

	// Verify it closes the first file and replaces current file with f2
	app.CloseFile()
	assert.True(t, len(app.Files) == 1)
	assert.True(t, app.CurrentFile == f2_file)

	// Verify all files are closed and current file is nil
	app.CloseFile()
	assert.True(t, len(app.Files) == 0)
	assert.True(t, app.CurrentFile == nil)
}
