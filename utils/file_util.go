package utils

import (
	"log"
	"os"
	"strings"
)

const DEFAULT_FILE_PERMISSIONS os.FileMode = 0640

func ReadFileUtil(filename string) ([]string, error) {
	fptr, err := os.OpenFile(filename, os.O_CREATE|os.O_RDONLY, DEFAULT_FILE_PERMISSIONS)
	if err != nil {
		return nil, err
	}

	defer fptr.Close()

	stat, err := fptr.Stat()
	if err != nil {
		return nil, err
	}

	read_buffer := make([]byte, stat.Size())
	n, err := fptr.Read(read_buffer)
	if err != nil {
		return nil, err
	}

	log.Printf("Read %d bytes from file %s", n, filename)
	return strings.Split(string(read_buffer), "\n"), nil
}

func WriteFileUtil(filename string, lines []string) error {
	fptr, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, DEFAULT_FILE_PERMISSIONS)
	if err != nil {
		return err
	}

	defer fptr.Close()

	write_buffer := []byte(strings.Join(lines, "\n"))
	n, err := fptr.Write(write_buffer)
	if err != nil {
		return err
	}

	log.Printf("Written %d bytes to file %s", n, filename)
	return nil
}
