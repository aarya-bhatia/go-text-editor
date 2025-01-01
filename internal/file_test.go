package internal

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const TEMP_FILE = "example"

func TestReadWriteFile(t *testing.T) {
	f, err := os.CreateTemp("", TEMP_FILE)
	assert.Nil(t, err)

	defer f.Close()

	file := NewFile(f.Name())
	assert.NotNil(t, file, "file should not be nil")

	assert.NoError(t, file.ReadFile(), "read file should work")

	file.Lines = []*Line{{Text: "Hello"}, {Text: "World"}}
	assert.NoError(t, file.WriteFile(), "write file should work")

	content, err := os.ReadFile(f.Name())
	assert.Nil(t, err)
	assert.Equal(t, string(content), "Hello\nWorld", "unexpected file content")
}

func TestNewFile(t *testing.T) {
	tmp_file, err := os.CreateTemp("", TEMP_FILE)
	assert.Nil(t, err)
	var file = NewFile(tmp_file.Name())
	assert.NotNil(t, file)
	assert.True(t, len(file.Lines) == 0)
}

func TestFileRead(t *testing.T) {
	tmp_file, _ := os.CreateTemp("", TEMP_FILE)
	assert.Nil(t, os.WriteFile(tmp_file.Name(), []byte("Hello"), DEFAULT_FILE_PERMISSIONS))
	var file = NewFile(tmp_file.Name())
	assert.Nil(t, file.ReadFile())
}

func TestWriteFile(t *testing.T) {
}

func TestFile_ReadFile(t *testing.T) {

}

func TestFile_WriteFile(t *testing.T) {

}

func TestFile_InsertLineAboveCursor(t *testing.T) {

}

func TestFile_InsertLineBelowCursor(t *testing.T) {

}

func TestFile_Paste(t *testing.T) {

}

func TestFile_DeleteLine(t *testing.T) {

}

func TestFile_GetCurrentLine(t *testing.T) {

}
