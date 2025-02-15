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
