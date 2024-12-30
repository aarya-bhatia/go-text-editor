package internal

import (
	"log"
	"os"
	"strings"
)

type File struct {
	Lines    []Line
	Name     string
	Readonly bool
	Modified bool
	Cursor   Cursor
}

func NewFile(filename string) *File {
	var this *File = new(File)
	this.Name = filename
	this.Modified = false
	this.Readonly = false
	this.Lines = make([]Line, 1)
	this.Cursor.X = 0
	this.Cursor.Y = 0
	return this
}

const DEFAULT_FILE_PERMISSIONS os.FileMode = 0664

func (this *File) ReadFile() error {
	fptr, err := os.OpenFile(this.Name, os.O_CREATE|os.O_RDONLY, DEFAULT_FILE_PERMISSIONS)
	if err != nil {
		return err
	}

	defer fptr.Close()

	stat, err := fptr.Stat()
	if err != nil {
		return err
	}

	read_buffer := make([]byte, stat.Size())
	n, err := fptr.Read(read_buffer)
	if err != nil {
		return err
	}

	log.Printf("Read %d bytes from file %s", n, this.Name)

	raw_lines := strings.Split(string(read_buffer), "\n")
	for _, line := range raw_lines {
		this.Lines = append(this.Lines, Line{Text: line})
	}

	log.Printf("Parsed %d lines", len(this.Lines))

	return nil
}

func (this *File) WriteFile() error {
	fptr, err := os.OpenFile(this.Name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, DEFAULT_FILE_PERMISSIONS)
	if err != nil {
		return err
	}

	defer fptr.Close()

	raw_lines := make([]string, 0)
	for _, line := range this.Lines {
		raw_lines = append(raw_lines, line.Text)
	}

	write_buffer := []byte(strings.Join(raw_lines, "\n"))
	n, err := fptr.Write(write_buffer)
	if err != nil {
		return err
	}

	log.Printf("Written %d bytes to file %s", n, this.Name)
	return nil
}
