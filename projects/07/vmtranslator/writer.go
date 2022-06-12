package main

import (
	"log"
	"os"
)

type CodeWriter struct {
	output *os.File
}

func NewCodeWriter(output string) *CodeWriter {
	writeFile, err := os.Create(output)

	if err != nil {
		log.Fatalf("Cannot create new file %s", output)
	}

	return &CodeWriter{writeFile}
}

func (w *CodeWriter) Close() error {
	return w.output.Close()
}

func (w *CodeWriter) WriteArithmetic(arg string) {
	f, ok := Arithmetic[arg]

	if !ok {
		log.Fatalf("Cannot find command %s in arithmetic commands", arg)
	}

	line := f()
	w.Write(line)
}

func (w *CodeWriter) WritePush(segment string, val string, currentFile string) {
	var line string

	switch segment {
	case "constant":
		line = PushConstant(val)
	case "static":
    line = PushStatic(val, currentFile)
  default:
    line = Push(segment, val)
	}

	w.Write(line)
}

func (w *CodeWriter) WriteEnd() {
	line := `(END)
@END
0;JMP
`
	w.Write(line)
}

func (w *CodeWriter) Write(line string) {
	_, err := w.output.WriteString(line + "\n")

	if err != nil {
		log.Fatalf("Could not write line %s to a file %s", line, w.output.Name())
	}
}
