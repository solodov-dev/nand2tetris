package main

import (
	"log"
	"strconv"
)

// x + y
func add() string {
	return twoArgCommand("M=D+M")
}

// x - y
func sub() string {
	return twoArgCommand("M=M-D")
}

// -y
func neg() string {
	return oneArgCommand("M=-M")
}

// x = y
// true -1
// false 0
func eq() string {
	return compareCommand("JEQ")
}

// if x > y
// true -1
// false 0
func gt() string {
	return compareCommand("JGT")
}

// if x < y
// true -1
// false 0
func lt() string {
	return compareCommand("JLT")
}

// x and y
func and() string {
	return twoArgCommand("M=M&D")
}

// x or y
func or() string {
	return twoArgCommand("M=M|D")
}

// not y
func not() string {
	return oneArgCommand("M=!M")
}

// Create a two arg command on the stack
// Puts X in M register and Y in D register
// Stack pointer points at X
//
//         -----------
//  SP ->  |    x    |    in M
//         -----------
//         |    y    |    in D
//         -----------
// After command is added stack pointer will point at X + 1
func twoArgCommand(cmd string) string {
	return `@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
` + cmd + `
@SP
M=M+1`
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
func oneArgCommand(cmd string) string {
	return `@SP
M=M-1
A=M
` + cmd + `
@SP
M=M+1`
}

var compareCount = 0

// Compare the last stack value with the compare command
// If false write 0 to M
// If true write -1 to M
func compareCommand(compare string) string {
	compareCount++
	index := strconv.Itoa(compareCount)
	return sub() + "\n" + `@SP
M=M-1
A=M
D=M
@TRUE` + index + `
D;` + compare + `
D=0
@WRITE` + index + `
0;JMP
(TRUE` + index + `)
D=-1
(WRITE` + index + `)
@SP
A=M
M=D
@SP
M=M+1`
}

var Arithmetic = map[string]func() string{
	"add": add,
	"sub": sub,
	"neg": neg,
	"eq":  eq,
	"lt":  lt,
	"gt":  gt,
	"and": and,
	"or":  or,
	"not": not,
}

var BaseAddresses = map[string]string{
	"local":    "@LCL",
	"argument": "@ARG",
	"this":     "@THIS",
	"that":     "@THAT",
	"pointer":  "@R3",
	"temp":     "@R5",
}

// PushConstant pushes value val onto stack
func PushConstant(val string) string {
	return `@` + val + `
D=A
` + pushDtoStack()
}

func PushStatic(index string, filename string) string {
	return `@` + filename + `.` + index + `
D=M
` + pushDtoStack()
}

func PushTempPointer(val string, addr int) string {
	idx, err := strconv.Atoi(val)

	if err != nil {
		log.Fatalf("Value %s is not an index", val)
	}

	addr += idx

	return `@` + strconv.Itoa(addr) + `
  D=M
  ` + pushDtoStack()
}

func Pop(segment string, index string) string {
	seg, ok := BaseAddresses[segment]

	if !ok {
		log.Fatalf("Segment %s is not in known segments list", segment)
	}

	return seg + `
D=M
@` + index + `
D=A+D
@R13
M=D
` + popStackToD() + `
@R13
A=M
M=D`
}

func PopTempPointer(index string, addr int) string {
	idx, err := strconv.Atoi(index)

	if err != nil {
		log.Fatalf("Value %s is not an index", index)
	}

	addr += idx

  return popStackToD() + `
@` + strconv.Itoa(addr) + `
M=D`
}

func PopStatic(index string, filename string) string {
	return popStackToD() + `
@` + filename + `.` + index + `
M=D
`
}

func Push(segment string, index string) string {
	seg, ok := BaseAddresses[segment]

	if !ok {
		log.Fatalf("Segment %s is not in known segments list", segment)
	}

	return seg + ` 
D=M
@` + index + `
A=D+A
D=M
` + pushDtoStack()
}

// Pushes the contents of D register to stack
func pushDtoStack() string {
	return `@SP
A=M
M=D
@SP
M=M+1`
}

// Pops stack and stores the value in D
func popStackToD() string {
	return `@SP
M=M-1
A=M
D=M`
}
