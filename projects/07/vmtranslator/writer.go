package main

import (
	"log"
	"os"
)

type CodeWriter struct {
	output *os.File
	parser *Parser
}

func NewCodeWriter(output string, p *Parser) *CodeWriter {
	writeFile, err := os.Create(output)

	if err != nil {
		log.Fatalf("Cannot create new file %s", output)
	}

	return &CodeWriter{writeFile, p}
}

func (w *CodeWriter) Close() error {
	return w.output.Close()
}

func (w *CodeWriter) WriteArithmetic() {
	command := w.parser.arg1

  f, ok := Arithmetic[command]

  if !ok {
    log.Fatalf("Cannot find command %s in arithmetic commands", command)
  }

  line := f()
	w.Write(line)
}


func (w *CodeWriter) WritePush() {
	segment := w.parser.arg1
	val := w.parser.arg2
  var line string

  switch segment {
  case "constant":
   line = Push(val)    
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
		log.Fatalf("Could not write line %d to a file %s", w.parser.lineNumber, w.output.Name())
	}
}
