package utils

import (
	"fmt"
	"log"
	"runtime"
	"unicode/utf8"
)

func Assert(value bool, message ...interface{}) {
	if !value {

		// Get the caller information
		_, file, line, ok := runtime.Caller(1)
		if ok {
			log.Printf("Panic at %s:%d - %s\n", file, line, message)
		}
		panic(message)
	}
}

func FlattenList[T any](array [][]T) []T {
	res := []T{}
	for _, element := range array {
		res = append(res, element...)
	}
	return res
}

func GetStringMatrix(nRows int, nCols int, char rune) [][]rune {
	lines := make([][]rune, nRows)

	for y := 0; y < nRows; y++ {
		lines[y] = make([]rune, nCols)
		for x := 0; x < nCols; x++ {
			lines[y][x] = char
		}
	}

	return lines
}

func Map[T any, R any](elements []T, mapper func(element T) R) []R {
	newElements := make([]R, len(elements))
	for i, element := range elements {
		newElements[i] = mapper(element)
	}
	return newElements
}

func Filter[T any](elements []T, filter func(element T) bool) []T {
	newElements := make([]T, 0, len(elements)) // pre-allocate max capacity
	for i, element := range elements {
		if filter(elements[i]) {
			newElements = append(newElements, element)
		}
	}
	return newElements
}

func Find[T any](elements []T, predicate func(element T) bool) *T {
	for i, element := range elements {
		if predicate(element) {
			return &elements[i] // return actual pointer not copy
		}
	}

	return nil
}

func IndexOf[T comparable](elements []T, target T) int {
	for i, element := range elements {
		if element == target {
			return i
		}
	}

	return -1
}

func ArrayRemove[T any](slice []T, index int) []T {
	newSlice := make([]T, 0, len(slice))
	for i, element := range slice {
		if i == index {
			continue
		}

		newSlice = append(newSlice, element)
	}
	return newSlice
}

// RightAlign trims the left side if needed and right-aligns text within the given width.
func RightAlign(text string, maxWidth int) string {
	textLen := utf8.RuneCountInString(text)

	if textLen > maxWidth {
		// Clip leading characters to fit within maxWidth
		runes := []rune(text)
		return string(runes[textLen-maxWidth:])
	}

	// Pad with spaces to right-align
	return fmt.Sprintf("%*s", maxWidth, text)
}
