package model

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLine_NewLine(t *testing.T) {
	line := NewLine()
	assert.True(t, line.Cursor == 0)
	assert.True(t, line.Modified == false)
	assert.True(t, string(line.Text) == "")

	line = NewLine1("Hello World")
	assert.True(t, line.Cursor == 0)
	assert.True(t, line.Modified == false)
	assert.True(t, string(line.Text) == "Hello World")
	assert.True(t, line.Size() == len("Hello World"))
}

func TestLine_MoveCursor(t *testing.T) {
	line := NewLine1("Hello World")
	assert.True(t, line.Cursor == 0)
	line.MoveToEnd()
	assert.True(t, line.Cursor == len("Hello World")-1)
	line.MoveToStart()
	assert.True(t, line.Cursor == 0)
	line.MoveForward()
	assert.True(t, line.Cursor == 1)
	line.Cursor = len("Hello World") - 1
	line.MoveForward()
	assert.True(t, line.Cursor == len("Hello World")-1)
	line.MoveBackward()
	assert.True(t, line.Cursor == len("Hello World")-2)
	line.MoveToStart()
	line.MoveBackward()
	assert.True(t, line.Cursor == 0)
}

func TestLine_Insert(t *testing.T) {
	line := NewLine1("H")
	assert.True(t, line.Cursor == 0)
	line.Cursor = 1
	line.Insert('e')
	assert.True(t, string(line.Text) == "He")
	assert.True(t, line.Cursor == 2)
	assert.True(t, line.Modified == true)
	line.Insert('l')
	assert.True(t, string(line.Text) == "Hel")
	assert.True(t, line.Cursor == 3)
	assert.True(t, line.Modified == true)
	line.Insert('l')
	assert.True(t, string(line.Text) == "Hell")
	assert.True(t, line.Cursor == 4)
	assert.True(t, line.Modified == true)
}

func TestLine_RemoveChar(t *testing.T) {
	t.Run("remove char from middle", func(t *testing.T) {
		line := NewLine1("Hello World")
		line.Cursor = 5
		line.RemoveChar()
		assert.True(t, string(line.Text) == "HelloWorld")
		assert.True(t, line.Cursor == 5)
		assert.True(t, line.Modified == true)
	})

	t.Run("remove char from empty line", func(t *testing.T) {
		line := NewLine()
		line.RemoveChar()
		assert.True(t, string(line.Text) == "")
		assert.True(t, line.Cursor == 0)
		assert.True(t, line.Modified == false)
	})

	t.Run("remove char from start", func(t *testing.T) {
		line := NewLine1("Hello World")
		line.Cursor = 0
		line.RemoveChar()
		assert.True(t, string(line.Text) == "Hello World")
		assert.True(t, line.Cursor == 0)
		assert.True(t, line.Modified == false)
	})

	t.Run("remove char from end", func(t *testing.T) {
		line := NewLine1("Hello World")
		line.Cursor = len("Hello World") - 1
		line.RemoveChar()
		assert.True(t, string(line.Text) == "Hello Worl")
		assert.True(t, line.Cursor == len("Hello Worl")-1)
		assert.True(t, line.Modified == true)
	})
}

func TestLine_Append(t *testing.T) {
	line := NewLine1("Hello")
	line.MoveToEnd()
	line.Append(' ')
	assert.True(t, string(line.Text) == "Hello ")
	assert.True(t, line.Cursor == len("Hello ")-1)
	line.Append('W')
	assert.True(t, string(line.Text) == "Hello W")
	assert.True(t, line.Cursor == len("Hello W")-1)
	line.Append('o')
	line.Append('r')
	line.Append('l')
	line.Append('d')
	assert.True(t, string(line.Text) == "Hello World")
	assert.True(t, line.Cursor == len("Hello World")-1)
	assert.True(t, line.Modified == true)
}

func TestLine_AppendString(t *testing.T) {
	line := NewLine1("Hello")
	line.Cursor = 4
	line.AppendString(" World")
	log.Println(string(line.Text))
	assert.True(t, string(line.Text) == "Hello World")
	assert.True(t, line.Cursor == len("Hello World")-1)
	assert.True(t, line.Modified == true)
}

func TestLine_JumpToNextChar(t *testing.T) {
	line := NewLine1("Hello World")
	line.JumpToNextChar(' ')
	assert.True(t, line.Cursor == 5)
	line.JumpToNextChar(' ')
	assert.True(t, line.Cursor == 5)
	line.JumpToNextChar('W')
	assert.True(t, line.Cursor == 6)
	line.JumpToNextChar('o')
	assert.True(t, line.Cursor == 7)
	line.JumpToNextChar('r')
	assert.True(t, line.Cursor == 8)
	line.JumpToNextChar('l')
	assert.True(t, line.Cursor == 9)
	line.JumpToNextChar('d')
	assert.True(t, line.Cursor == 10)
	line.JumpToNextChar(' ')
	assert.True(t, line.Cursor == 10)
}
