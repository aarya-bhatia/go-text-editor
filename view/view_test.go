package view

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestView_AddText(t *testing.T) {
	box := NewViewBuffer(0, 0, 10, 2)
	lines := [][]rune{
		[]rune("Hello"),
		[]rune("World"),
	}
	box.AddText(lines)
	expected := [][]rune{
		[]rune("Hello     "),
		[]rune("World     "),
	}
	assert.True(t, reflect.DeepEqual(box.Text, expected))
}

func TestView_AddBorder(t *testing.T) {
	box := NewViewBuffer(0, 0, 10, 4)
	box.AddBorder()
	expected := [][]rune{
		[]rune("----------"),
		[]rune("|        |"),
		[]rune("|        |"),
		[]rune("----------"),
	}
	assert.True(t, reflect.DeepEqual(box.Text, expected))
}

func TestView_AddViewBuffer(t *testing.T) {
	box1 := NewViewBuffer(0, 0, 10, 4)
	box1.AddBorder()

	box2 := NewViewBuffer(1, 1, 8, 2)
	box2.AddText([][]rune{
		[]rune("Hello"),
		[]rune("World"),
	})

	box1.Add(box2)

	expected := [][]rune{
		[]rune("----------"),
		[]rune("|Hello   |"),
		[]rune("|World   |"),
		[]rune("----------"),
	}

	assert.True(t, box1.Height == 4)
	assert.True(t, box1.Width == 10)
	assert.True(t, box1.TopX == 0)
	assert.True(t, box1.TopY == 0)

	assert.True(t, reflect.DeepEqual(box1.Text, expected))
}
