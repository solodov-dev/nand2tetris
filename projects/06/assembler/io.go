package main

import (
	"log"
	"os"
)

type IO struct {
	readFile  *os.File
	writeFile *os.File
}

func NewIO(input string, output string) *IO {
	readFile, err := os.Open(input)

	if err != nil {
		log.Fatalf("Cannot open file %s", input)
	}

	writeFile, err := os.Create(output)

	if err != nil {
		log.Fatalf("Cannot create new file %s", output)
	}

	return &IO{readFile, writeFile}
}

func (f *IO) Close() error {
	return f.readFile.Close()
}

func (f *IO) Write(line string) {
	_, err := f.writeFile.WriteString(line + "\n")

	if err != nil {
		log.Fatalf("Could not write line %s to a new file %s", line, f.writeFile.Name())
	}
}
