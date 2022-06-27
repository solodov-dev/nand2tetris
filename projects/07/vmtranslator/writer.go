package main

import (
	"log"
	"os"
	"strconv"
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

func (w *CodeWriter) WriteFunction(funcName string, numLocals string) {
	locals, err := strconv.Atoi(numLocals)

	if err != nil {
		log.Fatalf("Could not parse function %s locals %s. Not a number.", funcName, numLocals)
	}

	line := `(` + funcName + `)`

	for i := 0; i < locals; i++ {
		line += `
` + PushConstant("0")
	}

	w.Write(line)
}

func (w *CodeWriter) WriteReturn() {
	line := `@LCL  // FRAME=LCL  Save LCL in a temp variable
D=M
@R13
M=D
@5  // RET=*(FRAME-5)  Put the return address in a temp var R14
A=D-A
D=M
@R14
M=D
@SP // *ARG=pop()  Reposition the return value
M=M-1
@ARG
AD=M
@R15
M=D
@SP
A=M
D=M
@R15
A=M
M=D
@R2
D=M
@R0
M=D+1
@R13  // THAT=*(FRAME-1)  Restore THAT of the caller
D=M
D=D-1
@R13
M=D
A=D
D=M
@THAT
M=D
@R13  // THIS=*(FRAME-2)  Restore THIS of the caller
D=M
D=D-1
@R13
M=D
A=D
D=M
@THIS
M=D
@R13  // ARG=*(FRAME-3)  Restore ARG of the caller
D=M
D=D-1
@R13
M=D
A=D
D=M
@ARG
M=D
@R13  // LCL=*(FRAME-4)  Restore LCL of the caller
D=M
D=D-1
@R13
M=D
A=D
D=M
@LCL
M=D
@R14
A=M
0;JMP`

	w.Write(line)
}
