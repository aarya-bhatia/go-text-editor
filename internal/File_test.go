package internal

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFile(t *testing.T) {
	file := NewFile("testdata/test")
	assert.NotNil(t, file)
}

func TestReadWriteFile(t *testing.T) {
	file := NewFile("testdata/test")
	assert.NotNil(t, file)
	assert.NoError(t, file.ReadFile())

	file.Lines = []Line{{Text: "Hello"}, {Text: "World"}}
	assert.NoError(t, file.WriteFile())

	content, err := ioutil.ReadFile("testdata/test")
	assert.NoError(t, err)
	assert.Equal(t, string(content), "Hello\nWorld\n")
}
