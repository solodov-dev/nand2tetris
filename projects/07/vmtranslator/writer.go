package main

import (
	"log"
	"os"
	"strconv"
)

type CodeWriter struct {
	output       *os.File
	returnIdx    int
	compareCount int
}

var BaseAddresses = map[string]string{
	"local":    "LCL",
	"argument": "ARG",
	"this":     "THIS",
	"that":     "THAT",
	"pointer":  "R3",
	"temp":     "R5",
}

func NewCodeWriter(output string) *CodeWriter {
	writeFile, err := os.Create(output)

	if err != nil {
		log.Fatalf("Cannot create new file %s", output)
	}

	return &CodeWriter{writeFile, 0, 0}
}

func (w *CodeWriter) Close() error {
	return w.output.Close()
}

func (w *CodeWriter) WriteArithmetic(arg string) {
	switch arg {
	case "add":
		w.BinaryCommand("M=D+M")
	case "sub":
		w.BinaryCommand("M=M-D")
	case "neg":
		w.UnaryCommand("M=-M")
	case "eq":
		w.CompareCommand("JEQ")
	case "gt":
		w.CompareCommand("JGT")
	case "lt":
		w.CompareCommand("JLT")
	case "and":
		w.BinaryCommand("M=M&D")
	case "or":
		w.BinaryCommand("M=M|D")
	case "not":
		w.UnaryCommand("M=!M")
	}
}

func (w *CodeWriter) WritePush(segment string, val string, currentFile string) {
	switch segment {
	case "constant":
		w.NumToD(val)
		w.PushDToStack()
		w.IncrementStackPointer()
	case "static":
		w.Push(currentFile + "." + val)
	case "pointer":
		base, err := strconv.Atoi(val)
		HandleError(err, "Pointer index is not a value")
		w.Push(strconv.Itoa(base + 3))
	case "temp":
		base, err := strconv.Atoi(val)
		HandleError(err, "Pointer index is not a value")
		w.Push(strconv.Itoa(base + 5))
	default:
		seg, ok := BaseAddresses[segment]

		if !ok {
			log.Fatalf("Segment %s is not in known segments list", segment)
		}

		w.MToD(seg)
		w.Write("@" + val)
		w.Write("A=D+A")
		w.Write("D=M")
		w.PushDToStack()
		w.IncrementStackPointer()
	}
}

func (w *CodeWriter) Push(val string) {
	w.MToD(val)
	w.PushDToStack()
	w.IncrementStackPointer()
}

func (w *CodeWriter) WritePop(segment string, val string, currentFile string) {
	switch segment {
	case "static":
		w.DecrementStackPointer()
		w.PopStackToD()
		w.DToM(currentFile + "." + val)
	case "pointer":
		base, err := strconv.Atoi(val)
		HandleError(err, "Index is not a number")
		w.DecrementStackPointer()
		w.PopStackToD()
		w.DToM(strconv.Itoa(base + 3))
	case "temp":
		base, err := strconv.Atoi(val)
		HandleError(err, "Index is not a number")
		w.DecrementStackPointer()
		w.PopStackToD()
		w.DToM(strconv.Itoa(base + 5))
	default:
		seg, ok := BaseAddresses[segment]

		if !ok {
			log.Fatalf("Segment %s is not in known segments list", segment)
		}

		w.MToD(seg)
		w.Write("@" + val)
		w.Write("D=D+A")
		w.DToM("R13")
		w.DecrementStackPointer()
		w.PopStackToD()
		w.Write("@R13")
		w.Write("A=M")
		w.Write("M=D")
	}

}

func (w *CodeWriter) WriteEnd() {
	w.Write("(END)")
	w.Write("@END")
  w.UnconditionalJump()
}

func (w *CodeWriter) Write(line string) {
	_, err := w.output.WriteString(line + "\n")

	if err != nil {
		log.Fatalf("Could not write line %s to a file %s", line, w.output.Name())
	}
}

func (w *CodeWriter) WriteLabel(label string, f string) {
	w.Write("(" + f + "$" + label + ")")
}

func (w *CodeWriter) WriteGoTo(label string, f string) {
	w.Write("@" + f + "$" + label)
  w.UnconditionalJump()
}

// False = 0; True anything else
func (w *CodeWriter) WriteIfGoTo(label string, f string) {
	w.DecrementStackPointer()
	w.PopStackToD()
	w.Write("@" + f + "$" + label)
	w.Write("D;JNE")
}

func (w *CodeWriter) WriteFunction(funcName string, numLocals string) {
	locals, err := strconv.Atoi(numLocals)

	if err != nil {
		log.Fatalf("Could not parse function %s locals %s. Not a number.", funcName, numLocals)
	}

	line := `(` + funcName + `)`

	for i := 0; i < locals; i++ {
		w.NumToD("0")
		w.PushDToStack()
		w.IncrementStackPointer()
	}

	w.Write(line)
}

func (w *CodeWriter) WriteRestore() {
	w.MToD("R13")
	w.Write("D=D-1")
	w.DToM("R13")
	w.Write("A=D")
	w.Write("D=M")
}

