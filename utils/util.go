package utils

import (
	"log"
	"runtime"
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
