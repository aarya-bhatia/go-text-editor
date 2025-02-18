package utils

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtils_FlattenList(t *testing.T) {
	a := [][]int{
		{1, 2, 3},
		{4, 5},
		{},
		{6},
	}

	b := FlattenList(a)

	expected := []int{1, 2, 3, 4, 5, 6}
	assert.True(t, reflect.DeepEqual(expected, b))
}

func TestUtils_GetStringMatric(t *testing.T) {
	a := GetStringMatrix(2, 2, '.')

	aTransform := []string{}

	for _, row := range a {
		aTransform = append(aTransform, string(row))
	}

	expected := []string{
		"..",
		"..",
	}

	assert.True(t, reflect.DeepEqual(expected, aTransform))
}

type TestData struct {
	value int
}

func TestUtil_Map(t *testing.T) {
	elements := []TestData{
		{value: 1},
		{value: 2},
		{value: 3},
	}

	newElements := Map(elements, func(element TestData) TestData {
		return TestData{value: element.value * 2}
	})

	expected := []TestData{
		{value: 2},
		{value: 4},
		{value: 6},
	}

	assert.True(t, reflect.DeepEqual(expected, newElements))
}

func TestUtil_Filter(t *testing.T) {
	elements := []TestData{
		{value: 1},
		{value: 2},
		{value: 3},
	}

	newElements := Filter(elements, func(element TestData) bool {
		return element.value%2 == 0
	})

	expected := []TestData{
		{value: 2},
	}

	assert.True(t, reflect.DeepEqual(expected, newElements))
}

func TestUtil_Find(t *testing.T) {
	elements := []TestData{
		{value: 1},
		{value: 2},
		{value: 3},
	}

	found := Find(elements, func(element TestData) bool {
		return element.value == 1
	})

	assert.True(t, found != nil)
	assert.True(t, found == &elements[0])

	found = Find(elements, func(element TestData) bool {
		return element.value == 0
	})

	assert.True(t, found == nil)
}