func (w *CodeWriter) WriteReturn() {
	w.MoveMemToMem("LCL", "R13")
	w.Write("@5")
	w.Write("A=D-A")
	w.Write("D=M")
  w.DToM("R14")
	w.DecrementStackPointer()
	w.Write("@ARG")
	w.Write("AD=M")
	w.DToM("R15")
  w.PopStackToD()
	w.Write("@R15")
	w.Write("A=M")
	w.Write("M=D")
	w.MToD("R2")
	w.Write("@R0")
	w.Write("M=D+1")
	w.WriteRestore()
	w.DToM("THAT")
	w.WriteRestore()
	w.DToM("THIS")
	w.WriteRestore()
	w.DToM("ARG")
	w.WriteRestore()
	w.DToM("LCL")
	w.Write("@R14")
	w.Write("A=M")
  w.UnconditionalJump()
}

func (w *CodeWriter) WriteCall(f string, nArgs string) {
	w.returnIdx++
	count := strconv.Itoa(w.returnIdx)

	w.MoveMemToMem("SP", "R13")
	w.Write("@RETURN_" + count)
	w.Write("D=A")
  w.PushDToStack()
	w.IncrementStackPointer()
	w.SavePointer("LCL")
	w.IncrementStackPointer()
	w.SavePointer("ARG")
	w.IncrementStackPointer()
	w.SavePointer("THIS")
	w.IncrementStackPointer()
	w.SavePointer("THAT")
	w.IncrementStackPointer()
	w.MToD("R13")
	w.Write("@" + nArgs)
	w.Write("D=D-A")
	w.DToM("ARG")
	w.MToD("SP")
	w.DToM("LCL")
	w.Write("@" + f)
  w.UnconditionalJump()
	w.Write("(RETURN_" + count + ")")
}

func (w *CodeWriter) SavePointer(p string) {
	w.MToD(p)
	w.PushDToStack()
}

// Pushes the numeric value into D register
func (w *CodeWriter) NumToD(n string) {
	w.Write("@" + n)
	w.Write("D=A")
}

// Pushes the value located at address adr to D register
func (w *CodeWriter) MToD(adr string) {
	w.Write("@" + adr)
	w.Write("D=M")
}

// Pushes the value located at D register to address adr
func (w *CodeWriter) DToM(adr string) {
	w.Write("@" + adr)
	w.Write("M=D")
}

// Pushes value stored in D register to stack
func (w *CodeWriter) PushDToStack() {
	w.Write("@SP")
	w.Write("A=M")
	w.Write("M=D")
}

// Pops value stored on a stack to D register
func (w *CodeWriter) PopStackToD() {
  w.Write("@SP")
	w.Write("A=M")
	w.Write("D=M")
}

// Moves value from memory A to memory B
func (w *CodeWriter) MoveMemToMem(memA string, memB string) {
	w.Write("@" + memA)
	w.Write("D=M")
	w.Write("@" + memB)
	w.Write("M=D")
}

// Decrements stack pointer
func (w *CodeWriter) DecrementStackPointer() {
	w.Write("@SP")
	w.Write("M=M-1")
}

// Increments stack pointer
func (w *CodeWriter) IncrementStackPointer() {
	w.Write("@SP")
	w.Write("M=M+1")
}

// Create a binary command on the stack
// Puts X in M register and Y in D register
// Stack pointer points at X
//
//         -----------
//  SP ->  |    x    |    in M
//         -----------
//         |    y    |    in D
//         -----------
// After command is added stack pointer will point at X + 1
func (w *CodeWriter) BinaryCommand(cmd string) {
	w.DecrementStackPointer()
  w.PopStackToD()
	w.DecrementStackPointer()
	w.Write("A=M")
	w.Write(cmd)
	w.IncrementStackPointer()
}

// Create a one arg command on the stack
// Puts Y in M register
// Stack pointer points at Y
//
//         -----------
//         |    x    |
//         -----------
//  SP ->  |    y    |    in M
//         -----------
// After command is added stack pointer will point at X + 1
func (w *CodeWriter) UnaryCommand(cmd string) {
	w.DecrementStackPointer()
	w.Write("A=M")
	w.Write(cmd)
	w.IncrementStackPointer()
}

func (w *CodeWriter) CompareCommand(cmd string) {
	w.compareCount++
	index := strconv.Itoa(w.compareCount)
	w.BinaryCommand("M=M-D")
	w.DecrementStackPointer()
  w.PopStackToD()
	w.Write("@TRUE" + index)
	w.Write("D;" + cmd)
	w.Write("D=0")
	w.Write("@WRITE" + index)
  w.UnconditionalJump()
	w.Write("(TRUE" + index + ")")
	w.Write("D=-1")
	w.Write("(WRITE" + index + ")")
  w.PushDToStack()
	w.IncrementStackPointer()
}

func (w *CodeWriter) WriteInit(callSysInit bool) {
	w.Write("@256\nD=A\n@SP\nM=D")
	// if sys.vm provided call Sys.init 0
	if callSysInit {
		w.WriteCall("Sys.init", "0")
    w.UnconditionalJump()
	}
}

func (w *CodeWriter) UnconditionalJump() {
  w.Write("0;JMP")
}
