package util

import "testing"

func TestSomeNumber(t *testing.T) {
	if SomeNumber() != 0 {
		t.Fail()
	}
}
