package internal

import "errors"

func ErrorOutOfBounds() error {
  return errors.New("Out of bounds")
}

func ErrorFileNotModifiable() error {
  return errors.New("File is not modifiable")
}

func ErrorIllegalArguments() error {
  return errors.New("Illegal Arguments")
}
