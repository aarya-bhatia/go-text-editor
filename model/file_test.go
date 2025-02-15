package model

import (
	"log"
	"os"
	"strings"
	"testing"
  "go-editor/utils"

	"github.com/stretchr/testify/assert"
)

var testFile *File = nil

// setup test case to initialise testFile as a file struct with a temporary file path.
// Returns callback function to clean up test file before next test case.
func setupTest(t *testing.T) func(t *testing.T) {
	tempDir := t.TempDir()
	log.Println("Created temp dir: ", tempDir)
	testFile = NewFile(tempDir + "/tempfile")
	return func(t *testing.T) {
		if testFile == nil {
			t.Fatal("expected test file to be non-nil")
		}

		os.Remove(testFile.Name)
		testFile = nil
	}
}

func setupTestWithFileContent(t *testing.T, content string) func(t *testing.T) {
	tearDown := setupTest(t)

	lines := strings.Split(content, "\n")
	testFile.Lines = make([]*Line, 0)
	for _, line := range lines {
		testFile.Lines = append(testFile.Lines, NewLine1(line))
	}

	return tearDown
}

// Verify if a new file is created without error and no lines
func TestFile_NewFile(t *testing.T) {
	tearDown := setupTest(t)
	defer tearDown(t)

	assert.NotNil(t, testFile)
	assert.True(t, len(testFile.Lines) == 0)
}

// Writes to temp file on disk and expect ReadFile() to read content into file struct.
func TestFile_ReadFile(t *testing.T) {
	tearDown := setupTest(t)
	defer tearDown(t)

	assert.Nil(t, os.WriteFile(testFile.Name, []byte("Hello"), utils.DEFAULT_FILE_PERMISSIONS))
	assert.Nil(t, testFile.ReadFile())
}

// Verify that a file struct can update the contents of a real file.
func TestFile_WriteFile(t *testing.T) {
	tearDown := setupTest(t)
	defer tearDown(t)

	testFile.Lines = []*Line{{Text: []rune("Hello")}, {Text: []rune("World")}}
	assert.NoError(t, testFile.WriteFile())

	content, err := os.ReadFile(testFile.Name)
	assert.Nil(t, err)
	assert.Equal(t, string(content), "Hello\nWorld")
}

func TestFile_GetCurrentLine(t *testing.T) {
	tearDown := setupTestWithFileContent(t, "First line\nSecond line\nThird line")
	defer tearDown(t)

	testFile.CursorLine = 0
	assert.True(t, testFile.GetCurrentLine().Cursor == 0)
	assert.True(t, string(testFile.GetCurrentLine().Text) == "First line")

	testFile.CursorLine = 1
	assert.True(t, testFile.GetCurrentLine().Cursor == 0)
	assert.True(t, string(testFile.GetCurrentLine().Text) == "Second line")

	testFile.CursorLine = 2
	assert.True(t, testFile.GetCurrentLine().Cursor == 0)
	assert.True(t, string(testFile.GetCurrentLine().Text) == "Third line")

	testFile.CursorLine = 3
	assert.True(t, testFile.GetCurrentLine() == nil)

	testFile.Lines = []*Line{}
	testFile.CursorLine = 0
	assert.True(t, testFile.GetCurrentLine() == nil)
}

func TestFile_InsertLineAboveCursor(t *testing.T) {
	tearDown := setupTestWithFileContent(t, "First line\nSecond line\nThird line")
	defer tearDown(t)

	assert.True(t, len(testFile.Lines) == 3)
	testFile.CursorLine = 1
	assert.True(t, string(testFile.GetCurrentLine().Text) == "Second line")
	testFile.InsertLineAboveCursor()
	assert.True(t, len(testFile.Lines) == 4)
	assert.True(t, testFile.CursorLine == 2)
	assert.True(t, string(testFile.GetCurrentLine().Text) == "Second line")

	testFile.CursorLine = 0
	assert.True(t, string(testFile.GetCurrentLine().Text) == "First line")
	testFile.InsertLineAboveCursor()
	assert.True(t, len(testFile.Lines) == 5)
	assert.True(t, testFile.CursorLine == 1)
	assert.True(t, string(testFile.GetCurrentLine().Text) == "First line")
}

func TestFile_InsertLineBelowCursor(t *testing.T) {
	t.Run("insert line below second line", func(t *testing.T) {
		tearDown := setupTestWithFileContent(t, "First line\nSecond line\nThird line")
		defer tearDown(t)

		testFile.CursorLine = 1
		assert.True(t, string(testFile.GetCurrentLine().Text) == "Second line")
		testFile.InsertLineBelowCursor()
		assert.True(t, len(testFile.Lines) == 4)
		assert.True(t, testFile.CursorLine == 1)
		assert.True(t, string(testFile.GetCurrentLine().Text) == "Second line")
	})

	t.Run("insert line below first line", func(t *testing.T) {
		tearDown := setupTestWithFileContent(t, "First line\nSecond line\nThird line")
		defer tearDown(t)

		testFile.CursorLine = 0
		assert.True(t, string(testFile.GetCurrentLine().Text) == "First line")
		testFile.InsertLineBelowCursor()
		assert.True(t, len(testFile.Lines) == 4)
		assert.True(t, testFile.CursorLine == 0)
		assert.True(t, string(testFile.GetCurrentLine().Text) == "First line")
	})

	t.Run("insert line below third line", func(t *testing.T) {
		tearDown := setupTestWithFileContent(t, "First line\nSecond line\nThird line")
		defer tearDown(t)

		testFile.CursorLine = 2
		assert.True(t, string(testFile.GetCurrentLine().Text) == "Third line")
		testFile.InsertLineBelowCursor()
		assert.True(t, len(testFile.Lines) == 4)
		assert.True(t, testFile.CursorLine == 2)
		assert.True(t, string(testFile.GetCurrentLine().Text) == "Third line")
	})
}

func TestFile_DeleteLine(t *testing.T) {
	tearDown := setupTestWithFileContent(t, "First line\nSecond line\nThird line")
	defer tearDown(t)

	assert.True(t, len(testFile.Lines) == 3)
	testFile.CursorLine = 1
	assert.True(t, string(testFile.GetCurrentLine().Text) == "Second line")
	testFile.DeleteLine()
	assert.True(t, len(testFile.Lines) == 2)
	assert.True(t, testFile.CursorLine == 0)
	assert.True(t, string(testFile.GetCurrentLine().Text) == "First line")

	testFile.DeleteLine()
	assert.True(t, len(testFile.Lines) == 1)
	assert.True(t, testFile.CursorLine == 0)
	assert.True(t, string(testFile.GetCurrentLine().Text) == "Third line")

	testFile.DeleteLine()
	assert.True(t, len(testFile.Lines) == 0)
	assert.True(t, testFile.CursorLine == 0)
	assert.True(t, testFile.GetCurrentLine() == nil)

	testFile.DeleteLine()
	assert.True(t, len(testFile.Lines) == 0)
	assert.True(t, testFile.CursorLine == 0)
	assert.True(t, testFile.GetCurrentLine() == nil)
}
