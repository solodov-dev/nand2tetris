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
	case "pointer":
		line = PushTempPointer(val, 3)
	case "temp":
		line = PushTempPointer(val, 5)
	default:
		line = Push(segment, val)
	}

	w.Write(line)
}

func (w *CodeWriter) WritePop(segment string, val string, currentFile string) {
	var line string

	switch segment {
	case "static":
		line = PopStatic(val, currentFile)
	case "pointer":
		line = PopTempPointer(val, 3)
	case "temp":
		line = PopTempPointer(val, 5)
	default:
		line = Pop(segment, val)
	}

	w.Write(line)
}

func (w *CodeWriter) WriteEnd() {
	line := `(END)
@END
0;JMP`

	w.Write(line)
}

func (w *CodeWriter) Write(line string) {
	_, err := w.output.WriteString(line + "\n")

	if err != nil {
		log.Fatalf("Could not write line %s to a file %s", line, w.output.Name())
	}
}

func (w *CodeWriter) WriteLabel(label string, f string) {
	line := `(` + f + `$` + label + `)`
	w.Write(line)
}

func (w *CodeWriter) WriteGoTo(label string, f string) {
	line := `@` + f + `$` + label + `
0;JMP`

	w.Write(line)
}

// False = 0; True anything else
func (w *CodeWriter) WriteIfGoTo(label string, f string) {
	line := popStackToD() + `
@` + f + `$` + label + `
D;JNE`

	w.Write(line)
}
